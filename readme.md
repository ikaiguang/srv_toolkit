# srv_toolkit

toolkit for srv

## protobuf

```shell script
cd ./api

# windows
protoc -I. -I%GOSRCPATH% -I%GOPBPATH% --go_out=. --go_opt=paths=source_relative ./*.proto

# linux
protoc -I. -I$GOSRCPATH -I$GOPBPATH --go_out=. --go_opt=paths=source_relative ./*.proto

```