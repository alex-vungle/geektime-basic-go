// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/code/grpc"
	"gitee.com/geekbang/basic-go/webook/internal/code/ioc"
	"gitee.com/geekbang/basic-go/webook/internal/code/repository"
	"gitee.com/geekbang/basic-go/webook/internal/code/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/code/service"
	"github.com/google/wire"
)

// Injectors from wire.go:

func Init() *App {
	smsServiceClient := ioc.InitSmsRpcClient()
	cmdable := ioc.InitRedis()
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	codeService := service.NewSMSCodeService(smsServiceClient, codeRepository)
	codeServiceServer := grpc.NewCodeServiceServer(codeService)
	server := ioc.InitGRPCxServer(codeServiceServer)
	app := &App{
		server: server,
	}
	return app
}

// wire.go:

var thirdProvider = wire.NewSet(ioc.InitRedis, ioc.InitSmsRpcClient)
