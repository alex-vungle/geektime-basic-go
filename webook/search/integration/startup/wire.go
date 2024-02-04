//go:build wireinject

package startup

import (
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
	dao.NewTagESDAO,
	dao.NewAnyESDAO,
	repository.NewUserRepository,
	repository.NewAnyRepository,
	repository.NewArticleRepository,
	service.NewSyncService,
	service.NewSearchService,
)

var thirdProvider = wire.NewSet(
	InitESClient,
	ioc.InitLogger)

func InitSearchServer() *grpc.SearchServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSearchService,
	)
	return new(grpc.SearchServiceServer)
}

func InitSyncServer() *grpc.SyncServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSyncServiceServer,
	)
	return new(grpc.SyncServiceServer)
}
