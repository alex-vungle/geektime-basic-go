package events

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"gitee.com/geekbang/basic-go/webook/search/service"
	"github.com/IBM/sarama"
	"time"
)

// 通用的 sync data event
// 所有的业务方都可以通过这个 event 来同步数据
type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
}

type SyncDataEventConsumer struct {
	svc    service.SyncService
	client sarama.Client
	l      logger.LoggerV1
}

func (a *SyncDataEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("search_sync_data",
		a.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"sync_search_data"},
			saramax.NewHandler[SyncDataEvent](a.l, a.Consume))
		if err != nil {
			a.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (a *SyncDataEventConsumer) Consume(sg *sarama.ConsumerMessage,
	evt SyncDataEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return a.svc.InputAny(ctx, evt.IndexName, evt.DocID, evt.Data)
}
