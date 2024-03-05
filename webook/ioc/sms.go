package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/auth"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/localsms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/prometheus"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/ratelimit"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/tencent"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"time"
)

func InitSMSService() sms.Service {
	//return ratelimit.NewRateLimitSMSService(localsms.NewService(), limiter.NewRedisSlidingWindowLimiter())
	return localsms.NewService()
	// 如果有需要，就可以用这个
	//return initTencentSMSService()
}

// InitSMSServiceV1 2024.3.5 答疑
func InitSMSServiceV1(cmd redis.Cmdable) sms.Service {
	//return ratelimit.NewRateLimitSMSService(localsms.NewService(), limiter.NewRedisSlidingWindowLimiter())
	// 用的是本地内存的实现，输出控制台
	//return localsms.NewService()
	// 如果有需要，就可以用这个
	// 线上环境，我用腾讯云的实现
	//return initTencentSMSService()
	// 我要搞一个限流的实现，限制住我客户端这边，比如说不超过 300/s
	var res sms.Service = localsms.NewService()
	res = ratelimit.NewRateLimitSMSService(res,
		limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 300))
	// 假如说，我进一步叠加记录 trace 的功能
	//res = opentelemetry.NewDecorator(res, otel.Tracer("demo_qa"))
	//return aliyun.NewSMSService()
	// 再叠
	res = prometheus.NewDecorator(res, prometheus2.SummaryOpts{
		Name: "demo_qa",
	})
	// 我再叠加 auth
	res = auth.NewSMSService(res, []byte("abcd"))
	return res
}

func initTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯 SMS 的 secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯 SMS 的 secret key")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "1400842696", "妙影科技")
}
