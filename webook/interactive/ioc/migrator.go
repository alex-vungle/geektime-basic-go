package ioc

import (
	"gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx"
	"gitee.com/geekbang/basic-go/webook/pkg/gormx/connpool"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator/events"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator/events/fixer"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator/scheduler"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

const topic = "migrator_interactives"

func InitFixDataConsumer(l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	client sarama.Client) *fixer.Consumer[dao.Interactive] {
	res, err := fixer.NewConsumer[dao.Interactive](client, l,
		topic, src, dst)
	if err != nil {
		panic(err)
	}
	return res
}

func InitMigradatorProducer(p sarama.SyncProducer) events.Producer {
	return events.NewSaramaProducer(p, topic)
}

func InitMigratorWeb(
	l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	pool *connpool.DoubleWritePool,
	producer events.Producer,
) *ginx.Server {
	// 在这里，有多少张表，你就初始化多少个 scheduler
	intrSch := scheduler.NewScheduler[dao.Interactive](l, src, dst, pool, producer)
	engine := gin.Default()
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook_intr_admin",
		Name:      "http_biz_code",
		Help:      "HTTP 的业务错误码",
	})
	intrSch.RegisterRoutes(engine.Group("/migrator"))
	//intrSch.RegisterRoutes(engine.Group("/migrator/interactive"))
	addr := viper.GetString("migrator.web.addr")
	return &ginx.Server{
		Addr:   addr,
		Engine: engine,
	}
}
