package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/payment/domain"
	"gitee.com/geekbang/basic-go/webook/payment/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

var _ LocalMsgRepository = (*LocalMsgGORMRepository)(nil)

type LocalMsgRepository interface {
	AddMsg(ctx context.Context, content string) (int64, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
	FindInitMsg(ctx context.Context, offset, limit int) ([]domain.Msg, error)
}

type LocalMsgGORMRepository struct {
	dao *dao.LocalMsgGORMDAO
}

func (l *LocalMsgGORMRepository) FindInitMsg(ctx context.Context, offset, limit int) ([]domain.Msg, error) {
	msgs, err := l.dao.FindInitMsg(ctx, offset, limit)
	return slice.Map(msgs, func(idx int, src dao.Msg) domain.Msg {
		return domain.Msg{
			Id:      src.Id,
			Content: src.Content,
			Ctime:   time.UnixMilli(src.Ctime),
		}
	}), err
}

func (l *LocalMsgGORMRepository) MarkFailed(ctx context.Context, id int64) error {
	return l.dao.UpdateStatus(ctx, id, dao.MsgStatusFailed)
}

func (l *LocalMsgGORMRepository) MarkSuccess(ctx context.Context, id int64) error {
	return l.dao.UpdateStatus(ctx, id, dao.MsgStatusSuccess)
}

func NewLocalMsgGORMRepository(db *gorm.DB) *LocalMsgGORMRepository {
	return &LocalMsgGORMRepository{dao: dao.NewLocalMsgGORMDAO(db)}
}

func (l *LocalMsgGORMRepository) AddMsg(ctx context.Context, content string) (int64, error) {
	return l.dao.AddMsg(ctx, dao.Msg{
		Content: content,
	})
}
