# srv_hello

hello

## protobuf 说明

请查看 ./api/readme.md

## package 说明

说明时间：20201209

由于kratos的warden暂未实现grpc v1.34.0的接口，
原因：krotos依赖上游etcd使用v1.29.1，

google.golang.org/grpc v1.29.1
google.golang.org/protobuf v1.25.0

如能使用grpc v1.34.0，需更新protoc-gen-go/protoc-gen-go-grpc
