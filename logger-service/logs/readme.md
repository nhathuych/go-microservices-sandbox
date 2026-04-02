## gRPC stuff

Install binaries:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6
```

Generate proto.

Then, run command (inside logs directory):

```bash
protoc --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  logs.proto 
```
