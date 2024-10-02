package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
)

const CollectIndexName = "collect_index"

type collectElasticDAO struct {
	client *elastic.Client
}

func NewCollectDAO(client *elastic.Client) CollectDAO {
	return &collectElasticDAO{
		client: client,
	}
}

func (c *collectElasticDAO) Search(ctx context.Context, uid int64, biz string) ([]int64, error) {
	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("uid", uid),
		elastic.NewTermQuery("biz", biz),
	)
	resp, err := c.client.Search(CollectIndexName).Query(query).Do(ctx)
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
