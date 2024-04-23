package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

//go:generate mockgen -source=./user.go -package=cachemocks -destination=./mocks/user.mock.go UserCache
type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
	Del(ctx context.Context, id int64) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *RedisUserCache) Del(ctx context.Context, id int64) error {
	return c.cmd.Del(ctx, c.key(id)).Err()
}

func (c *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.key(uid)
	// 我假定这个地方用 JSON 来
	data, err := c.cmd.Get(ctx, key).Result()
	//data, err := c.cmd.Get(ctx, firstKey).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	//if err != nil {
	//	return domain.User{}, err
	//}
	//return u, nil
	return u, err
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.Id)
	// 我假定这个地方用 JSON
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisUserCache) key(uid int64) string {
	// user-info-
	// user.info.
	// user/info/
	// user_info_
	return fmt.Sprintf("user:info:%d", uid)
}

type UserCacheV1 struct {
	client *redis.Client
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisUserCache) IsLogoutBF(ctx context.Context, id int64) (bool, error) {
	// 是否已经退出登录
	// 考虑假阳性的问题，你就再次调用 Redis 检测 session（token） 是不是过期了
	return c.cmd.BFExists(ctx, "logout_users", id).Result()
}

// 一定不要自己去初始化你需要的东西，让外面传进来
//func NewUserCacheV1(addr string) *RedisUserCache {
//	cmd := redis.NewClient(&redis.Options{Addr: addr})
//	return &RedisUserCache{
//		cmd:        cmd,
//		expiration: time.Minute * 15,
//	}
//}
