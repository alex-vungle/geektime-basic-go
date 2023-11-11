// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	article2 "gitee.com/geekbang/basic-go/webook/internal/events/article"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao/article"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/async"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from wire.go:

//go:generate wire
func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	handler := jwt.NewRedisHandler(cmdable)
	loggerV1 := InitLog()
	v := ioc.GinMiddlewares(cmdable, handler, loggerV1)
	gormDB := InitTestDB()
	userDAO := dao.NewGORMUserDAO(gormDB)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	smsService := ioc.InitSmsMemoryService()
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	codeService := service.NewSMSCodeService(smsService, codeRepository)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	articleDAO := article.NewGORMArticleDAO(gormDB)
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := repository.NewArticleRepository(articleDAO, articleCache, userRepository, loggerV1)
	client := InitKafka()
	syncProducer := NewSyncProducer(client)
	producer := article2.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, loggerV1, producer)
	interactiveDAO := dao.NewGORMInteractiveDAO(gormDB)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	articleHandler := web.NewArticleHandler(articleService, interactiveService, loggerV1)
	observabilityHandler := web.NewObservabilityHandler()
	wechatService := InitPhantomWechatService(loggerV1)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	engine := ioc.InitWebServer(v, userHandler, articleHandler, observabilityHandler, oAuth2WechatHandler, loggerV1)
	return engine
}

func InitArticleHandler(dao2 article.ArticleDAO) *web.ArticleHandler {
	cmdable := InitRedis()
	articleCache := cache.NewRedisArticleCache(cmdable)
	gormDB := InitTestDB()
	userDAO := dao.NewGORMUserDAO(gormDB)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	loggerV1 := InitLog()
	articleRepository := repository.NewArticleRepository(dao2, articleCache, userRepository, loggerV1)
	client := InitKafka()
	syncProducer := NewSyncProducer(client)
	producer := article2.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, loggerV1, producer)
	interactiveDAO := dao.NewGORMInteractiveDAO(gormDB)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	articleHandler := web.NewArticleHandler(articleService, interactiveService, loggerV1)
	return articleHandler
}

func InitUserSvc() service.UserService {
	gormDB := InitTestDB()
	userDAO := dao.NewGORMUserDAO(gormDB)
	cmdable := InitRedis()
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	return userService
}

func InitAsyncSmsService(svc sms.Service) *async.Service {
	gormDB := InitTestDB()
	asyncSmsDAO := dao.NewGORMAsyncSmsDAO(gormDB)
	asyncSmsRepository := repository.NewAsyncSMSRepository(asyncSmsDAO)
	loggerV1 := InitLog()
	asyncService := async.NewService(svc, asyncSmsRepository, loggerV1)
	return asyncService
}

func InitRankingService() service.RankingService {
	gormDB := InitTestDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(gormDB)
	cmdable := InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	loggerV1 := InitLog()
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	articleDAO := article.NewGORMArticleDAO(gormDB)
	articleCache := cache.NewRedisArticleCache(cmdable)
	userRepository := _wireCachedUserRepositoryValue
	articleRepository := repository.NewArticleRepository(articleDAO, articleCache, userRepository, loggerV1)
	client := InitKafka()
	syncProducer := NewSyncProducer(client)
	producer := article2.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, loggerV1, producer)
	redisRankingCache := cache.NewRedisRankingCache(cmdable)
	rankingLocalCache := cache.NewRankingLocalCache()
	rankingRepository := repository.NewCachedRankingRepository(redisRankingCache, rankingLocalCache)
	rankingService := service.NewBatchRankingService(interactiveService, articleService, rankingRepository)
	return rankingService
}

var (
	_wireCachedUserRepositoryValue = &repository.CachedUserRepository{}
)

func InitInteractiveService() service.InteractiveService {
	gormDB := InitTestDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(gormDB)
	cmdable := InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	loggerV1 := InitLog()
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	return interactiveService
}

func InitJwtHdl() jwt.Handler {
	cmdable := InitRedis()
	handler := jwt.NewRedisHandler(cmdable)
	return handler
}

// wire.go:

var thirdProvider = wire.NewSet(InitRedis, InitTestDB,
	InitLog,
	NewSyncProducer,
	InitKafka,
)

var userSvcProvider = wire.NewSet(dao.NewGORMUserDAO, cache.NewRedisUserCache, repository.NewCachedUserRepository, service.NewUserService)

var articlSvcProvider = wire.NewSet(article.NewGORMArticleDAO, article2.NewSaramaSyncProducer, cache.NewRedisArticleCache, repository.NewArticleRepository, service.NewArticleService)

var interactiveSvcProvider = wire.NewSet(service.NewInteractiveService, repository.NewCachedInteractiveRepository, dao.NewGORMInteractiveDAO, cache.NewRedisInteractiveCache)

var rankServiceProvider = wire.NewSet(service.NewBatchRankingService, repository.NewCachedRankingRepository, cache.NewRedisRankingCache, cache.NewRankingLocalCache)
