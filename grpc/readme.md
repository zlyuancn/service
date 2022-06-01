# grpc服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 先决条件

1. 安装protoc编译器

    从 https://github.com/protocolbuffers/protobuf/releases 下载protoc编译器, 解压protoc.exe到$GOPATH/bin/

2. 安装 ProtoBuffer Golang 支持. 推荐使用 [gogoslick](https://github.com/gogo/protobuf)
    + gogoslick<sup>**推荐**</sup> `go install github.com/gogo/protobuf/protoc-gen-gogoslick@latest`
    + go `go install github.com/golang/protobuf/protoc-gen-go@latest`

3. 安装 ProtoBuffer GRpc Golang 支持. [文档](https://grpc.io/docs/languages/go/quickstart/)

   ```shell
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

4. 如果启用http协议, 需要安装 [protoc-gen-grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

   ```shell
   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
   ```

# 示例项目

+ [grpc服务端](./example/grpc-s)
+ [grpc客户端](https://github.com/zly-app/component/tree/master/grpc-client/example/grpc-c)

# 配置

添加配置文件 `configs/default.yml`. 更多配置参考[这里](./config.go)

```yaml
services:
   grpc:
      bind: ":3000"
```
