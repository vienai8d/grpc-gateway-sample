package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	v1 "github.com/vienai8d/grpc-gateway-sample/internal/api/v1"
)

type HTTPConfig struct {
	Host     string
	Port     int
	GrpcPort int
}

func RunHTTP(ctx context.Context, c *HTTPConfig) chan error {
	errChan := make(chan error, 1)
	srv := http.Server{}

	go func() {
		defer close(errChan)
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		httpEndpoint := fmt.Sprintf("%s:%d", c.Host, c.Port)
		grpcEndpoint := fmt.Sprintf("%s:%d", c.Host, c.GrpcPort)
		err := v1.RegisterExampleHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		if err != nil {
			errChan <- err
			return
		}
		srv.Addr = httpEndpoint
		srv.Handler = mux
		errChan <- srv.ListenAndServe()
	}()

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
		<-errChan
	}()

	return errChan
}
