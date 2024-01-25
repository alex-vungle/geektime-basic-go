//go:build wireinject

package startup

import (
	grpc2 "gitee.com/geekbang/basic-go/webook/comment/grpc"
	"gitee.com/geekbang/basic-go/webook/comment/repository"
	"gitee.com/geekbang/basic-go/webook/comment/repository/dao"
	"gitee.com/geekbang/basic-go/webook/comment/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc2.NewGrpcServer,
)

var thirdProvider = wire.NewSet(
	logger.NewNoOpLogger,
	InitTestDB,
)

func InitGRPCServer() *grpc2.CommentServiceServer {
	wire.Build(thirdProvider, serviceProviderSet)
	return new(grpc2.CommentServiceServer)
}
