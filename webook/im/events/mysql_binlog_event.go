package events

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/im/domain"
	"gitee.com/geekbang/basic-go/webook/im/service"
	"gitee.com/geekbang/basic-go/webook/pkg/canalx"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	svc    service.UserService
}

func (r *MySQLBinlogConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("im_users_sync",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"webook_binlog"},
			// 监听 User 的数据
			// 不能直接用过 DAO，因为 DAO 是别的模块的
			saramax.NewHandler[canalx.Message[User]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	cmsg canalx.Message[User]) error {
	// 别的表的 binlog，你不关心
	// 可以考虑，不同的表用不同的 topic，那么你这里就不需要判定了
	if cmsg.Table != "users" {
		return nil
	}

	// 删除用户
	if cmsg.Type == "DELETE" {
		// 你不管都可以
		return nil
	}
	// 要在这里更新缓存了
	// 增删改的消息，实际上在 publish article 里面是没有删的消息的
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, data := range cmsg.Data {
		// 在这里处理了
		// 在这里，把用户数据同步过去
		// 这边能不能把 SyncUser 改成批量接口？可以的，只是说当下是没有必要的。
		// 万一后面有批量插入用户，或者批量更新，你就把这里改成批量接口
		err := r.svc.SyncUser(ctx, domain.User{
			ID:       data.Id,
			Nickname: data.Nickname,
			// 已有代码没有头像字段
			//Avatar: data.
		})
		if err != nil {
			// 记录日志下一条
			continue
		}
	}
	return nil
}

type User struct {
	Id            int64  `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Nickname      string `json:"nickname"`
	Phone         string `json:"phone"`
	WechatUnionID string `json:"wechat_union_id"`
	WechatOpenID  string `json:"wechat_open_id"`
	Ctime         int64  `json:"ctime"`
	Utime         int64  `json:"utime"`
}
