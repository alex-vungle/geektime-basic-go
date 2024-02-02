//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/search/events"
	"gitee.com/geekbang/basic-go/webook/search/grpc"
	"gitee.com/geekbang/basic-go/webook/search/ioc"
	"gitee.com/geekbang/basic-go/webook/search/repository"
	"gitee.com/geekbang/basic-go/webook/search/repository/dao"
	"gitee.com/geekbang/basic-go/webook/search/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewUserElasticDAO,
	dao.NewArticleElasticDAO,
	dao.NewAnyESDAO,
	dao.NewTagESDAO,
	repository.NewUserRepository,
	repository.NewArticleRepository,
	repository.NewAnyRepository,
	service.NewSyncService,
	service.NewSearchService,
)

var thirdProvider = wire.NewSet(
	ioc.InitESClient,
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitKafka)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSyncServiceServer,
		grpc.NewSearchService,
		events.NewUserConsumer,
		events.NewArticleConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
