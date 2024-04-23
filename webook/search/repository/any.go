package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/search/repository/dao"
)

type AnyRepository interface {
	Input(ctx context.Context, index string, docID string, data string) error
	Delete(ctx context.Context, index, docID string) error
}

type anyRepository struct {
	dao dao.AnyDAO
}

func (repo *anyRepository) Delete(ctx context.Context, index string, docID string) error {
	return repo.dao.Delete(ctx, index, docID)
}

func NewAnyRepository(dao dao.AnyDAO) AnyRepository {
	return &anyRepository{dao: dao}
}

func (repo *anyRepository) Input(ctx context.Context, index string, docID string, data string) error {
	return repo.dao.Input(ctx, index, docID, data)
}
