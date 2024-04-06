package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/account/domain"
	"gitee.com/geekbang/basic-go/webook/account/repository/cache"
	"gitee.com/geekbang/basic-go/webook/account/repository/dao"
	"time"
)

type accountRepository struct {
	dao   dao.AccountDAO
	cache cache.Cache
}

func NewAccountRepository(dao dao.AccountDAO, cache cache.Cache) *accountRepository {
	return &accountRepository{dao: dao, cache: cache}
}

func (a *accountRepository) CheckUnique(ctx context.Context, c domain.Credit) error {
	return a.cache.GetUnique(ctx, c)
}

func (a *accountRepository) SetUnique(ctx context.Context, c domain.Credit) error {
	return a.cache.SetUnique(ctx, c)
}

func (a *accountRepository) AddCredit(ctx context.Context, c domain.Credit) error {
	activities := make([]dao.AccountActivity, 0, len(c.Items))
	now := time.Now().UnixMilli()
	for _, itm := range c.Items {
		activities = append(activities, dao.AccountActivity{
			Uid:         itm.Uid,
			Biz:         c.Biz,
			BizId:       c.BizId,
			Account:     itm.Account,
			AccountType: itm.AccountType.AsUint8(),
			Amount:      itm.Amt,
			Currency:    itm.Currency,
			Ctime:       now,
			Utime:       now,
		})
	}
	return a.dao.AddActivities(ctx, activities...)
}
