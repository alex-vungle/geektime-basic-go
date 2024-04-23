package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const interactiveTopic = "interactive_sync"

type InteractiveProducer interface {
	ProduceInteractiveEvent(ctx context.Context, evt InteractiveEvent) error
}

type InteractiveEvent struct {
	Type  InteractiveEventType `json:"type"` // 1-喜欢 2-收藏 3-取消喜欢
	Biz   string               `json:"biz"`
	BizId int64                `json:"bizId"`
	Uid   int64                `json:"uid"`
}

type InteractiveEventType int64

const (
	LikeEventType       InteractiveEventType = 1
	CollectEventType    InteractiveEventType = 2
	CancelLikeEventType InteractiveEventType = 3
)

type interactiveProducer struct {
	producer sarama.SyncProducer
}

func NewInteractiveProducer(p sarama.SyncProducer) InteractiveProducer {
	return &interactiveProducer{
		producer: p,
	}
}

func (i *interactiveProducer) ProduceInteractiveEvent(ctx context.Context, interactive InteractiveEvent) error {
	data, err := json.Marshal(interactive)
	if err != nil {
		return err
	}
	_, _, err = i.producer.SendMessage(&sarama.ProducerMessage{
		Topic: interactiveTopic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
