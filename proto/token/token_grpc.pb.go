// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.0
// source: token.proto

package token

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
	TokenService_RecordTokenUsage_FullMethodName      = "/proto.TokenService/RecordTokenUsage"
	TokenService_UpdateUserTotalTokens_FullMethodName = "/proto.TokenService/UpdateUserTotalTokens"
)

// TokenServiceClient is the client API for TokenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// TokenService defines the methods for managing tokens.
type TokenServiceClient interface {
	// RecordTokenUsage creates a record of token usage for a specific user.
	RecordTokenUsage(ctx context.Context, in *RecordTokenUsageRequest, opts ...grpc.CallOption) (*TokenOperationResponse, error)
	// UpdateUserTotalTokens updates the total number of tokens for a specific user.
	UpdateUserTotalTokens(ctx context.Context, in *UpdateUserTotalTokensRequest, opts ...grpc.CallOption) (*TokenOperationResponse, error)
}

type tokenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenServiceClient(cc grpc.ClientConnInterface) TokenServiceClient {
	return &tokenServiceClient{cc}
}

func (c *tokenServiceClient) RecordTokenUsage(ctx context.Context, in *RecordTokenUsageRequest, opts ...grpc.CallOption) (*TokenOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TokenOperationResponse)
	err := c.cc.Invoke(ctx, TokenService_RecordTokenUsage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) UpdateUserTotalTokens(ctx context.Context, in *UpdateUserTotalTokensRequest, opts ...grpc.CallOption) (*TokenOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TokenOperationResponse)
	err := c.cc.Invoke(ctx, TokenService_UpdateUserTotalTokens_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenServiceServer is the server API for TokenService service.
// All implementations must embed UnimplementedTokenServiceServer
// for forward compatibility.
//
// TokenService defines the methods for managing tokens.
type TokenServiceServer interface {
	// RecordTokenUsage creates a record of token usage for a specific user.
	RecordTokenUsage(context.Context, *RecordTokenUsageRequest) (*TokenOperationResponse, error)
	// UpdateUserTotalTokens updates the total number of tokens for a specific user.
	UpdateUserTotalTokens(context.Context, *UpdateUserTotalTokensRequest) (*TokenOperationResponse, error)
	mustEmbedUnimplementedTokenServiceServer()
}

// UnimplementedTokenServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTokenServiceServer struct{}

func (UnimplementedTokenServiceServer) RecordTokenUsage(context.Context, *RecordTokenUsageRequest) (*TokenOperationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordTokenUsage not implemented")
}
func (UnimplementedTokenServiceServer) UpdateUserTotalTokens(context.Context, *UpdateUserTotalTokensRequest) (*TokenOperationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserTotalTokens not implemented")
}
func (UnimplementedTokenServiceServer) mustEmbedUnimplementedTokenServiceServer() {}
func (UnimplementedTokenServiceServer) testEmbeddedByValue()                      {}

// UnsafeTokenServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenServiceServer will
// result in compilation errors.
type UnsafeTokenServiceServer interface {
	mustEmbedUnimplementedTokenServiceServer()
}

func RegisterTokenServiceServer(s grpc.ServiceRegistrar, srv TokenServiceServer) {
	// If the following call pancis, it indicates UnimplementedTokenServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TokenService_ServiceDesc, srv)
}

func _TokenService_RecordTokenUsage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordTokenUsageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).RecordTokenUsage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_RecordTokenUsage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).RecordTokenUsage(ctx, req.(*RecordTokenUsageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_UpdateUserTotalTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserTotalTokensRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).UpdateUserTotalTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_UpdateUserTotalTokens_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).UpdateUserTotalTokens(ctx, req.(*UpdateUserTotalTokensRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenService_ServiceDesc is the grpc.ServiceDesc for TokenService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TokenService",
	HandlerType: (*TokenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RecordTokenUsage",
			Handler:    _TokenService_RecordTokenUsage_Handler,
		},
		{
			MethodName: "UpdateUserTotalTokens",
			Handler:    _TokenService_UpdateUserTotalTokens_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token.proto",
}
