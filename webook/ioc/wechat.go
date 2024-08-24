package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	appID := "fuck"
	appSecret := "fuck"
	return wechat.NewService(appID, appSecret, l)
}
