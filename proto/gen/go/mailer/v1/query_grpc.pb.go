// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: mailer/v1/query.proto

package mailerv1

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
	MailerQueryService_Get_FullMethodName       = "/rashop.mailer.v1.MailerQueryService/Get"
	MailerQueryService_FindByIds_FullMethodName = "/rashop.mailer.v1.MailerQueryService/FindByIds"
)

// MailerQueryServiceClient is the client API for MailerQueryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MailerQueryServiceClient interface {
	Get(ctx context.Context, in *GetMailsRequest, opts ...grpc.CallOption) (*GetMailsResponse, error)
	FindByIds(ctx context.Context, in *FindMailByIdsRequest, opts ...grpc.CallOption) (*FindMailByIdsResponse, error)
}

type mailerQueryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMailerQueryServiceClient(cc grpc.ClientConnInterface) MailerQueryServiceClient {
	return &mailerQueryServiceClient{cc}
}

func (c *mailerQueryServiceClient) Get(ctx context.Context, in *GetMailsRequest, opts ...grpc.CallOption) (*GetMailsResponse, error) {
	out := new(GetMailsResponse)
	err := c.cc.Invoke(ctx, MailerQueryService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailerQueryServiceClient) FindByIds(ctx context.Context, in *FindMailByIdsRequest, opts ...grpc.CallOption) (*FindMailByIdsResponse, error) {
	out := new(FindMailByIdsResponse)
	err := c.cc.Invoke(ctx, MailerQueryService_FindByIds_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MailerQueryServiceServer is the server API for MailerQueryService service.
// All implementations must embed UnimplementedMailerQueryServiceServer
// for forward compatibility
type MailerQueryServiceServer interface {
	Get(context.Context, *GetMailsRequest) (*GetMailsResponse, error)
	FindByIds(context.Context, *FindMailByIdsRequest) (*FindMailByIdsResponse, error)
	mustEmbedUnimplementedMailerQueryServiceServer()
}

// UnimplementedMailerQueryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMailerQueryServiceServer struct {
}

func (UnimplementedMailerQueryServiceServer) Get(context.Context, *GetMailsRequest) (*GetMailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedMailerQueryServiceServer) FindByIds(context.Context, *FindMailByIdsRequest) (*FindMailByIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByIds not implemented")
}
func (UnimplementedMailerQueryServiceServer) mustEmbedUnimplementedMailerQueryServiceServer() {}

// UnsafeMailerQueryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MailerQueryServiceServer will
// result in compilation errors.
type UnsafeMailerQueryServiceServer interface {
	mustEmbedUnimplementedMailerQueryServiceServer()
}

func RegisterMailerQueryServiceServer(s grpc.ServiceRegistrar, srv MailerQueryServiceServer) {
	s.RegisterService(&MailerQueryService_ServiceDesc, srv)
}

func _MailerQueryService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerQueryServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerQueryService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerQueryServiceServer).Get(ctx, req.(*GetMailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailerQueryService_FindByIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindMailByIdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerQueryServiceServer).FindByIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerQueryService_FindByIds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerQueryServiceServer).FindByIds(ctx, req.(*FindMailByIdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MailerQueryService_ServiceDesc is the grpc.ServiceDesc for MailerQueryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MailerQueryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rashop.mailer.v1.MailerQueryService",
	HandlerType: (*MailerQueryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _MailerQueryService_Get_Handler,
		},
		{
			MethodName: "FindByIds",
			Handler:    _MailerQueryService_FindByIds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mailer/v1/query.proto",
}
