/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package grpc

import (
	"context"
	"net"
	"reflect"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/opentracing/opentracing-go"
	open_log "github.com/opentracing/opentracing-go/log"
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ServiceRegistrar = grpc.ServiceRegistrar

var typeOfComponentConnInterface = reflect.TypeOf((*core.IComponent)(nil)).Elem()
var typeOfGrpcServer = reflect.TypeOf((*ServiceRegistrar)(nil)).Elem()

type RegistryGrpcServiceFunc = func(c core.IComponent, server ServiceRegistrar)

type GrpcService struct {
	app    core.IApp
	server *grpc.Server
	conf   *Config
}

func NewGrpcService(app core.IApp) core.IService {
	conf := newConfig()
	err := app.GetConfig().ParseServiceConfig(nowServiceType, conf, true)
	if err != nil {
		app.Fatal("创建服务失败", zap.String("serviceType", string(nowServiceType)), zap.Error(err))
	}
	conf.Check()

	chainUnaryClientList := []grpc.UnaryServerInterceptor{}

	if conf.EnableOpenTrace {
		chainUnaryClientList = append(chainUnaryClientList, UnaryServerOpenTraceInterceptor)
	}
	chainUnaryClientList = append(chainUnaryClientList,
		UnaryServerLogInterceptor(app, conf),   // 日志
		grpc_ctxtags.UnaryServerInterceptor(),  // 设置标记
		grpc_recovery.UnaryServerInterceptor(), // panic恢复
	)
	if conf.ReqDataValidate && !conf.ReqDataValidateAllField {
		chainUnaryClientList = append(chainUnaryClientList, UnaryServerReqDataValidateInterceptor)
	}
	if conf.ReqDataValidate && conf.ReqDataValidateAllField {
		chainUnaryClientList = append(chainUnaryClientList, UnaryServerReqDataValidateAllInterceptor)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(chainUnaryClientList...)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: time.Duration(conf.HeartbeatTime) * time.Millisecond, // 心跳
		}),
	)

	// 在app关闭前优雅的关闭服务
	zapp.AddHandler(zapp.BeforeExitHandler, func(app core.IApp, handlerType zapp.HandlerType) {
		server.GracefulStop()
		app.Warn("grpc服务已关闭")
	})

	return &GrpcService{
		app:    app,
		server: server,
		conf:   conf,
	}
}

func (g *GrpcService) Inject(a ...interface{}) {
	cValue := reflect.ValueOf(g.app.GetComponent())
	serverValue := reflect.ValueOf(g.server)
	for _, v := range a {
		vType := reflect.TypeOf(v)
		if vType.Kind() != reflect.Func {
			g.app.Fatal("grpc服务端注入参数必须是func")
			return
		}
		if vType.NumIn() != 2 {
			g.app.Fatal("grpc服务端注入func入参为2个")
			return
		}
		arg0 := vType.In(0)
		if !arg0.AssignableTo(typeOfComponentConnInterface) {
			g.app.Fatal("注入的func第1个入参必须能转为 core.IComponent")
			return
		}
		arg1 := vType.In(1)
		if !arg1.AssignableTo(typeOfGrpcServer) {
			g.app.Fatal("注入的func第2个入参必须能转为 *grpc.GrpcServer")
			return
		}

		vValue := reflect.ValueOf(v)
		vValue.Call([]reflect.Value{cValue, serverValue})
	}
}

func (g *GrpcService) Start() error {
	listener, err := net.Listen("tcp", g.conf.Bind)
	if err != nil {
		return err
	}

	g.app.Info("正在启动grpc服务", zap.String("bind", g.conf.Bind))
	return g.server.Serve(listener)
}

func (g *GrpcService) Close() error {
	return nil
}

// 日志拦截器
func UnaryServerLogInterceptor(app core.IApp, conf *Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := app.NewTraceLogger(ctx, zap.String("grpc.method", info.FullMethod))
		se := &Session{
			ILogger: log,
		}
		ctx = context.WithValue(ctx, sessionContextKey, se)

		startTime := time.Now()
		if conf.ReqLogLevelIsInfo {
			log.Info("grpc.request", zap.Any("req", req))
		} else {
			log.Debug("grpc.request", zap.Any("req", req))
		}

		reply, err := handler(ctx, req)
		if err != nil {
			log.Error("grpc.response", zap.String("latency", time.Since(startTime).String()), zap.Error(err))
			return reply, err
		}

		if conf.RspLogLevelIsInfo {
			log.Info("grpc.response", zap.String("latency", time.Since(startTime).String()), zap.Any("reply", reply))
		} else {
			log.Debug("grpc.response", zap.String("latency", time.Since(startTime).String()), zap.Any("reply", reply))
		}

		return reply, err
	}
}

type TextMapCarrier struct {
	metadata.MD
}

func (t TextMapCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range t.MD {
		for _, vv := range v {
			if err := handler(k, vv); err != nil {
				return err
			}
		}
	}
	return nil
}

// 开放链路追踪hook
func UnaryServerOpenTraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 取出元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// 如果对元数据修改必须使用它的副本
		md = md.Copy()
	}

	// 从元数据中取出span
	carrier := TextMapCarrier{md}
	parentSpan, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, carrier)

	span := opentracing.StartSpan(info.FullMethod, opentracing.ChildOf(parentSpan))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	span.LogFields(open_log.Object("req", req))
	reply, err := handler(ctx, req)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(open_log.Error(err))
	} else {
		span.LogFields(open_log.Object("reply", reply))
	}
	return reply, err
}

type ValidateInterface interface {
	Validate() error
}
type ValidateAllInterface interface {
	ValidateAll() error
}

// 数据校验
func UnaryServerReqDataValidateInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if v, ok := req.(ValidateInterface); ok {
		if err := v.Validate(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}

// 数据校验, 总是校验所有字段
func UnaryServerReqDataValidateAllInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	v, ok := req.(ValidateAllInterface)
	if !ok {
		// 降级
		return UnaryServerReqDataValidateInterceptor(ctx, req, info, handler)
	}
	if err := v.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return handler(ctx, req)
}
