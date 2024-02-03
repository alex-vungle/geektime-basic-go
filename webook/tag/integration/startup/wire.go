//go:build wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/tag/grpc"
	"gitee.com/geekbang/basic-go/webook/tag/repository/cache"
	"gitee.com/geekbang/basic-go/webook/tag/repository/dao"
	"gitee.com/geekbang/basic-go/webook/tag/service"
	"github.com/google/wire"
)

func InitGRPCService() *grpc.TagServiceServer {
	wire.Build(InitTestDB, InitRedis,
		InitLog,
		dao.NewGORMTagDAO,
		InitRepository,
		cache.NewRedisTagCache,
		service.NewTagService,
		grpc.NewTagServiceServer,
	)
	return new(grpc.TagServiceServer)
}
