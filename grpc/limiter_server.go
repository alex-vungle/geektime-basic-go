package grpc

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/ratelimit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LimiterUserServer struct {
	limiter ratelimit.Limiter
	l       logger.LoggerV1
	UserServiceServer
}

func (s *LimiterUserServer) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	limited, err := s.limiter.Limit(ctx,

		//fmt.Sprintf("limiter:user:get_by_id:%d", req.Id))
		// limiter:user:456
		fmt.Sprintf("limiter:user:%d", req.Id))
	if err != nil {
		// err 不为nil，你要考虑你用保守的，还是用激进的策略
		// 这是保守的策略
		s.l.Error("判定限流出现问题", logger.Error(err))
		return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
		// 这是激进的策略
		// return handler(ctx, req)
	}
	if limited {
		return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
	}

	resp, err := s.UserServiceServer.GetById(ctx, req)
	return resp, err
}

//func (s *LimiterUserServer) UpdateById(ctx context.Context, req *UpdateByIdRequest) (*UpdateByIdRequest, error) {
//
//}
