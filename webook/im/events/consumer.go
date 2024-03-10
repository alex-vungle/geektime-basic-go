package events

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/im/domain"
	"gitee.com/geekbang/basic-go/webook/im/service"
	"gitee.com/geekbang/basic-go/webook/pkg/canalx"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	svc    service.UserService
}

func (r *MySQLBinlogConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("openim_sync",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"webook_binlog"},
			saramax.NewHandler[canalx.Message[User]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[User]) error {
	if val.Table != "users" || val.Type != "INSERT" {
		// 我不需要处理
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	for _, data := range val.Data {
		// 我这里怎么办？
		err := r.svc.Sync(ctx, domain.User{
			UserID:   strconv.FormatInt(data.Id, 10),
			Nickname: data.Nickname,
		})
		if err != nil {
			// 记录日志
			continue
		}
	}
	return nil
}

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Nickname string `json:"nickname"`
	Birthday int64  `json:"birthday"`
	AboutMe  string `json:"about_me"`
	Phone    string `json:"phone"`

	WechatOpenId  string `json:"wechat_open_id"`
	WechatUnionId string `json:"wechat_union_id"`

	Ctime int64 `json:"ctime"`
	Utime int64 `json:"utime"`
}
