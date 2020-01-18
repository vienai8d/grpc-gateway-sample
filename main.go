package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"

	"github.com/vienai8d/grpc-gateway-sample/internal/server"
)

type Env struct {
	Port     int `default:"8888"`
	GrpcPort int `split_words:"true" default:"9999"`
}

func main() {
	defer glog.Flush()

	flag.Parse()

	e := Env{}
	err := envconfig.Process("", &e)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info(e)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	gErrChan := server.RunGrpc(ctx, &server.GrpcConfig{
		Port: e.GrpcPort,
	})
	hErrChan := server.RunHTTP(ctx, &server.HTTPConfig{
		Port:     e.Port,
		GrpcPort: e.GrpcPort,
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
