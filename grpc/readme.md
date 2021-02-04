
# grpc服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

```text
grpc.RegistryService()          # 注册服务
grpc.WithGrpcService()          # 启用服务
grpc.RegistryGrpcServerObject(...)    # 服务注入(注册服务实体)
```

# 示例

```protobuf
// test.proto
syntax = "proto3";
option go_package = "./pb";
service Test{
  rpc Test(TestReq) returns (TestResp);
}
message TestReq{
  string Arg = 1;
}
message TestResp{
  string Result = 1;
}
```

```go
type TestService struct {
    c core.IComponent
}

func NewTestService(c core.IComponent, service *grpc.GrpcServer) {
    t := &TestService{c: c}
    pb.RegisterTestServer(service, t)
}
func (t *TestService) Test(ctx context.Context, req *pb.TestReq) (*pb.TestResp, error) {
    log := grpc.MustGetLoggerFromContext(ctx) // 获取会话logger
    log.Info("grpc处理过程中的日志")
    return &pb.TestResp{Result: req.Arg + "已处理"}, nil
}

func main() {
    // 注册服务
    grpc.RegistryService()
    // 启用服务
    app := zapp.NewApp("test", grpc.WithGrpcService())
    // 服务注入
    grpc.RegistryGrpcServerObject(app, NewTestService)
    // 启动
    app.Run()
}
```

# 配置

> 默认服务类型为 `grpc`

```toml
[services.grpc]
# bind地址
Bind="localhost:3000"
# 心跳时间(毫秒), 默认20000
HeartbeatTime=20000
```
