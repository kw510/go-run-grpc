# Go Run gRPC
A very basic starting point for go projects to run within Google Cloud Run containers.

Supports

- HTTP/1.1 h2c upgrade into gRPC service
- Google Pub/Sub Push subscription

Inprogress gateway provided by [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

## Generate Protos
```
 protoc \
  -I ./proto \
  -I ~/path/to/googleapis \
  --go_out=.      --go_opt=module=github.com/kw510/go-run-grpc \
  --go-grpc_out=. --go-grpc_opt=module=github.com/kw510/go-run-grpc \
  --go-mock_out=. --go-mock_opt=module=github.com/kw510/go-run-grpc \
  --grpc-gateway_out=. --grpc-gateway_opt=module=github.com/kw510/go-run-grpc \
  ./proto/**/*.proto
```

## Future

Remove dependency on grpc-gateway (Create a new protoc plugin to use raw server?)
