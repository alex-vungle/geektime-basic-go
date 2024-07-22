package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "demo", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = consumer.Consume(ctx, []string{"test_topic"}, ConsumerHandler{})
	assert.NoError(t, err)
}

type ConsumerHandler struct{}

func (c ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("这是Setup")
	var offset int64 = -1
	partitions := session.Claims()["test_topic"]
	for _, p := range partitions {
		session.ResetOffset("test_topic", p, offset, "")
	}
	return nil
}

func (c ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("这是Cleanup")
	return nil
}

func (c ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Println("这是ConsumeClaim")
	msgs := claim.Messages()
	const batchSize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, 10)
		var eg errgroup.Group
		for i := 0; i < 10; i++ {
			msg := <-msgs
			batch = append(batch, msg)
			eg.Go(func() error {
				log.Println(string(msg.Value))
				return nil
			})
		}
		err := eg.Wait()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

func (c ConsumerHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Println("这是ConsumeClaim")
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println(string(msg.Value))
	}
	return nil
}
