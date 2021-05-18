/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/30
   Description :
-------------------------------------------------
*/

package api

import (
	"fmt"
	"reflect"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"

	app_config "github.com/zly-app/zapp/config"

	"github.com/zly-app/service/api/config"
)

// 处理程序
//
// 如果返回bytes会直接返回给客户端
// 返回其它值会经过处理后再返回给客户端
type Handler = func(ctx *Context) interface{}

type Response struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

var typeOfContext = reflect.TypeOf((*Context)(nil))
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// 检查 handler 指纹
func checkHandlerFingerprint(handler interface{}, handlerType reflect.Type) {
	if handlerType.Kind() != reflect.Func {
		logger.Log.Fatal("handler必须是函数", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 检查入参
	if handlerType.NumIn() < 1 || handlerType.NumIn() > 2 {
		logger.Log.Fatal("handler的入参数量为1个或2个", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 检查第一个参数
	arg0 := handlerType.In(0)
	if !arg0.AssignableTo(typeOfContext) {
		logger.Log.Fatal("handler的第一个入参必须是 *api.Context", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 检查出参
	if handlerType.NumOut() < 1 || handlerType.NumOut() > 2 {
		logger.Log.Fatal("handler的出参数量为1个或2个", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 如果出参数为2个, 最后一个出参必须是error
	if handlerType.NumOut() == 2 {
		out1 := handlerType.Out(1)
		if !out1.AssignableTo(typeOfError) {
			logger.Log.Fatal("handler的第二个出参必须是 error", zap.String("handler", fmt.Sprintf("%T", handler)))
		}
	}
}

// 根据 handler 构建req建造者
//
// req 是 handler 的第二个入参, 如果入参数量小于 2 返回 nil
// req 必须是 struct 或 *struct
func mustMakeReqCreator(handler interface{}, handlerType reflect.Type) func(ctx *Context) (reflect.Value, error) {
	if handlerType.NumIn() < 2 {
		return nil
	}

	// 检查 req 是 struct 或 *struct
	arg1 := handlerType.In(1)              // 获取req的类型
	reqIsPtr := arg1.Kind() == reflect.Ptr // req参数是否为指针
	if reqIsPtr {
		arg1 = arg1.Elem() // 获取req的真实类型
	}
	if arg1.Kind() != reflect.Struct {
		logger.Log.Fatal("handler的第二个入参必须是 struct 或 *struct", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 返回建造者
	return func(ctx *Context) (reflect.Value, error) {
		req := reflect.New(arg1)                          // 创建req实例
		if err := ctx.Bind(req.Interface()); err != nil { // bind参数
			return req, err
		}

		if reqIsPtr { // 如果req是指针, 直接返回
			return req, nil
		}
		return req.Elem(), nil // 非指针要返回指向对象
	}
}

// 构建handler
func makeHandler(handler interface{}, handlerType reflect.Type) Handler {
	hValue := reflect.ValueOf(handler)
	reqCreator := mustMakeReqCreator(handler, handlerType)

	h := func(ctx *Context) interface{} {
		var outValues []reflect.Value

		// 调用handler
		if reqCreator == nil { // 如果没有req建造者, 表示不需要req参数
			outValues = hValue.Call([]reflect.Value{reflect.ValueOf(ctx)})
		} else {
			reqValue, err := reqCreator(ctx)
			if err != nil {
				return err
			}
			outValues = hValue.Call([]reflect.Value{reflect.ValueOf(ctx), reqValue})
		}

		// 检查结果
		if len(outValues) == 1 { // 如果只有一个结果直接返回
			return outValues[0].Interface()
		}

		err := outValues[1].Interface()
		if err != nil {
			return err.(error)
		}
		return outValues[0].Interface()
	}
	return h
}

// 包装处理程序
func wrap(handler interface{}, isMiddleware bool) iris.Handler {
	if handler == nil {
		logger.Log.Fatal("handler为nil", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	h, ok := handler.(Handler)
	if !ok {
		hType := reflect.TypeOf(handler)
		checkHandlerFingerprint(handler, hType) // 检查handler指纹
		h = makeHandler(handler, hType)         // 构建handler
	}

	return func(irisCtx *iris_context.Context) {
		ctx := makeContext(irisCtx) // 构建上下文
		result := h(ctx)            // 处理

		// 如果是中间件, 只有返回nil才能继续调用链, 非nil值表示拦截, 并将结果处理后返回给客户端
		if isMiddleware && result == nil { // 返回nil继续调用链
			ctx.Next()
			return
		}

		WriteToCtx(ctx, result) // 写入结果
		ctx.StopExecution()     // 停止调用链
	}
}

// 写入数据到ctx
//
// 如果返回bytes会直接返回给客户端
// 返回其它值会经过处理后再返回给客户端
func WriteToCtx(ctx *Context, result interface{}) {
	if err, ok := result.(error); ok {
		ctx.Values().Set("error", err)
		code, message := decodeErr(err)

		isProduction := !app_config.Conf.Config().Frame.Debug
		showDetailedErrorInProduction := config.Conf.ShowDetailedErrorInProduction
		if !isProduction || showDetailedErrorInProduction {
			message = err.Error()
		}
		_, _ = ctx.JSON(Response{
			ErrCode: code,
			ErrMsg:  message,
		})
		return
	}

	switch v := result.(type) {
	case []byte: // 直接写入
		ctx.Values().Set("result", fmt.Sprintf("bytes<len=%d>", len(v)))
		_, _ = ctx.Write(v)
	case *[]byte: // 直接写入
		ctx.Values().Set("result", fmt.Sprintf("bytes<len=%d>", len(*v)))
		_, _ = ctx.Write(*v)
	default:
		ctx.Values().Set("result", result)
		_, _ = ctx.JSON(Response{
			ErrCode: OK.Code,
			ErrMsg:  OK.Message,
			Data:    result,
		})
	}
}

// 包装处理程序
//
// handler 是一个 func
//      入参: 第一个入参必须是 *api.Context 类型, 如果有第二个入参必须是 struct, 第二个入参可以是指针
//      出参: 第一个出参可以是任何类型, 如果有第二个出参必须是error类型
//      示例:
//          func (ctx *api.Context) interface{}
//          func (ctx *api.Context) error
//          func (ctx *api.Context, req *ReqStruct) interface{}
//          func (ctx *api.Context, req *ReqStruct) error
//          func (ctx *api.Context, req *ReqStruct) (interface{}, error)
//          func (ctx *api.Context, req *ReqStruct) (*OutStruct, error)
func Wrap(handler interface{}) iris.Handler {
	return wrap(handler, false)
}

// 包装中间件, 类似 Wrap, 只有返回nil才能继续调用链, 非nil值表示拦截, 并将结果处理后返回给客户端
func WrapMiddleware(handler interface{}) iris.Handler {
	return wrap(handler, true)
}
