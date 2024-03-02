package dao

import (
	"context"
	"gorm.io/gorm"
)

// FeedPullEventDAO 拉模型
type FeedPullEventDAO interface {
	CreatePullEvent(ctx context.Context, event FeedPullEvent) error
	FindPullEventList(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
	FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
}

// FeedPullEvent 发件箱
type FeedPullEvent struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 发件人
	UID  int64 `gorm:"index;column:uid"`
	Type string
	// 这边放的就是关键的扩展字段，不同的事件类型，有不同的解析方式
	Content string
	Ctime   int64
	// 正常来说，这个表的数据是不会被更新的
	//Utime int64

	//
}

type feedPullEventDAO struct {
	db *gorm.DB
}

func NewFeedPullEventDAO(db *gorm.DB) FeedPullEventDAO {
	return &feedPullEventDAO{
		db: db,
	}
}

func (f *feedPullEventDAO) FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (f *feedPullEventDAO) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	return f.db.WithContext(ctx).Create(&event).Error
}

func (f *feedPullEventDAO) FindPullEventList(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}
