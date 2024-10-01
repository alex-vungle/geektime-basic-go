package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Msg struct {
	Id      int64 `gorm:"primaryKey,autoIncrement"`
	Content string
	Status  uint8
	Ctime   int64
	Utime   int64 `gorm:"index"`
}

const (
	MsgStatusUnknown = 0
	MsgStatusInit    = 1
	MsgStatusSuccess = 2
	MsgStatusFailed  = 3
)

var _ LocalMsgDAO = (*LocalMsgGORMDAO)(nil)

type LocalMsgDAO interface {
	AddMsg(ctx context.Context, msg Msg) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status uint8) error
}

type LocalMsgGORMDAO struct {
	db *gorm.DB
}

func (dao *LocalMsgGORMDAO) UpdateStatus(ctx context.Context, id int64, status uint8) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Msg{}).Where("id = ?", id).Updates(map[string]any{
		"status": status,
		"utime":  now,
	}).Error
}

func NewLocalMsgGORMDAO(db *gorm.DB) *LocalMsgGORMDAO {
	return &LocalMsgGORMDAO{db: db}
}

func (dao *LocalMsgGORMDAO) AddMsg(ctx context.Context, msg Msg) (int64, error) {
	now := time.Now().UnixMilli()
	msg.Ctime = now
	msg.Utime = now
	err := dao.db.WithContext(ctx).Create(&msg).Error
	return msg.Id, err
}

func (dao *LocalMsgGORMDAO) FindInitMsg(ctx context.Context, offset int, limit int) ([]Msg, error) {
	var msgs []Msg
	err := dao.db.WithContext(ctx).
		Where("status = ?", MsgStatusInit).
		Offset(offset).Limit(limit).
		Find(&msgs).Error
	return msgs, err
}
