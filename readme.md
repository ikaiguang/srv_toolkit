# srv_toolkit

toolkit for srv

- google.golang.org/protobuf v1.25.0
- google.golang.org/grpc v1.29.1

## protobuf

protoc-gen-go google.golang.org/protobuf@v1.25.0

```shell script

# 安装 protoc-gen-go
go get google.golang.org/protobuf/cmd/protoc-gen-go

protoc -I. -I%GOPATH%/src --go_out=. --go_opt=paths=source_relative ./*.proto
protoc -I. -I$GOPATH/src --go_out=. --go_opt=paths=source_relative ./*.proto

```

## 安装 protoc-gen-go-tkform protoc-gen-go-tkgrpc
   
```shell script

# 安装编译工具 protoc-gen-go@v1.25.0 protoc-gen-go-grpc@v1.4.3
go get github.com/ikaiguang/protoc-gen-go/cmd/protoc-gen-go-tkform
go get github.com/ikaiguang/protoc-gen-go/cmd/protoc-gen-go-tkgrpc

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

## 运行测试

go run ./srv_hello/main.go
