// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

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

// BitBoxClient is the client API for BitBox service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BitBoxClient interface {
	// Start initiates a process.
	Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartReply, error)
	// Stop halts a process.
	Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopReply, error)
	// Status returns the status of a process.
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error)
	// Query streams the output/result of a process.
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (BitBox_QueryClient, error)
}

type bitBoxClient struct {
	cc grpc.ClientConnInterface
}

func NewBitBoxClient(cc grpc.ClientConnInterface) BitBoxClient {
	return &bitBoxClient{cc}
}

func (c *bitBoxClient) Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartReply, error) {
	out := new(StartReply)
	err := c.cc.Invoke(ctx, "/grpc.BitBox/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bitBoxClient) Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopReply, error) {
	out := new(StopReply)
	err := c.cc.Invoke(ctx, "/grpc.BitBox/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bitBoxClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error) {
	out := new(StatusReply)
	err := c.cc.Invoke(ctx, "/grpc.BitBox/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bitBoxClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (BitBox_QueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &BitBox_ServiceDesc.Streams[0], "/grpc.BitBox/Query", opts...)
	if err != nil {
		return nil, err
	}
	x := &bitBoxQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BitBox_QueryClient interface {
	Recv() (*QueryReply, error)
	grpc.ClientStream
}

type bitBoxQueryClient struct {
	grpc.ClientStream
}

func (x *bitBoxQueryClient) Recv() (*QueryReply, error) {
	m := new(QueryReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BitBoxServer is the server API for BitBox service.
// All implementations must embed UnimplementedBitBoxServer
// for forward compatibility
type BitBoxServer interface {
	// Start initiates a process.
	Start(context.Context, *StartRequest) (*StartReply, error)
	// Stop halts a process.
	Stop(context.Context, *StopRequest) (*StopReply, error)
	// Status returns the status of a process.
	Status(context.Context, *StatusRequest) (*StatusReply, error)
	// Query streams the output/result of a process.
	Query(*QueryRequest, BitBox_QueryServer) error
	mustEmbedUnimplementedBitBoxServer()
}

// UnimplementedBitBoxServer must be embedded to have forward compatible implementations.
type UnimplementedBitBoxServer struct {
}

func (UnimplementedBitBoxServer) Start(context.Context, *StartRequest) (*StartReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedBitBoxServer) Stop(context.Context, *StopRequest) (*StopReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedBitBoxServer) Status(context.Context, *StatusRequest) (*StatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedBitBoxServer) Query(*QueryRequest, BitBox_QueryServer) error {
	return status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedBitBoxServer) mustEmbedUnimplementedBitBoxServer() {}

// UnsafeBitBoxServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BitBoxServer will
// result in compilation errors.
type UnsafeBitBoxServer interface {
	mustEmbedUnimplementedBitBoxServer()
}

func RegisterBitBoxServer(s grpc.ServiceRegistrar, srv BitBoxServer) {
	s.RegisterService(&BitBox_ServiceDesc, srv)
}

func _BitBox_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitBoxServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.BitBox/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitBoxServer).Start(ctx, req.(*StartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BitBox_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitBoxServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.BitBox/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitBoxServer).Stop(ctx, req.(*StopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BitBox_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitBoxServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.BitBox/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitBoxServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BitBox_Query_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BitBoxServer).Query(m, &bitBoxQueryServer{stream})
}

type BitBox_QueryServer interface {
	Send(*QueryReply) error
	grpc.ServerStream
}

type bitBoxQueryServer struct {
	grpc.ServerStream
}

func (x *bitBoxQueryServer) Send(m *QueryReply) error {
	return x.ServerStream.SendMsg(m)
}

// BitBox_ServiceDesc is the grpc.ServiceDesc for BitBox service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BitBox_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.BitBox",
	HandlerType: (*BitBoxServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Start",
			Handler:    _BitBox_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _BitBox_Stop_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _BitBox_Status_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _BitBox_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "bitbox.proto",
}
