/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/5/9
   Description :
-------------------------------------------------
*/

package api

import (
	"fmt"
	"reflect"

	"github.com/kataras/iris/v12"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"
)

var typeOfContext = reflect.TypeOf((*Context)(nil))
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

func wrapX(handler interface{}, isMiddleware bool) iris.Handler {
	if handler == nil {
		logger.Log.Fatal("handler为nil", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	if h, ok := handler.(Handler); ok {
		return Wrap(h)
	}

	hType := reflect.TypeOf(handler)
	if hType.Kind() != reflect.Func {
		logger.Log.Fatal("handler必须是函数", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 检查入参
	if hType.NumIn() < 1 || hType.NumIn() > 2 {
		logger.Log.Fatal("handler的入参数量为1个或2个", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 检查第一个参数
	arg0 := hType.In(0)
	if !arg0.AssignableTo(typeOfContext) {
		logger.Log.Fatal("handler的第一个入参必须是 *api.Context", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 检查req参数
	var reqCreator func(ctx *Context) (reflect.Value, error) // 请求参数创建者, 用于实例化请求参数
	if hType.NumIn() == 2 {
		arg1 := hType.In(1) // 获取req的类型
		reqIsPtr := arg1.Kind() == reflect.Ptr
		if reqIsPtr {
			arg1 = arg1.Elem() // 获取req的真实类型
		}
		if arg1.Kind() != reflect.Struct {
			logger.Log.Fatal("handler的第二个入参必须是 struct 或 *struct", zap.String("handler", fmt.Sprintf("%T", handler)))
		}

		// req创建者
		reqCreator = func(ctx *Context) (reflect.Value, error) {
			req := reflect.New(arg1)                          // 创建req实例
			if err := ctx.Bind(req.Interface()); err != nil { // bind参数
				return reflect.Value{}, err
			}

			if !reqIsPtr { // 如果req不是指针, 则返回指针指向的值
				return req.Elem(), nil
			}
			return req, nil
		}
	}

	// 检查出参
	if hType.NumOut() < 1 || hType.NumOut() > 2 {
		logger.Log.Fatal("handler的出参数量为1个或2个", zap.String("handler", fmt.Sprintf("%T", handler)))
	}

	// 如果出参数为2个, 最后一个出参必须是error
	if hType.NumOut() == 2 {
		out1 := hType.Out(1)
		if !out1.AssignableTo(typeOfError) {
			logger.Log.Fatal("handler的第二个出参必须是 error", zap.String("handler", fmt.Sprintf("%T", handler)))
		}
	}

	wrapFunc := Wrap
	if isMiddleware {
		wrapFunc = WrapMiddleware
	}

	hValue := reflect.ValueOf(handler)
	return wrapFunc(func(ctx *Context) interface{} {
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

		// 检查参数
		if len(outValues) == 1 {
			return outValues[0].Interface()
		}

		err := outValues[1].Interface()
		if err != nil {
			return err.(error)
		}
		return outValues[0].Interface()
	})
}

// 包装处理程序
//
// handler 是一个 func
//      第一个入参必须是 *api.Context 类型, 第二个入参是可选的 struct, 它可以是指针
//      第一个出参可以是任何类型, 第二个出参必须是error类型, 第二个出参是可选的
//      示例:
//          func (ctx *api.Context) interface{}
//          func (ctx *api.Context) error
//          func (ctx *api.Context, req *ReqStruct) interface{}
//          func (ctx *api.Context, req *ReqStruct) error
//          func (ctx *api.Context, req *ReqStruct) (interface{}, error)
//          func (ctx *api.Context, req *ReqStruct) (*OutStruct, error)
func WrapX(handler interface{}) iris.Handler {
	return wrapX(handler, false)
}

// 包装中间件, 类似 WrapX, 只有返回nil才能继续调用链, 非nil值表示拦截, 并将结果处理后返回给客户端
func WrapMiddlewareX(handler interface{}) iris.Handler {
	return wrapX(handler, true)
}
