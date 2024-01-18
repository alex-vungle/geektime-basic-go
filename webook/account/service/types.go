package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/account/domain"
)

type AccountService interface {
	Credit(ctx context.Context, cr domain.Credit) error
}
