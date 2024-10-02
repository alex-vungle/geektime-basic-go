// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: article/v1/article.proto

package articlev1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ArticleService_Save_FullMethodName             = "/article.v1.ArticleService/Save"
	ArticleService_Publish_FullMethodName          = "/article.v1.ArticleService/Publish"
	ArticleService_Withdraw_FullMethodName         = "/article.v1.ArticleService/Withdraw"
	ArticleService_List_FullMethodName             = "/article.v1.ArticleService/List"
	ArticleService_GetById_FullMethodName          = "/article.v1.ArticleService/GetById"
	ArticleService_GetPublishedById_FullMethodName = "/article.v1.ArticleService/GetPublishedById"
	ArticleService_ListPub_FullMethodName          = "/article.v1.ArticleService/ListPub"
)

// ArticleServiceClient is the client API for ArticleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ArticleServiceClient interface {
	Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error)
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error)
	GetPublishedById(ctx context.Context, in *GetPublishedByIdRequest, opts ...grpc.CallOption) (*GetPublishedByIdResponse, error)
	ListPub(ctx context.Context, in *ListPubRequest, opts ...grpc.CallOption) (*ListPubResponse, error)
}

type articleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewArticleServiceClient(cc grpc.ClientConnInterface) ArticleServiceClient {
	return &articleServiceClient{cc}
}

func (c *articleServiceClient) Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SaveResponse)
	err := c.cc.Invoke(ctx, ArticleService_Save_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, ArticleService_Publish_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WithdrawResponse)
	err := c.cc.Invoke(ctx, ArticleService_Withdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, ArticleService_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetByIdResponse)
	err := c.cc.Invoke(ctx, ArticleService_GetById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) GetPublishedById(ctx context.Context, in *GetPublishedByIdRequest, opts ...grpc.CallOption) (*GetPublishedByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPublishedByIdResponse)
	err := c.cc.Invoke(ctx, ArticleService_GetPublishedById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articleServiceClient) ListPub(ctx context.Context, in *ListPubRequest, opts ...grpc.CallOption) (*ListPubResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPubResponse)
	err := c.cc.Invoke(ctx, ArticleService_ListPub_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ArticleServiceServer is the server API for ArticleService service.
// All implementations must embed UnimplementedArticleServiceServer
// for forward compatibility.
type ArticleServiceServer interface {
	Save(context.Context, *SaveRequest) (*SaveResponse, error)
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error)
	GetPublishedById(context.Context, *GetPublishedByIdRequest) (*GetPublishedByIdResponse, error)
	ListPub(context.Context, *ListPubRequest) (*ListPubResponse, error)
	mustEmbedUnimplementedArticleServiceServer()
}

// UnimplementedArticleServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedArticleServiceServer struct{}

func (UnimplementedArticleServiceServer) Save(context.Context, *SaveRequest) (*SaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedArticleServiceServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedArticleServiceServer) Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}
func (UnimplementedArticleServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedArticleServiceServer) GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedArticleServiceServer) GetPublishedById(context.Context, *GetPublishedByIdRequest) (*GetPublishedByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPublishedById not implemented")
}
func (UnimplementedArticleServiceServer) ListPub(context.Context, *ListPubRequest) (*ListPubResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPub not implemented")
}
func (UnimplementedArticleServiceServer) mustEmbedUnimplementedArticleServiceServer() {}
func (UnimplementedArticleServiceServer) testEmbeddedByValue()                        {}

// UnsafeArticleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ArticleServiceServer will
// result in compilation errors.
type UnsafeArticleServiceServer interface {
	mustEmbedUnimplementedArticleServiceServer()
}

func RegisterArticleServiceServer(s grpc.ServiceRegistrar, srv ArticleServiceServer) {
	// If the following call pancis, it indicates UnimplementedArticleServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ArticleService_ServiceDesc, srv)
}

func _ArticleService_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ArticleService_Save_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).Save(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ArticleService_Publish_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ArticleService_Withdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).Withdraw(ctx, req.(*WithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ArticleService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ArticleService_GetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_GetPublishedById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPublishedByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).GetPublishedById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ArticleService_GetPublishedById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).GetPublishedById(ctx, req.(*GetPublishedByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticleService_ListPub_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPubRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).ListPub(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ArticleService_ListPub_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).ListPub(ctx, req.(*ListPubRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ArticleService_ServiceDesc is the grpc.ServiceDesc for ArticleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ArticleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "article.v1.ArticleService",
	HandlerType: (*ArticleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Save",
			Handler:    _ArticleService_Save_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _ArticleService_Publish_Handler,
		},
		{
			MethodName: "Withdraw",
			Handler:    _ArticleService_Withdraw_Handler,
		},
		{
			MethodName: "List",
			Handler:    _ArticleService_List_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _ArticleService_GetById_Handler,
		},
		{
			MethodName: "GetPublishedById",
			Handler:    _ArticleService_GetPublishedById_Handler,
		},
		{
			MethodName: "ListPub",
			Handler:    _ArticleService_ListPub_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "article/v1/article.proto",
}
