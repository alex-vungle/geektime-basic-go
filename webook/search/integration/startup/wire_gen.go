// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/search/grpc"
	"gitee.com/geekbang/basic-go/webook/search/ioc"
	"gitee.com/geekbang/basic-go/webook/search/repository"
	"gitee.com/geekbang/basic-go/webook/search/repository/dao"
	"gitee.com/geekbang/basic-go/webook/search/service"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitSearchServer() *grpc.SearchServiceServer {
	client := InitESClient()
	userDAO := dao.NewUserElasticDAO(client)
	userRepository := repository.NewUserRepository(userDAO)
	articleDAO := dao.NewArticleElasticDAO(client)
	articleRepository := repository.NewArticleRepository(articleDAO)
	searchService := service.NewSearchService(userRepository, articleRepository)
	searchServiceServer := grpc.NewSearchService(searchService)
	return searchServiceServer
}

func InitSyncServer() *grpc.SyncServiceServer {
	client := InitESClient()
	userDAO := dao.NewUserElasticDAO(client)
	userRepository := repository.NewUserRepository(userDAO)
	articleDAO := dao.NewArticleElasticDAO(client)
	articleRepository := repository.NewArticleRepository(articleDAO)
	syncService := service.NewSyncService(userRepository, articleRepository)
	syncServiceServer := grpc.NewSyncServiceServer(syncService)
	return syncServiceServer
}

// wire.go:

var serviceProviderSet = wire.NewSet(dao.NewUserElasticDAO, dao.NewArticleElasticDAO, repository.NewUserRepository, repository.NewArticleRepository, service.NewSyncService, service.NewSearchService)

var thirdProvider = wire.NewSet(
	InitESClient, ioc.InitLogger,
)
