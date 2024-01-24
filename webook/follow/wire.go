//go:build wireinject

package main

import (
	grpc2 "gitee.com/geekbang/basic-go/webook/follow/grpc"
	"gitee.com/geekbang/basic-go/webook/follow/ioc"
	"gitee.com/geekbang/basic-go/webook/follow/repository"
	"gitee.com/geekbang/basic-go/webook/follow/repository/dao"
	"gitee.com/geekbang/basic-go/webook/follow/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMFollowRelationDAO,
	repository.NewFollowRelationRepository,
	service.NewFollowRelationService,
	grpc2.NewFollowRelationServiceServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
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
