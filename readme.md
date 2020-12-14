# srv_toolkit

toolkit for srv

## protobuf

```shell script

# 安装 protoc-gen-go
go get google.golang.org/protobuf/cmd/protoc-gen-go

# 切换目录
cd ./api

# linux
protoc -I. -I$GOSRCPATH -I$GOPBPATH --go_out=. --go_opt=paths=source_relative ./*.proto

# windows
protoc -I. -I%GOSRCPATH% -I%GOPBPATH% --go_out=. --go_opt=paths=source_relative ./*.proto

```