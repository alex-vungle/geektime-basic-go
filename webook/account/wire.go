//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/account/grpc"
	"gitee.com/geekbang/basic-go/webook/account/ioc"
	"gitee.com/geekbang/basic-go/webook/account/repository"
	"gitee.com/geekbang/basic-go/webook/account/repository/dao"
	"gitee.com/geekbang/basic-go/webook/account/service"
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
	"github.com/google/wire"
)

func Init() *wego.App {
	wire.Build(
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitGRPCxServer,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc.NewAccountServiceServer,
		wire.Struct(new(wego.App), "GRPCServer"))
	return new(wego.App)
}
