# grpc服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 先决条件

1. 安装protoc编译器

    从 https://github.com/protocolbuffers/protobuf/releases 下载protoc编译器, 解压protoc.exe到$GOPATH/bin/

2. 安装 ProtoBuffer Golang 支持

   ```shell
   go install github.com/golang/protobuf/protoc-gen-go@latest`
   ```

3. 安装 ProtoBuffer GRpc Golang 支持. [文档](https://grpc.io/docs/languages/go/quickstart/)

   ```shell
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

4. 数据校验需要安装 [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate)

   ```shell
   go install github.com/envoyproxy/protoc-gen-validate@latest
   ```

# 示例项目

+ [grpc服务端](./example/grpc-s)
+ [grpc客户端](https://github.com/zly-app/component/tree/master/grpc-client/example/grpc-c)

# 快速开始

1. 创建一个项目

   ```shell
   mkdir server && cd server
   go mod init server
   ```

2. 添加 `hello/hello.proto` 文件

   ```protobuf
   syntax = 'proto3';
   package hello; // 决定proto引用路径和rpc路由
   option go_package = "server/hello/hello"; // 用于对golang包管理的定位
   
   service helloService{
      rpc Hello(HelloReq) returns (HelloResp);
   }
   
   message HelloReq{
      string msg = 1;
   }
   message HelloResp{
      string msg = 1;
   }
   ```

3. 编译 proto
   
   ```shell
   protoc \
   --gogoslick_out=. --gogoslick_opt=paths=source_relative \
   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   hello/hello.proto
   ````

4. 添加 `main.go` 文件

   ```go
   package main
   
   import (
       "context"
   
       "github.com/zly-app/zapp"
       "github.com/zly-app/zapp/core"
       "go.uber.org/zap"
   
       "github.com/zly-app/service/grpc"
   
       "server/hello"
   )
   
   // 服务实现
   type HelloService struct {
       hello.UnimplementedHelloServiceServer
   }
   
   func (h *HelloService) Hello(ctx context.Context, req *hello.HelloReq) (*hello.HelloResp, error) {
       session := grpc.GetSession(ctx) // 获取session
       session.Info("收到请求", zap.String("req", req.Msg))
       return &hello.HelloResp{Msg: req.GetMsg() + "world"}, nil
   }
   
   func main() {
       app := zapp.NewApp("test",
           grpc.WithService(), // 启用grpc服务
       )
   
       // 注册服务对象
       grpc.RegistryServerObject(func(c core.IComponent, server grpc.ServiceRegistrar) {
           hello.RegisterHelloServiceServer(server, new(HelloService)) // 注册服务
       })
   
       app.Run()
   }
   ```

5. 运行

   ```shell
   go mod tidy && go run .
   ```

# 配置文件

添加配置文件 `configs/default.yml`. 更多配置参考[这里](./config.go)

```yaml
services:
   grpc:
      bind: ":3000"
```
