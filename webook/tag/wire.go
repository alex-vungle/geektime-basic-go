package main

import (
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
	"gitee.com/geekbang/basic-go/webook/tag/grpc"
	"gitee.com/geekbang/basic-go/webook/tag/ioc"
	"gitee.com/geekbang/basic-go/webook/tag/repository/cache"
	"gitee.com/geekbang/basic-go/webook/tag/repository/dao"
	"gitee.com/geekbang/basic-go/webook/tag/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitDB,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		cache.NewRedisTagCache,
		dao.NewGORMTagDAO,
		ioc.InitRepository,
		service.NewTagService,
		grpc.NewTagServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
