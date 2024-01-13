package grpc

import (
	"context"
	"go.opentelemetry.io/otel"
	"log"
	"time"
)

type Server struct {
	UnimplementedUserServiceServer
	Name string
}

var _ UserServiceServer = &Server{}

func (s *Server) GetById(ctx context.Context, request *GetByIdRequest) (*GetByIdResponse, error) {
	ctx, span := otel.Tracer("user_server").Start(ctx, "get_by_id_biz")
	defer span.End()
	ddl, ok := ctx.Deadline()
	if ok {
		// 打印剩余超时时间
		log.Println(ddl.Sub(time.Now()).String())
	}
	time.Sleep(time.Millisecond * 50)

	return &GetByIdResponse{
		User: &User{
			Id:   123,
			Name: "abcd, from " + s.Name,
		},
	}, nil
}
