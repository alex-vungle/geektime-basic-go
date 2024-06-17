package cache

import (
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

type LocalCodeCache struct {
	cache *cache.Cache
}

func NewLocalCodeCache() CodeCache {
	return &LocalCodeCache{
		cache: cache.New(5*time.Minute, 1*time.Minute),
	}
}

func (c *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {

	key := c.key(biz, phone)
	_, found := c.cache.Get(key)

	if found {
		return ErrCodeSendTooMany
	}

	c.cache.Set(key, code, 5*time.Minute)
	return nil

}

func (c *LocalCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	key := c.key(biz, phone)
	val, found := c.cache.Get(key)

	if !found {
		return false, nil
	}
	if val == code {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
