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
		cache: cache.New(2*time.Minute, 20*time.Second),
	}
}

func (c *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {

	key := c.key(biz, phone)
	keyCnt := fmt.Sprintf("%s:cnt", key)
	_, found := c.cache.Get(key)

	if found {
		return ErrCodeSendTooMany
	}

	c.cache.Set(key, code, 2*time.Minute)
	c.cache.Set(keyCnt, 2, -1)

	return nil
}

func (c *LocalCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {

	key := c.key(biz, phone)
	keyCnt := fmt.Sprintf("%s:cnt", key)

	actualCode, foundKey := c.cache.Get(key)
	cnt, foundKeyCnt := c.cache.Get(keyCnt)
	iCnt := cnt.(int)

	if !foundKey {
		return false, nil
	}

	if !foundKeyCnt || cnt == nil || iCnt < 0 {
		return false, ErrCodeVerifyTooMany
	}

	// Verify Code 比较
	if actualCode == code {
		c.cache.Delete(key)
		c.cache.Delete(keyCnt)
		return true, nil
	} else {
		iCnt = iCnt - 1
		c.cache.Set(keyCnt, iCnt, -1)
		return false, nil
	}
}

func (c *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
