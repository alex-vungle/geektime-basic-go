package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/tag/domain"
	"gitee.com/geekbang/basic-go/webook/tag/repository/cache"
	"gitee.com/geekbang/basic-go/webook/tag/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type TagRepository interface {
	CreateTag(ctx context.Context, tag domain.Tag) (int64, error)
	BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error)
	GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error)
}

type CachedTagRepository struct {
	dao   dao.TagDAO
	cache cache.TagCache
	l     logger.LoggerV1
}

// PreloadUserTags 在 toB 的场景下，你可以提前预加载缓存
func (repo *CachedTagRepository) PreloadUserTags(ctx context.Context) error {
	// 我们要存的是 uid => 我的所有标签 [{tag1}, {}]
	// 你这边分批次预加载
	// 数据里面取出来，调用append
	offset := 0
	batch := 100
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		tags, err := repo.dao.GetTags(dbCtx, offset, batch)
		cancel()
		if err != nil {
			// 你也可以 continue
			return err
		}
		for _, tag := range tags {
			rctx, cancel := context.WithTimeout(ctx, time.Second)
			err = repo.cache.Append(rctx, tag.Uid, repo.toDomain(tag))
			cancel()
			if err != nil {
				continue
			}
		}
		if len(tags) < batch {
			return nil
		}
		offset = offset + batch
	}
}

func (repo *CachedTagRepository) GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error) {
	tags, err := repo.dao.GetTagsById(ctx, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	}), nil
}

func (repo *CachedTagRepository) BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	return repo.dao.CreateTagBiz(ctx, slice.Map(tags, func(idx int, src int64) dao.TagBiz {
		return dao.TagBiz{
			Tid:   src,
			BizId: bizId,
			Biz:   biz,
			Uid:   uid,
		}
	}))
}

func (repo *CachedTagRepository) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	res, err := repo.cache.GetTags(ctx, uid)
	if err == nil {
		return res, nil
	}
	tags, err := repo.dao.GetTagsByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	res = slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	})
	err = repo.cache.Append(ctx, uid, res...)
	if err != nil {
		// 记录日志
	}
	return res, nil
}

func (repo *CachedTagRepository) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	// 这里要不要缓存？
	tags, err := repo.dao.GetTagsByBiz(ctx, uid, biz, bizId)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	}), nil
}

func (repo *CachedTagRepository) CreateTag(ctx context.Context, tag domain.Tag) (int64, error) {
	id, err := repo.dao.CreateTag(ctx, repo.toEntity(tag))
	if err != nil {
		return 0, err
	}
	err = repo.cache.Append(ctx, tag.Uid, tag)
	if err != nil {
		// 记录日志就可以
	}
	return id, nil
}

func (repo *CachedTagRepository) toDomain(tag dao.Tag) domain.Tag {
	return domain.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}

func (repo *CachedTagRepository) toEntity(tag domain.Tag) dao.Tag {
	return dao.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}

func NewTagRepository(tagDAO dao.TagDAO, c cache.TagCache, l logger.LoggerV1) *CachedTagRepository {
	return &CachedTagRepository{
		dao:   tagDAO,
		l:     l,
		cache: c,
	}
}
