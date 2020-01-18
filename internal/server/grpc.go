package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	v1 "github.com/vienai8d/grpc-gateway-sample/internal/api/v1"
	"github.com/vienai8d/grpc-gateway-sample/internal/service"
)

type GrpcConfig struct {
	Host string
	Port int
}

func RunGrpc(ctx context.Context, c *GrpcConfig) chan error {
	errChan := make(chan error, 1)
	srv := grpc.NewServer()

	go func() {
		defer close(errChan)
		addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			errChan <- err
			return
		}
		v1.RegisterExampleServer(srv, &service.ExampleServer{})
		errChan <- srv.Serve(lis)
	}()

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
		<-errChan
	}()

	return errChan
}
