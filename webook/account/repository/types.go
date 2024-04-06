package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/account/domain"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
	// CheckUnique 如果返回了 error 就说明重复记账了
	CheckUnique(ctx context.Context, c domain.Credit) error
	SetUnique(ctx context.Context, c domain.Credit) error
}
