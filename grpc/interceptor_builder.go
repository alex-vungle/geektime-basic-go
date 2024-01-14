package grpc

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	limiter limiter.Limiter
	key     string
}

// NewInterceptorBuilder key 1. limiter:interactive-service => 整个点赞的应用限流
func NewInterceptorBuilder(limiter limiter.Limiter, key string) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key}
}

func (b *InterceptorBuilder) BuildServerUnaryInterceptorBiz() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if getReq, ok := req.(*GetByIDRequest); ok {
			key := fmt.Sprintf("limiter:user:get_by_id:%d", getReq.Id)
			limited, err := b.limiter.Limit(ctx, key)
			if err != nil {
				// 你有保守的做法，也有激进的做法
				// 这个是保守的做法
				return nil, status.Errorf(codes.ResourceExhausted, "限流")
			}

			if limited {
				return nil, status.Errorf(codes.ResourceExhausted, "限流")
			}
		}
		return handler(ctx, req)
	}
}
