package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/tag/domain"
	"gitee.com/geekbang/basic-go/webook/tag/events"
	"gitee.com/geekbang/basic-go/webook/tag/repository"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type TagService interface {
	CreateTag(ctx context.Context, uid int64, name string) (int64, error)
	AttachTags(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error)
}

type tagService struct {
	repo     repository.TagRepository
	logger   logger.LoggerV1
	producer events.Producer
}

func (svc *tagService) AttachTags(ctx context.Context, uid int64, biz string, bizId int64, tagIds []int64) error {
	err := svc.repo.BindTagToBiz(ctx, uid, biz, bizId, tagIds)
	if err != nil {
		return err
	}
	// 异步发送
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		tags, err := svc.repo.GetTagsById(ctx, tagIds)
		cancel()
		if err != nil {
			return
		}
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		err = svc.producer.ProduceSyncEvent(ctx, events.BizTags{
			Biz:   biz,
			BizId: bizId,
			Uid:   uid,
			Tags: slice.Map(tags, func(idx int, src domain.Tag) string {
				return src.Name
			}),
		})
		cancel()
		if err != nil {
			// 记录一下日志
		}
	}()
	return err
}

func (svc *tagService) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	return svc.repo.GetBizTags(ctx, uid, biz, bizId)
}

func (svc *tagService) CreateTag(ctx context.Context, uid int64, name string) (int64, error) {
	return svc.repo.CreateTag(ctx, domain.Tag{
		Uid:  uid,
		Name: name,
	})
}

func (svc *tagService) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	return svc.repo.GetTags(ctx, uid)
}

func NewTagService(repo repository.TagRepository,
	producer events.Producer,
	l logger.LoggerV1) TagService {
	return &tagService{
		producer: producer,
		repo:     repo,
		logger:   l,
	}
}
