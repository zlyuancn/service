package main

import (
	"context"

	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/service/grpc"

	"github.com/zly-app/service/grpc/example/grpc-s/pb/hello"
)

var _ hello.HelloServiceServer = (*HelloService)(nil)

type HelloService struct {
	hello.UnimplementedHelloServiceServer
}

func (h *HelloService) Hello(ctx context.Context, req *hello.HelloReq) (*hello.HelloResp, error) {
	session := grpc.GetSession(ctx) // 获取session
	session.Info("收到请求", zap.String("req", req.Msg))
	return &hello.HelloResp{Msg: req.GetMsg() + "world"}, nil
}

func main() {
	app := zapp.NewApp("grpc-s",
		grpc.WithService())

	grpc.RegistryServerObject(func(c core.IComponent, server grpc.ServiceRegistrar) {
		hello.RegisterHelloServiceServer(server, new(HelloService)) // 注册服务
	})

	app.Run()
}
