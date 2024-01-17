//go:build wireinject

package main

import (
	grpc2 "gitee.com/geekbang/basic-go/webook/comment/grpc"
	"gitee.com/geekbang/basic-go/webook/comment/ioc"
	"gitee.com/geekbang/basic-go/webook/comment/repository"
	"gitee.com/geekbang/basic-go/webook/comment/repository/dao"
	"gitee.com/geekbang/basic-go/webook/comment/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc2.NewGrpcServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
