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

func (s *Server) GetByID(ctx context.Context, request *GetByIDRequest) (*GetByIDResponse, error) {
	ddl, ok := ctx.Deadline()
	if ok {
		rest := ddl.Sub(time.Now())
		log.Println(rest.String())
	}
	return &GetByIDResponse{
		User: &User{
			Id:   123,
			Name: "from" + s.Name,
		},
	}, nil
}
