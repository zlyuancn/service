package main

import (
	"context"

	"github.com/zly-app/service/grpc"
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"

	"grpc-s/pb/hello"
)

var _ hello.HelloServiceServer = (*HelloService)(nil)

type HelloService struct{}

func (h *HelloService) Hello(ctx context.Context, req *hello.HelloReq) (*hello.HelloResp, error) {
	return &hello.HelloResp{Msg: req.GetMsg() + "world"}, nil
}

func main() {
	app := zapp.NewApp("grpc-s",
		grpc.WithService())

	grpc.RegistryServerObject(func(c core.IComponent, server *grpc.GrpcServer) {
		hello.RegisterHelloServiceServer(server, new(HelloService)) // 注册服务
	})

	app.Run()
}
