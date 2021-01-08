# generate

```shell script

# 安装编译工具 protoc-gen-go@v1.25.0 protoc-gen-go-grpc@v1.4.3
go get github.com\ikaiguang\protoc-gen-go\cmd\protoc-gen-go-tkform
go get github.com\ikaiguang\protoc-gen-go\cmd\protoc-gen-go-tkgrpc

# 添加环境变量 $GOSRCPATH $GOPBPATH
# export GOSRCPATH=$GOPATH/src
# export GOPBPATH=$GOPATH/src/github.com/go-kratos/kratos/third_party

# 切换到需要编译的文件夹
# cd $GOPATH/src/github.com/ikaiguang/srv_hello/api/xxx

# google.golang.org/protobuf@v1.25.0 + google.golang.org/grpc@v1.34.0 
# google.golang.org/protobuf@v1.25.0_form + github.com/golang/protobuf@1.4.3_grpc
# protoc-gen-go-tkform + protoc-gen-go-tkgrpc
protoc -I. -I%GOSRCPATH% -I%GOPBPATH% --go-tkform_out=. --go-tkform_opt=paths=source_relative --go-tkgrpc_out=. --go-tkgrpc_opt=paths=source_relative ./*.proto
protoc -I. -I$GOSRCPATH -I$GOPBPATH --go-tkform_out=. --go-tkform_opt=paths=source_relative --go-tkgrpc_out=. --go-tkgrpc_opt=paths=source_relative ./*.proto

# generate BM HTTP
kratos tool protoc --bm *.proto

```

## openapi v2

> 参考链接 [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
>
> 参考例子 [a_bit_of_everything.proto](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/internal/proto/examplepb/a_bit_of_everything.proto)

```shell script

# 安装 openapiv2
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

protoc -I. -I$GOSRCPATH -I$GOPBPATH --openapiv2_out ./ --openapiv2_opt logtostderr=true ./*.proto

```

## kratos tool protoc

```shell script

# generate all
#kratos tool protoc api.proto
kratos tool protoc *.proto
# generate gRPC
#kratos tool protoc --grpc api.proto
kratos tool protoc --grpc *.proto
# generate BM HTTP
#kratos tool protoc --bm api.proto
kratos tool protoc --bm *.proto
# generate ecode
#kratos tool protoc --ecode api.proto
kratos tool protoc --ecode *.proto
# generate swagger
#kratos tool protoc --swagger api.proto
kratos tool protoc --swagger *.proto

```