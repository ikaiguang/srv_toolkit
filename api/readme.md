# protobuf

protoc-gen-go google.golang.org/protobuf@v1.25.0

```shell script

# 安装protoc-gen-go
go get google.golang.org/protobuf/cmd/protoc-gen-go

protoc -I. -I%GOPATH%/src --go_out=. --go_opt=paths=source_relative ./*.proto
protoc -I. -I$GOPATH/src --go_out=. --go_opt=paths=source_relative ./*.proto

```