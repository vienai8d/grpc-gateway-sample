package service

import (
	"context"

	"github.com/golang/glog"

	v1 "github.com/vienai8d/grpc-gateway-sample/internal/api/v1"
)

type ExampleServer struct {
}

func (s *ExampleServer) Echo(ctx context.Context, in *v1.EchoRequest) (*v1.EchoResponse, error) {
	glog.Info(in)
	return &v1.EchoResponse{Text: in.GetText()}, nil
}
