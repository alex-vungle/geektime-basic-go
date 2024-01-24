//go:build wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/follow/grpc"
	"gitee.com/geekbang/basic-go/webook/follow/repository"
	"gitee.com/geekbang/basic-go/webook/follow/repository/cache"
	"gitee.com/geekbang/basic-go/webook/follow/repository/dao"
	"gitee.com/geekbang/basic-go/webook/follow/service"
	"github.com/google/wire"
)

func InitServer() *grpc.FollowServiceServer {
	wire.Build(
		InitRedis,
		InitLog,
		InitTestDB,
		dao.NewGORMFollowRelationDAO,
		cache.NewRedisFollowCache,
		repository.NewFollowRelationRepository,
		service.NewFollowRelationService,
		grpc.NewFollowRelationServiceServer,
	)
	return new(grpc.FollowServiceServer)
}
