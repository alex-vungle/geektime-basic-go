package validator

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator"
	events2 "gitee.com/geekbang/basic-go/webook/pkg/migrator/events"
	"gorm.io/gorm"
)

type CanalIncrValidator[T migrator.Entity] struct {
	baseValidator
}

func NewCanalIncrValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.LoggerV1,
	producer events2.Producer,
) *CanalIncrValidator[T] {
	return &CanalIncrValidator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
	}
}

// Validate 一次校验一条
// id 是被修改的数据的主键
func (v *CanalIncrValidator[T]) Validate(ctx context.Context, id int64) error {
	var base T
	err := v.base.WithContext(ctx).Where("id = ? ", id).First(&base).Error
	switch err {
	case nil:
		// 找到了
		var target T
		err = v.target.WithContext(ctx).Where("id = ?", id).First(&target).Error
		switch err {
		case nil:
			// 两边都找到了
			if !base.CompareTo(target) {
				v.notify(id, events2.InconsistentEventTypeNEQ)
			}
			return nil
		case gorm.ErrRecordNotFound:
			// base 有，target 没有
			v.notify(id, events2.InconsistentEventTypeTargetMissing)
			return nil
		default:
			return err
		}
	case gorm.ErrRecordNotFound:
		// 没找到
		var target T
		err = v.target.WithContext(ctx).Where("id = ?", id).First(&target).Error
		switch err {
		case nil:
			// target 找到了, base 没有
			v.notify(id, events2.InconsistentEventTypeBaseMissing)
			return nil
		case gorm.ErrRecordNotFound:
			return nil
		default:
			return err
		}
	default:
		//	不知道啥错误
		return err
	}
}
