package cache

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/account/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	SetUnique(ctx context.Context, cr domain.Credit) error
	GetUnique(ctx context.Context, cr domain.Credit) error
}

type RedisCache struct {
	client redis.Cmdable
}

func (r *RedisCache) SetUnique(ctx context.Context, cr domain.Credit) error {
	return r.client.Set(ctx, r.key(cr.Biz, cr.BizId), "", time.Minute*30).Err()
}

func (r *RedisCache) GetUnique(ctx context.Context, cr domain.Credit) error {
	res, err := r.client.Exists(ctx, r.key(cr.Biz, cr.BizId)).Result()
	if err != nil {
		return err
	}
	if res > 0 {
		return errors.New("该业务已经处理过了")
	}
	return nil
}

func (r *RedisCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("credit:biz:%s_%d", biz, bizId)
}
