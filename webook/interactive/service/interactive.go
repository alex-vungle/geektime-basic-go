package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/interactive/domain"
	"gitee.com/geekbang/basic-go/webook/interactive/events"
	"gitee.com/geekbang/basic-go/webook/interactive/repository"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"golang.org/x/sync/errgroup"
	"time"
)

//go:generate mockgen -source=./interactive.go -package=svcmocks -destination=./mocks/interactive.mock.go InteractiveService
type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(c context.Context, biz string, id int64, uid int64) error
	CancelLike(c context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, bizId, cid, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	repo     repository.InteractiveRepository
	producer events.InteractiveProducer
	l        logger.LoggerV1
}

func (i *interactiveService) GetByIds(ctx context.Context,
	biz string, ids []int64) (map[int64]domain.Interactive, error) {
	intrs, err := i.repo.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Interactive, len(intrs))
	for _, intr := range intrs {
		res[intr.BizId] = intr
	}
	return res, nil
}

func (i *interactiveService) Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	intr, err := i.repo.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	var eg errgroup.Group
	eg.Go(func() error {
		var er error
		intr.Liked, er = i.repo.Liked(ctx, biz, id, uid)
		return er
	})

	eg.Go(func() error {
		var er error
		intr.Collected, er = i.repo.Collected(ctx, biz, id, uid)
		return er
	})
	return intr, eg.Wait()
}

func (i *interactiveService) Collect(ctx context.Context, biz string, bizId, cid, uid int64) error {
	err := i.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
	if err != nil {
		return err
	}
	go func() {
		err = i.CollectSync(uid, biz, bizId)
		if err != nil {
			i.l.Error("同步收藏数失败",
				logger.String("biz", biz),
				logger.Int64("bizId", bizId),
				logger.Error(err))
		}
	}()
	return nil
}

func (i *interactiveService) Like(c context.Context, biz string, id int64, uid int64) error {
	err := i.repo.IncrLike(c, biz, id, uid)
	if err != nil {
		return err
	}
	go func() {
		err = i.likeSync(uid, biz, id)
		if err != nil {
			i.l.Error("同步点赞数失败",
				logger.String("biz", biz),
				logger.Int64("bizId", id),
				logger.Error(err))
		}
	}()
	return nil
}

func (i *interactiveService) CancelLike(c context.Context, biz string, id int64, uid int64) error {
	err := i.repo.DecrLike(c, biz, id, uid)
	if err != nil {
		return err
	}
	go func() {
		err = i.cancelLikeSync(uid, biz, id)
		if err != nil {
			i.l.Error("同步点赞数失败",
				logger.String("biz", biz),
				logger.Int64("bizId", id),
				logger.Error(err))
		}
	}()
	return nil
}

func NewInteractiveService(repo repository.InteractiveRepository, l logger.LoggerV1, producer events.InteractiveProducer) InteractiveService {
	return &interactiveService{
		repo:     repo,
		l:        l,
		producer: producer,
	}
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}

func (i *interactiveService) likeSync(uid int64, biz string, bizId int64) error {
	return i.sync(events.LikeEventType, uid, biz, bizId)
}
func (i *interactiveService) CollectSync(uid int64, biz string, bizId int64) error {
	return i.sync(events.CollectEventType, uid, biz, bizId)
}
func (i *interactiveService) cancelLikeSync(uid int64, biz string, bizId int64) error {
	return i.sync(events.CancelLikeEventType, uid, biz, bizId)
}

func (i *interactiveService) sync(typ events.InteractiveEventType, uid int64, biz string, bizId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return i.producer.ProduceInteractiveEvent(ctx, events.InteractiveEvent{
		Type:  typ,
		Biz:   biz,
		BizId: bizId,
		Uid:   uid,
	})
}
