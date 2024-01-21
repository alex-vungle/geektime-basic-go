//go:build wireinject

package startup

import (
	pmtv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/payment/v1"
	"gitee.com/geekbang/basic-go/webook/reward/repository"
	"gitee.com/geekbang/basic-go/webook/reward/repository/cache"
	"gitee.com/geekbang/basic-go/webook/reward/repository/dao"
	"gitee.com/geekbang/basic-go/webook/reward/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(InitTestDB, InitLogger, InitRedis)

func InitWechatNativeSvc(client pmtv1.WechatPaymentServiceClient) *service.WechatNativeRewardService {
	wire.Build(service.NewWechatNativeRewardService,
		thirdPartySet,
		cache.NewRewardRedisCache,
		repository.NewRewardRepository, dao.NewRewardGORMDAO)
	return new(service.WechatNativeRewardService)
}
