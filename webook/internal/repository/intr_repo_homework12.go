package repository

import (
	"context"
	intrv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/intr/v1"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
)

type InteractiveRepository interface {
	// GetById 获得点赞
	GetById(ctx context.Context, uid, id int64) (domain.Interactive, error)
}

type GRPCInteractiveRepository struct {
	client intrv1.InteractiveServiceClient
}

func NewGRPCInteractiveRepository(client intrv1.InteractiveServiceClient) *GRPCInteractiveRepository {
	return &GRPCInteractiveRepository{client: client}
}

func (repo *GRPCInteractiveRepository) GetById(ctx context.Context, uid, id int64) (domain.Interactive, error) {
	resp, err := repo.client.Get(ctx, &intrv1.GetRequest{
		Biz:   "article",
		BizId: id,
		Uid:   uid,
	})
	if err != nil {
		return domain.Interactive{}, err
	}
	return domain.Interactive{
		Biz:        "article",
		BizId:      id,
		ReadCnt:    resp.Intr.ReadCnt,
		LikeCnt:    resp.Intr.LikeCnt,
		CollectCnt: resp.Intr.CollectCnt,
		Liked:      resp.Intr.Liked,
		Collected:  resp.Intr.Collected,
	}, nil
}
