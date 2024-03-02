package simpleim

import (
	"context"
	"github.com/IBM/sarama"
)

type IMService struct {
	producer sarama.SyncProducer
}

func (s *IMService) Receive(ctx context.Context, sender int64, msg Message) error {

}

// 这里模拟根据 cid，也就是聊天 ID 来查找参与了该聊天的成员
func (s *IMService) findMembers() []int64 {
	// 固定返回 1，2，3，4
	return []int64{1, 2, 3, 4}
}
