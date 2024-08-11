package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	cache cache.RankingCache
}

func (repo *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return repo.cache.Get(ctx)
}

func NewCachedRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CachedRankingRepository{cache: cache}
}

func (repo *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return repo.cache.Set(ctx, arts)
}
