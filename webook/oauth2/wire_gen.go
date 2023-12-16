// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire.go/cmd/wire.go
//go:build !wireinject
// +build !wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/oauth2/grpc"
	ioc2 "gitee.com/geekbang/basic-go/webook/internal/oauth2/ioc"
	"github.com/google/wire"
)

// Injectors from wire.go.go:

func Init() *App {
	loggerV1 := ioc2.InitLogger()
	service := ioc2.InitPrometheus(loggerV1)
	oauth2ServiceServer := grpc.NewOauth2ServiceServer(service)
	server := ioc2.InitGRPCxServer(oauth2ServiceServer)
	app := &App{
		server: server,
	}
	return app
}

// wire.go.go:

var thirdProvider = wire.NewSet(ioc2.InitLogger)
