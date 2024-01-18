//go:build wireinject

package startup

import (
	"gitee.com/geekbang/basic-go/webook/payment/ioc"
	"gitee.com/geekbang/basic-go/webook/payment/repository"
	"gitee.com/geekbang/basic-go/webook/payment/repository/dao"
	"gitee.com/geekbang/basic-go/webook/payment/service/wechat"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(ioc.InitLogger, InitTestDB)

var wechatNativeSvcSet = wire.NewSet(
	ioc.InitWechatClient,
	dao.NewPaymentGORMDAO,
	repository.NewPaymentRepository,
	ioc.InitWechatNativeService,
	ioc.InitWechatConfig)

func InitWechatNativeService() *wechat.NativePaymentService {
	wire.Build(wechatNativeSvcSet, thirdPartySet)
	return new(wechat.NativePaymentService)
}
