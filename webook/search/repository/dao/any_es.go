package dao

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type AnyESDAO struct {
	client *elastic.Client
}

func NewAnyESDAO(client *elastic.Client) AnyDAO {
	return &AnyESDAO{client: client}
}

func (a *AnyESDAO) Input(ctx context.Context, index, docId, data string) error {
	_, err := a.client.Index().
		// 直接整个 data 从 Kafka/grpc 里面一路透传到这里
		Index(index).Id(docId).BodyString(data).Do(ctx)
	return err
}
