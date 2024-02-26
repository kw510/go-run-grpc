package main

import (
	"context"

	"github.com/kw510/go-run-grpc/genproto/acme/helloworld/v1"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (server) Hello(ctx context.Context, request *helloworld.GreeterHelloRequest) (response *helloworld.GreeterHelloResponse, err error) {
	response = &helloworld.GreeterHelloResponse{
		Message: "Hello, " + request.Name + "",
	}
	return
}
