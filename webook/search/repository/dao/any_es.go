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
		Index(index).Id(docId).BodyString(data).Do(ctx)
	return err
}

func (a *AnyESDAO) Delete(ctx context.Context, index string, docID string) error {
	_, err := a.client.Delete().Index(index).Id(docID).Do(ctx)
	return err
}
