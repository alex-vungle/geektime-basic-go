//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/feed/events"
	"gitee.com/geekbang/basic-go/webook/feed/grpc"
	"gitee.com/geekbang/basic-go/webook/feed/ioc"
	"gitee.com/geekbang/basic-go/webook/feed/repository"
	"gitee.com/geekbang/basic-go/webook/feed/repository/cache"
	"gitee.com/geekbang/basic-go/webook/feed/repository/dao"
	"gitee.com/geekbang/basic-go/webook/feed/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewFeedPushEventDAO,
	dao.NewFeedPullEventDAO,
	cache.NewFeedEventCache,
	repository.NewFeedEventRepo,
)

var thirdProvider = wire.NewSet(
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitKafka,
	ioc.InitDB,
	ioc.InitFollowClient,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.RegisterHandler,
		service.NewFeedService,
		grpc.NewFeedEventGrpcSvc,
		events.NewArticleEventConsumer,
		events.NewFeedEventConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
