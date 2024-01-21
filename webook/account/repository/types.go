package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/account/domain"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
}
