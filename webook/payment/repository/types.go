package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/payment/domain"
	"time"
)

//go:generate mockgen -source=types.go -destination=mocks/payment.mock.go --package=repomocks PaymentRepository
type PaymentRepository interface {
	AddPayment(ctx context.Context, pmt domain.Payment) error
	// UpdatePayment 这个设计有点差，因为
	UpdatePayment(ctx context.Context, pmt domain.Payment) error
	FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error)
}
