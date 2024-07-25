package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
	//go:embed lua/ranking_cnt.lua
	luaRankingCnt string
	//go:embed lua/ranking_set.lua
	luaRankingSet string
)

var RankingUpdateErr = errors.New("指定数据不存在")

const fieldReadCnt = "read_cnt"
const fieldLikeCnt = "like_cnt"
const fieldCollectCnt = "collect_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, id int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, id int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, id int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, res domain.Interactive) error
	IncrLikeRankingIfPresent(ctx context.Context, biz string, bizId int64) error
	SetLikeRankingScore(ctx context.Context, biz string, bizId int64, score int64) error
	TopLikes(ctx context.Context, biz string) ([]domain.Interactive, error)
}

type InteractiveRedisCache struct {
	client redis.Cmdable
}

func NewInteractiveRedisCache(client redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{
		client: client,
	}
}

func (i *InteractiveRedisCache) Set(ctx context.Context,
	biz string, bizId int64,
	res domain.Interactive) error {
	key := i.key(biz, bizId)
	err := i.client.HSet(ctx, key, fieldCollectCnt, res.CollectCnt,
		fieldReadCnt, res.ReadCnt,
		fieldLikeCnt, res.LikeCnt,
	).Err()
	if err != nil {
		return err
	}
	return i.client.Expire(ctx, key, time.Minute*15).Err()
}

func (i *InteractiveRedisCache) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	key := i.key(biz, id)
	res, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	var intr domain.Interactive
	// 这边是可以忽略错误的
	intr.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	intr.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	intr.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	return intr, nil
}

func (i *InteractiveRedisCache) IncrCollectCntIfPresent(ctx context.Context,
	biz string, id int64) error {
	key := i.key(biz, id)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldCollectCnt, 1).Err()
}

func (i *InteractiveRedisCache) IncrLikeCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, 1).Err()
}

func (i *InteractiveRedisCache) DecrLikeCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (i *InteractiveRedisCache) IncrReadCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	// 不是特别需要处理 res
	//res, err := i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Int()
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Err()
}

func (i *InteractiveRedisCache) IncrLikeRankingIfPresent(ctx context.Context, biz string, bizId int64) error {
	res, err := i.client.Eval(ctx, luaRankingCnt, []string{i.rankingKey(biz)}, bizId).Result()
	if err != nil {
		return err
	}
	if res.(int64) == 0 {
		return RankingUpdateErr
	}
	return nil
}

func (i *InteractiveRedisCache) SetLikeRankingScore(ctx context.Context, biz string, bizId int64, score int64) error {
	return i.client.Eval(ctx, luaRankingSet, []string{i.rankingKey(biz)}, bizId, score).Err()
}

func (i *InteractiveRedisCache) TopLikes(ctx context.Context, biz string) ([]domain.Interactive, error) {
	var start int64 = 0
	var end int64 = 99
	res, err := i.client.ZRevRangeWithScores(ctx, i.rankingKey(biz), start, end).Result()
	if err != nil {
		return nil, err
	}
	interactives := make([]domain.Interactive, 0, 100)
	for i := 0; i < len(res); i++ {
		id, _ := strconv.ParseInt(res[i].Member.(string), 10, 64)
		interactives = append(interactives, domain.Interactive{
			Biz:     biz,
			BizId:   id,
			LikeCnt: int64(res[i].Score),
		})
	}
	return interactives, nil
}

func (i *InteractiveRedisCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}

func (i *InteractiveRedisCache) rankingKey(biz string) string {
	return fmt.Sprintf("top_100_%s", biz)
}
