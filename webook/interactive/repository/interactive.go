package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/interactive/domain"
	"gitee.com/geekbang/basic-go/webook/interactive/repository/cache"
	"gitee.com/geekbang/basic-go/webook/interactive/repository/dao"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// BatchIncrReadCnt biz 和 bizId 长度必须一致
	BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logger.LoggerV1
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO,
	l logger.LoggerV1,
	cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao, cache: cache, l: l}
}

func (c *CachedInteractiveRepository) GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	intrs, err := c.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map(intrs, func(idx int, src dao.Interactive) domain.Interactive {
		return c.toDomain(src)
	}), nil
}

func (c *CachedInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, id)
	if err == nil {
		return intr, nil
	}
	ie, err := c.dao.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	if err == nil {
		res := c.toDomain(ie)
		err = c.cache.Set(ctx, biz, id, res)
		if err != nil {
			c.l.Error("回写缓存失败",
				logger.String("biz", biz),
				logger.Int64("bizId", id),
				logger.Error(err))
		}
		return res, nil
	}
	return intr, err
}

func (c *CachedInteractiveRepository) Liked(ctx context.Context,
	biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedInteractiveRepository) Collected(ctx context.Context,
	biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedInteractiveRepository) AddCollectionItem(ctx context.Context,
	biz string, id int64, cid int64, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		BizId: id,
		Cid:   cid,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	// 这两个操作可以考虑合并为一个操作
	err = c.cache.IncrLikeCntIfPresent(ctx, biz, id)
	if err != nil {
		return err
	}
	err = c.cache.IncrRankingIfPresent(ctx, biz, id)
	if err == cache.RankingUpdateErr {
		// 这是一个可选的，跟你的模型有关
		val, err := c.dao.Get(ctx, biz, id)
		if err != nil {
			return err
		}
		return c.cache.SetRankingScore(ctx, biz, id, val.LikeCnt)
	}
	return err
}

func (c *CachedInteractiveRepository) LikeTop(ctx context.Context, biz string) ([]domain.Interactive, error) {
	return c.cache.LikeTop(ctx, biz)
}

func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error {
	err := c.dao.BatchIncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	go func() {
		for i := 0; i < len(biz); i++ {
			er := c.cache.IncrReadCntIfPresent(ctx, biz[i], bizId[i])
			if er != nil {
				// 记录日志
			}
		}
	}()
	return nil
}

func (c *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	// 你要更新缓存了
	// 部分失败问题 —— 数据不一致
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CachedInteractiveRepository) toDomain(ie dao.Interactive) domain.Interactive {
	return domain.Interactive{
		BizId:      ie.BizId,
		ReadCnt:    ie.ReadCnt,
		LikeCnt:    ie.LikeCnt,
		CollectCnt: ie.CollectCnt,
	}
}
