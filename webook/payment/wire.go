//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/payment/ioc"
	"gitee.com/geekbang/basic-go/webook/payment/web"
	"gitee.com/geekbang/basic-go/webook/pkg/wego"
	"github.com/google/wire"
)

func InitApp() *wego.App {
	wire.Build(
		//ioc.InitWechatClient,
		//ioc.InitWechatNativeService,
		ioc.InitWechatConfig,
		ioc.InitWechatNotifyHandler,
		web.NewWechatHandler,
		ioc.InitGinServer,
		ioc.InitLogger,
		wire.Struct(new(wego.App), "WebServer"))
	return new(wego.App)
}
