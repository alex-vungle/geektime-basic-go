package events

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"gitee.com/geekbang/basic-go/webook/search/service"
	"github.com/IBM/sarama"
	"log"
	"time"
)

const InteractiveTopic = "sync_interactive"

type InteractiveConsumer struct {
	syncSvc  service.SyncService
	client   sarama.Client
	l        logger.LoggerV1
	handlers map[int64]InteractiveHandler
}

func NewInteractiveConsumer(client sarama.Client,
	l logger.LoggerV1,
	svc service.SyncService) *InteractiveConsumer {
	handlers := map[int64]InteractiveHandler{
		1: &LikeHandler{syncSvc: svc},
		2: &CollectHandler{syncSvc: svc},
		3: &CancelLikeHandler{syncSvc: svc},
	}
	return &InteractiveConsumer{
		syncSvc:  svc,
		client:   client,
		l:        l,
		handlers: handlers,
	}
}

type InteractiveEvent struct {
	Type  int64  `json:"type,omitempty"`
	Uid   int64  `json:"uid"`
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
}

func (a *InteractiveConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_interactive_group1",
		a.client)
	if err != nil {
		return err
	}
	go func() {
		log.Println("开启消费进程")
		err := cg.Consume(context.Background(),
			[]string{InteractiveTopic},
			saramax.NewHandler[InteractiveEvent](a.l, a.Consume))
		if err != nil {
			a.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (a *InteractiveConsumer) Consume(sg *sarama.ConsumerMessage,
	evt InteractiveEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100000*time.Second)
	defer cancel()
	log.Println(fmt.Sprintf("开始消费%v", evt))
	handleFunc, ok := a.handlers[evt.Type]
	if !ok {
		a.l.Error("未知类型")
		return nil
	}
	return handleFunc.Handle(ctx, evt)
}

type InteractiveHandler interface {
	Handle(ctx context.Context, data InteractiveEvent) error
}

type LikeHandler struct {
	syncSvc service.SyncService
}

func (l *LikeHandler) Handle(ctx context.Context, data InteractiveEvent) error {
	return handle(ctx, l.syncSvc, "like_index", getDocId(data), data)
}

type CollectHandler struct {
	syncSvc service.SyncService
}

func (c *CollectHandler) Handle(ctx context.Context, data InteractiveEvent) error {
	return handle(ctx, c.syncSvc, "collect_index", getDocId(data), data)
}

type CancelLikeHandler struct {
	syncSvc service.SyncService
}

func (c *CancelLikeHandler) Handle(ctx context.Context, data InteractiveEvent) error {
	return c.syncSvc.Delete(ctx, "like_index", getDocId(data))
}

func getDocId(data InteractiveEvent) string {
	return fmt.Sprintf("%d_%s_%d", data.Uid, data.Biz, data.BizId)
}

func handle(ctx context.Context, syncSvc service.SyncService, index, key string, data any) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return syncSvc.InputAny(ctx, index, key, string(val))
}
