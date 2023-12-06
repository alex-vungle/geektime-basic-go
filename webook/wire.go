//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/interactive/events"
	repository2 "gitee.com/geekbang/basic-go/webook/interactive/repository"
	cache2 "gitee.com/geekbang/basic-go/webook/interactive/repository/cache"
	dao2 "gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	service2 "gitee.com/geekbang/basic-go/webook/interactive/service"
	"gitee.com/geekbang/basic-go/webook/internal/events/article"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	article2 "gitee.com/geekbang/basic-go/webook/internal/repository/article"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	article3 "gitee.com/geekbang/basic-go/webook/internal/repository/dao/article"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	ijwt "gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/google/wire"
)

var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewCachedInteractiveRepository,
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)

var rankingServiceSet = wire.NewSet(
	repository.NewCachedRankingRepository,
	cache.NewRankingRedisCache,
	cache.NewRankingLocalCache,
	service.NewBatchRankingService,
)

func InitWebServer() *App {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis, ioc.InitRLockClient,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		interactiveSvcProvider,
		ioc.InitIntrGRPCClient,
		rankingServiceSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// consumer
		events.NewInteractiveReadEventBatchConsumer,
		article.NewKafkaProducer,

		// 初始化 DAO
		dao.NewUserDAO,
		article3.NewGORMArticleDAO,

		cache.NewUserCache,
		cache.NewCodeCache,
		cache.NewRedisArticleCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		article2.NewArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// 直接基于内存实现
		ioc.InitSMSService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		//ioc.NewWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,
		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		//gin.Default,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
		// 组装我这个结构体的所有字段
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
