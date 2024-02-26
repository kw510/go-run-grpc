package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	helloworld "github.com/kw510/go-run-grpc/genproto/acme/helloworld/v1"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	// Main running context during app lifecycle
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create gRPC server
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	// Register service
	helloworld.RegisterGreeterServer(grpcServer, server{})

	// In progress proxy
	// Register services to gateway
	mux := runtime.NewServeMux()
	err = helloworld.RegisterGreeterHandlerFromEndpoint(ctx, mux, "localhost:8080", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		return
	}

	// Create mux main http handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	})

	// Create h2c http server
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      h2c.NewHandler(handler, &http2.Server{}),
	}

	// Serve
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	log.Println("server initialised")

	select {
	case err = <-srvErr:
		// Error when starting HTTP server
		return
	case <-ctx.Done():
		// Wait for first CTRL+C
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	err = srv.Shutdown(context.Background())

	return
}
