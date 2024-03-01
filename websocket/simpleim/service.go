package simpleim

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"strconv"
)

type IMService struct {
	producer sarama.SyncProducer
}

func (s *IMService) Receive(ctx context.Context, sender int64, msg Message) error {
	// 转发到 Kafka 里面，以通知别的网关节点
	// 1. 查找目标
	members := s.findMembers()
	// 2. 通知 Kafka，让别的节点能够订阅到消息
	for _, mem := range members {
		if mem == sender {
			// 本人就不用转发了
			continue
		}
		msgJson, err := json.Marshal(Event{Receiver: mem, Msg: msg})
		if err != nil {
			continue
		}
		_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
			Topic: eventName,
			Key:   sarama.ByteEncoder(strconv.FormatInt(mem, 10)),
			Value: sarama.ByteEncoder(msgJson),
		})
		if err != nil {
			log.Println("发送消息失败", err)
			continue
		}
	}
	return nil
}

// 这里模拟根据 cid，也就是聊天 ID 来查找参与了该聊天的成员
func (s *IMService) findMembers() []int64 {
	// 固定返回 1，2，3，4
	return []int64{1, 2, 3, 4}
}
