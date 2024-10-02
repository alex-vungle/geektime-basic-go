package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
)

const LikeIndexName = "like_index"

type likeElasticDAO struct {
	client *elastic.Client
}

func NewLikeDAO(client *elastic.Client) LikeDAO {
	return &likeElasticDAO{
		client: client,
	}
}

func (l *likeElasticDAO) Search(ctx context.Context, uid int64, biz string) ([]int64, error) {
	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("uid", uid),
		elastic.NewTermQuery("biz", biz),
	)
	resp, err := l.client.Search(LikeIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]int64, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var bt BizTags
		err = json.Unmarshal(hit.Source, &bt)
		if err != nil {
			return nil, err
		}
		res = append(res, bt.BizId)
	}
	return res, nil

}
