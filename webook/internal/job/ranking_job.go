package job

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/google/uuid"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
	"sync"
	"sync/atomic"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	l       logger.LoggerV1
	timeout time.Duration
	client  *rlock.Client
	key     string

	localLock *sync.Mutex
	lock      *rlock.Lock

	// 作业提示
	// 随机生成一个，就代表当前负载。你可以每隔一分钟生成一个
	load *atomic.Int32

	redisClient    redis.Cmdable
	loadTicker     *time.Ticker
	nodeID         string
	rankingLoadKey string
	closeSignal    chan struct{}
}

func NewRankingJob(
	svc service.RankingService,
	l logger.LoggerV1,
	client *rlock.Client,
	timeout time.Duration,
	redisClient redis.Cmdable) *RankingJob {
	rankingJob := &RankingJob{svc: svc,
		key:            "job:ranking",
		l:              l,
		client:         client,
		localLock:      &sync.Mutex{},
		timeout:        timeout,
		redisClient:    redisClient,
		loadTicker:     time.NewTicker(5 * time.Second),
		nodeID:         uuid.New().String(),
		load:           &atomic.Int32{},
		rankingLoadKey: "ranking_job_load",
		closeSignal:    make(chan struct{}),
	}
	rankingJob.loadCycle()
	return rankingJob
}

func (r *RankingJob) loadCycle() {
	go func() {
		for range r.loadTicker.C {
			r.updateLoad()
			r.releaseLockIfNeeded()
		}
	}()
}

func (r *RankingJob) updateLoad() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	load := r.load.Load()
	r.redisClient.ZAdd(ctx, r.rankingLoadKey,
		redis.Z{Member: r.nodeID, Score: float64(load)})
	cancel()
	return
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	r.localLock.Lock()
	lock := r.lock
	if lock == nil {
		// 抢分布式锁
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout,
			&rlock.FixIntervalRetry{
				Interval: time.Millisecond * 100,
				Max:      3,
				// 重试的超时
			}, time.Second)
		if err != nil {
			r.localLock.Unlock()
			r.l.Warn("获取分布式锁失败", logger.Error(err))
			return nil
		}
		r.lock = lock
		r.localLock.Unlock()
		go func() {
			// 并不是非得一半就续约
			er := lock.AutoRefresh(r.timeout/2, r.timeout)
			if er != nil {
				// 续约失败了
				// 你也没办法中断当下正在调度的热榜计算（如果有）
				r.localLock.Lock()
				r.lock = nil
				r.localLock.Unlock()
			}
		}()
	}
	// 这边就是你拿到了锁
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}

func (r *RankingJob) releaseLockIfNeeded() {
	r.localLock.Lock()
	lock := r.lock
	r.localLock.Unlock()
	if lock != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		res, err := r.redisClient.ZPopMin(ctx, r.rankingLoadKey).Result()
		if err != nil {
			r.l.Error("获取节点负载数据失败")
			return
		}
		head := res[0]
		if head.Member.(string) != r.nodeID {
			r.l.Debug(r.nodeID+" 不是负载最低的节点，释放分布式锁",
				logger.Field{Key: "head", Val: head})
			r.localLock.Lock()
			r.lock = nil
			r.localLock.Unlock()
			lock.Unlock(ctx)
		}
	}
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if r.loadTicker != nil {
		r.loadTicker.Stop()
	}

	return lock.Unlock(ctx)
}
