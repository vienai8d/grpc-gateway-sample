package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"

	"github.com/vienai8d/grpc-gateway-sample/internal/server"
)

func main() {
	defer glog.Flush()

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	gErrChan := server.RunGrpc(ctx, &server.GrpcConfig{
		Port: 9999,
	})
	hErrChan := server.RunHTTP(ctx, &server.HTTPConfig{
		Port:     8888,
		GrpcPort: 9999,
	})

	select {
	case sig := <-sigChan:
		glog.Info("sigChan")
		glog.Info(sig)
	case err := <-gErrChan:
		glog.Info("gErrChan")
		glog.Info(err)
	case err := <-hErrChan:
		glog.Info("hErrChan")
		glog.Info(err)
	}

	cancel()
	<-hErrChan
	<-gErrChan
	glog.Info("done")
}
