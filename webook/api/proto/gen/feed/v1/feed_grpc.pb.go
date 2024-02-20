// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: feed/v1/feed.proto

package feedv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	FeedSvc_CreateFeedEvent_FullMethodName = "/feed.v1.FeedSvc/CreateFeedEvent"
	FeedSvc_FindFeedEvents_FullMethodName  = "/feed.v1.FeedSvc/FindFeedEvents"
)

// FeedSvcClient is the client API for FeedSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FeedSvcClient interface {
	CreateFeedEvent(ctx context.Context, in *CreateFeedEventRequest, opts ...grpc.CallOption) (*CreateFeedEventResponse, error)
	FindFeedEvents(ctx context.Context, in *FindFeedEventsRequest, opts ...grpc.CallOption) (*FindFeedEventsResponse, error)
}

type feedSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewFeedSvcClient(cc grpc.ClientConnInterface) FeedSvcClient {
	return &feedSvcClient{cc}
}

func (c *feedSvcClient) CreateFeedEvent(ctx context.Context, in *CreateFeedEventRequest, opts ...grpc.CallOption) (*CreateFeedEventResponse, error) {
	out := new(CreateFeedEventResponse)
	err := c.cc.Invoke(ctx, FeedSvc_CreateFeedEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedSvcClient) FindFeedEvents(ctx context.Context, in *FindFeedEventsRequest, opts ...grpc.CallOption) (*FindFeedEventsResponse, error) {
	out := new(FindFeedEventsResponse)
	err := c.cc.Invoke(ctx, FeedSvc_FindFeedEvents_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeedSvcServer is the server API for FeedSvc service.
// All implementations must embed UnimplementedFeedSvcServer
// for forward compatibility
type FeedSvcServer interface {
	CreateFeedEvent(context.Context, *CreateFeedEventRequest) (*CreateFeedEventResponse, error)
	FindFeedEvents(context.Context, *FindFeedEventsRequest) (*FindFeedEventsResponse, error)
	mustEmbedUnimplementedFeedSvcServer()
}

// UnimplementedFeedSvcServer must be embedded to have forward compatible implementations.
type UnimplementedFeedSvcServer struct {
}

func (UnimplementedFeedSvcServer) CreateFeedEvent(context.Context, *CreateFeedEventRequest) (*CreateFeedEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFeedEvent not implemented")
}
func (UnimplementedFeedSvcServer) FindFeedEvents(context.Context, *FindFeedEventsRequest) (*FindFeedEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindFeedEvents not implemented")
}
func (UnimplementedFeedSvcServer) mustEmbedUnimplementedFeedSvcServer() {}

// UnsafeFeedSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FeedSvcServer will
// result in compilation errors.
type UnsafeFeedSvcServer interface {
	mustEmbedUnimplementedFeedSvcServer()
}

func RegisterFeedSvcServer(s grpc.ServiceRegistrar, srv FeedSvcServer) {
	s.RegisterService(&FeedSvc_ServiceDesc, srv)
}

func _FeedSvc_CreateFeedEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFeedEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedSvcServer).CreateFeedEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedSvc_CreateFeedEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedSvcServer).CreateFeedEvent(ctx, req.(*CreateFeedEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedSvc_FindFeedEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindFeedEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedSvcServer).FindFeedEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedSvc_FindFeedEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedSvcServer).FindFeedEvents(ctx, req.(*FindFeedEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FeedSvc_ServiceDesc is the grpc.ServiceDesc for FeedSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FeedSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "feed.v1.FeedSvc",
	HandlerType: (*FeedSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFeedEvent",
			Handler:    _FeedSvc_CreateFeedEvent_Handler,
		},
		{
			MethodName: "FindFeedEvents",
			Handler:    _FeedSvc_FindFeedEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "feed/v1/feed.proto",
}
