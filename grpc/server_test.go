package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	// 这个是 grpc 的 Server
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(first, second))
	defer func() {
		// 优雅退出
		server.GracefulStop()
	}()
	// 我们的业务的 server
	userServer := &Server{}
	RegisterUserServiceServer(server, userServer)
	// 创建一个监听器，监听 tcp 协议，8090 端口
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = server.Serve(l)
	t.Log(err)
}

var first grpc.UnaryServerInterceptor = func(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("这是第一个前")
	resp, err = handler(ctx, req)
	log.Println("这是第一个后")
	return
}

var second grpc.UnaryServerInterceptor = func(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("这是第二个前")
	resp, err = handler(ctx, req)
	log.Println("这是第二个后")
	return
}
