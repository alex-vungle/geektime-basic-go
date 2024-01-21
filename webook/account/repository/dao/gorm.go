package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type AccountGORMDAO struct {
	db *gorm.DB
}

func NewCreditGORMDAO(db *gorm.DB) AccountDAO {
	return &AccountGORMDAO{db: db}
}

func (c *AccountGORMDAO) AddActivities(ctx context.Context, activities ...AccountActivity) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		// 针对每一个 activity，入账
		for _, act := range activities {
			err := c.db.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]any{
					"balance": gorm.Expr("`balance`+?", act.Amount),
					"utime":   now,
				}),
			}).Create(&Account{
				Uid:      act.Uid,
				Account:  act.Account,
				Type:     act.AccountType,
				Balance:  act.Amount,
				Currency: act.Currency,
				Utime:    now,
				Ctime:    now,
			}).Error
			if err != nil {
				return err
			}
		}
		return tx.Create(&activities).Error
	})
}
