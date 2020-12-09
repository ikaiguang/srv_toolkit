# protobuf

protoc-gen-go google.golang.org/protobuf@v1.25.0

```shell script

# windows
protoc -I. -I%GOSRCPATH% -I%GOPBPATH% --go_out=. --go_opt=paths=source_relative ./*.proto

# linux
protoc -I. -I$GOSRCPATH -I$GOPBPATH --go_out=. --go_opt=paths=source_relative ./*.proto

```