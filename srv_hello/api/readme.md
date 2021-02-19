# generate

```shell script

# 安装编译工具 protoc-gen-go@v1.25.0 protoc-gen-go-grpc@v1.4.3
go get github.com/ikaiguang/protoc-gen-go/cmd/protoc-gen-go-tkform
go get github.com/ikaiguang/protoc-gen-go/cmd/protoc-gen-go-tkgrpc

# 添加环境变量 $GOSRCPATH $GOPBPATH
# export GOSRCPATH=$GOPATH/src
# export GOPBPATH=$GOPATH/src/github.com/go-kratos/kratos/third_party

# 切换到需要编译的文件夹
# cd $GOPATH/src/github.com/ikaiguang/xxx/api/xxx

# google.golang.org/protobuf@v1.25.0 + google.golang.org/grpc@v1.34.0 
# google.golang.org/protobuf@v1.25.0_form + github.com/golang/protobuf@1.4.3_grpc
# protoc-gen-go-tkform + protoc-gen-go-tkgrpc
protoc -I. -I%GOSRCPATH% -I%GOPBPATH% --go-tkform_out=. --go-tkform_opt=paths=source_relative --go-tkgrpc_out=. --go-tkgrpc_opt=paths=source_relative ./*.proto
protoc -I. -I$GOSRCPATH -I$GOPBPATH --go-tkform_out=. --go-tkform_opt=paths=source_relative --go-tkgrpc_out=. --go-tkgrpc_opt=paths=source_relative ./*.proto

```

## 格式化代码

```shell script
clang-format -i -style="{BasedOnStyle: llvm, IndentWidth: 2, ColumnLimit: 1000}" ./*.proto
```