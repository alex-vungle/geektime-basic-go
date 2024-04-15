package events

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/follow/repository/cache"
	"gitee.com/geekbang/basic-go/webook/follow/repository/dao"
	"gitee.com/geekbang/basic-go/webook/pkg/canalx"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	// 直接操作缓存
	cache cache.FollowCache
}

func (r *MySQLBinlogConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("follow",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"webook_binlog"},
			saramax.NewHandler[canalx.Message[FollowRelation]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[FollowRelation]) error {
	if val.Table != "users" || val.Type != "INSERT" {
		// 我不需要处理
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	for _, data := range val.Data {
		// 开始处理
		var err error
		switch data.Status {
		case dao.FollowRelationStatusActive:
			err = r.cache.Follow(ctx, data.Followee, data.Followee)
		case dao.FollowRelationStatusInactive:
			err = r.cache.CancelFollow(ctx, data.Follower, data.Followee)
		default:
			// 记录日志
		}

		if err != nil {
			// 记录日志
		}
	}
	return nil
}

type FollowRelation struct {
	ID int64 `gorm:"column:id;autoIncrement;primaryKey;"`

	Follower int64 `gorm:"uniqueIndex:follower_followee"`
	Followee int64 `gorm:"uniqueIndex:follower_followee"`

	// 软删除策略
	Status uint8

	Ctime int64
	Utime int64
}
