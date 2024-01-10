package grpc

import (
	"context"
	"log"
	"time"
)

type Server struct {
	UnimplementedUserServiceServer
	Name string
}

var _ UserServiceServer = &Server{}

func (s *Server) GetById(ctx context.Context, request *GetByIdRequest) (*GetByIdResponse, error) {
	ddl, ok := ctx.Deadline()
	if ok {
		// 打印剩余超时时间
		log.Println(ddl.Sub(time.Now()).String())
	}

	return &GetByIdResponse{
		User: &User{
			Id:   123,
			Name: "abcd, from " + s.Name,
		},
	}, nil
}
