//go:build wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/account/grpc"
	"gitee.com/geekbang/basic-go/webook/account/repository"
	"gitee.com/geekbang/basic-go/webook/account/repository/dao"
	"gitee.com/geekbang/basic-go/webook/account/service"
	"github.com/google/wire"
)

func InitAccountService() *grpc.AccountServiceServer {
	wire.Build(InitTestDB,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc.NewAccountServiceServer)
	return new(grpc.AccountServiceServer)
}
