# hello

hello

## protobuf

generate protobuf

### gen go

protoc-gen-go

```shell script

protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/go-kratos/kratos/third_party --go_opt=paths=source_relative --go_out=. ./*.proto

protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/go-kratos/kratos/third_party --go_opt=paths=source_relative --go_out=plugins=grpc:. ./*.proto

```

### gen gogo

protoc-gen-gogofast (same as gofast, but imports gogoprotobuf)
protoc-gen-gogofaster (same as gogofast, without XXX_unrecognized, less pointer fields)
protoc-gen-gogoslick (same as gogofaster, but with generated string, gostring and equal methods)

```shell script

protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/go-kratos/kratos/third_party --gofast_out=. *.proto

protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/go-kratos/kratos/third_party --gofast_out=plugins=grpc:. *.proto

protoc -I. -I%GOPATH%/src -I%GOPATH%/src/github.com/go-kratos/kratos/third_party --gogofaster_out=plugins=grpc:. *.proto

```

### gen kratos

```shell script

kratos tool protoc --grpc api.proto

```
