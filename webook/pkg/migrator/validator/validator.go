package validator

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator"
	events2 "gitee.com/geekbang/basic-go/webook/pkg/migrator/events"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	baseValidator
	batchSize int
	utime     int64
	// 如果没有数据了，就睡眠
	// 如果不是正数，那么就说明直接返回，结束这一次的循环
	// 我很厌恶这种特殊值有特殊含义的做法，但是不得不搞
	sleepInterval time.Duration
}

func NewValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.LoggerV1,
	producer events2.Producer,
) *Validator[T] {
	return &Validator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
		batchSize: 100,
		// 默认是全量校验，并且数据没了就结束
		sleepInterval: 0,
	}
}

func (v *Validator[T]) Utime(utime int64) *Validator[T] {
	v.utime = utime
	return v
}

func (v *Validator[T]) SleepInterval(i time.Duration) *Validator[T] {
	v.sleepInterval = i
	return v
}

// Validate 执行校验。
// 分成两步：
// 1. from => to
func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		return v.baseToTarget(ctx)
	})
	eg.Go(func() error {
		return v.targetToBase(ctx)
	})
	return eg.Wait()
}

// baseToTarget 批量写法，第十三周答案
func (v *Validator[T]) baseToTargetV1(ctx context.Context) error {
	offset := 0
	// 假设说一次 100 条
	const limit = 100
	for {
		var srcs []T
		// 直接取出来一批
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.base.WithContext(dbCtx).
			Order("id").
			Where("utime >= ?", v.utime).
			Offset(offset).Limit(limit).Find(&srcs).Error
		cancel()
		switch err {
		// 在 find 里面其实不会有这个错误
		//case gorm.ErrRecordNotFound:
		case context.Canceled, context.DeadlineExceeded:
			// 超时你可以继续，也可以返回。一般超时都是因为数据库有了问题
			return err
		case nil:
			if len(srcs) == 0 {
				// 结束，没有数据
				return nil
			}
			err = v.dstDiffV1(srcs)
			if err != nil {
				// 直接中断，你也可以考虑继续重试
				return err
			}
		default:
			v.l.Error("src => dst 查询源表失败", logger.Error(err))
		}
		if len(srcs) < limit {
			// 没有数据了
			return nil
		}
		offset += len(srcs)
	}
}

func (v *Validator[T]) dstDiffV1(srcs []T) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ids := slice.Map(srcs, func(idx int, src T) int64 {
		return src.ID()
	})
	var dsts []T
	err := v.target.WithContext(ctx).Where("id IN ?", ids).
		Find(&dsts).Error
	// 让调用者来决定
	if err != nil {
		return err
	}
	dstMap := v.toMap(dsts)
	for _, src := range srcs {
		dst, ok := dstMap[src.ID()]
		if !ok {
			v.notify(src.ID(), events2.InconsistentEventTypeTargetMissing)
			continue
		}
		if !src.CompareTo(dst) {
			v.notify(src.ID(), events2.InconsistentEventTypeNEQ)
		}
	}
	return nil
}

func (v *Validator[T]) toMap(data []T) map[int64]T {
	res := make(map[int64]T, len(data))
	for _, val := range data {
		res[val.ID()] = val
	}
	return res
}

// baseToTarget 从 first 到 second 的验证
func (v *Validator[T]) baseToTarget(ctx context.Context) error {
	offset := 0
	for {
		var src T
		// 这里假定主键的规范都是叫做 id，基本上大部分公司都有这种规范
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.base.WithContext(dbCtx).
			Order("id").
			Where("utime >= ?", v.utime).
			Offset(offset).First(&src).Error
		cancel()
		switch err {
		case gorm.ErrRecordNotFound:
			// 已经没有数据了
			if v.sleepInterval <= 0 {
				return nil
			}
			time.Sleep(v.sleepInterval)
			continue
		case context.Canceled, context.DeadlineExceeded:
			// 退出循环
			return nil
		case nil:
			v.dstDiff(ctx, src)
		default:
			v.l.Error("src => dst 查询源表失败", logger.Error(err))
		}
		offset++
	}
}

func (v *Validator[T]) dstDiff(ctx context.Context, src T) {
	var dst T
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	err := v.target.WithContext(dbCtx).
		Where("id=?", src.ID()).First(&dst).Error
	cancel()
	// 这边要考虑不同的 error
	switch err {
	case gorm.ErrRecordNotFound:
		v.notify(src.ID(), events2.InconsistentEventTypeTargetMissing)
	case nil:
		// 查询到了数据
		equal := src.CompareTo(dst)
		if !equal {
			v.notify(src.ID(), events2.InconsistentEventTypeNEQ)
		}
	default:
		v.l.Error("src => dst 查询目标表失败", logger.Error(err))
	}
}

// targetToBase 反过来，执行 target 到 base 的验证
// 这是为了找出 dst 中多余的数据
func (v *Validator[T]) targetToBase(ctx context.Context) error {
	// 这个我们只需要找出 src 中不存在的 id 就可以了
	offset := 0
	for {
		var ts []T
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.target.WithContext(dbCtx).Model(new(T)).Select("id").Offset(offset).
			Limit(v.batchSize).Find(&ts).Error
		cancel()
		switch err {
		case gorm.ErrRecordNotFound:
			if v.sleepInterval > 0 {
				time.Sleep(v.sleepInterval)
				// 在 sleep 的时候。不需要调整偏移量
				continue
			}
		case context.DeadlineExceeded, context.Canceled:
			return nil
		case nil:
			v.srcMissingRecords(ctx, ts)
		default:
			v.l.Error("dst => src 查询目标表失败", logger.Error(err))
		}
		if len(ts) < v.batchSize {
			// 数据没了
			return nil
		}
		offset += v.batchSize
	}
}

func (v *Validator[T]) srcMissingRecords(ctx context.Context, ts []T) {
	ids := slice.Map(ts, func(idx int, src T) int64 {
		return src.ID()
	})
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	base := v.base.WithContext(dbCtx)
	var srcTs []T
	err := base.Select("id").Where("id IN ?", ids).Find(&srcTs).Error
	switch err {
	case gorm.ErrRecordNotFound:
		// 说明 ids 全部没有
		v.notifySrcMissing(ts)
	case nil:
		// 计算差集
		missing := slice.DiffSetFunc(ts, srcTs, func(src, dst T) bool {
			return src.ID() == dst.ID()
		})
		v.notifySrcMissing(missing)
	default:
		v.l.Error("dst => src 查询源表失败", logger.Error(err))
	}
}

func (v *Validator[T]) notifySrcMissing(ts []T) {
	for _, t := range ts {
		v.notify(t.ID(), events2.InconsistentEventTypeBaseMissing)
	}
}
