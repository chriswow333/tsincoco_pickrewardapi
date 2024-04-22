// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: channel.proto

package proto

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

// ChannelV1Client is the client API for ChannelV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChannelV1Client interface {
	GetShowLabels(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*ShowLabelsReply, error)
	GetChannelsByShowLabel(ctx context.Context, in *ShowLabelReq, opts ...grpc.CallOption) (*ChannelsReply, error)
	GetsByChannelIDs(ctx context.Context, in *ChannelIDsReq, opts ...grpc.CallOption) (*ChannelsReply, error)
	SearchChannel(ctx context.Context, in *SearchChannelReq, opts ...grpc.CallOption) (*SearchChannelsReply, error)
}

type channelV1Client struct {
	cc grpc.ClientConnInterface
}

func NewChannelV1Client(cc grpc.ClientConnInterface) ChannelV1Client {
	return &channelV1Client{cc}
}

func (c *channelV1Client) GetShowLabels(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*ShowLabelsReply, error) {
	out := new(ShowLabelsReply)
	err := c.cc.Invoke(ctx, "/channel.v1.ChannelV1/GetShowLabels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelV1Client) GetChannelsByShowLabel(ctx context.Context, in *ShowLabelReq, opts ...grpc.CallOption) (*ChannelsReply, error) {
	out := new(ChannelsReply)
	err := c.cc.Invoke(ctx, "/channel.v1.ChannelV1/GetChannelsByShowLabel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelV1Client) GetsByChannelIDs(ctx context.Context, in *ChannelIDsReq, opts ...grpc.CallOption) (*ChannelsReply, error) {
	out := new(ChannelsReply)
	err := c.cc.Invoke(ctx, "/channel.v1.ChannelV1/GetsByChannelIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelV1Client) SearchChannel(ctx context.Context, in *SearchChannelReq, opts ...grpc.CallOption) (*SearchChannelsReply, error) {
	out := new(SearchChannelsReply)
	err := c.cc.Invoke(ctx, "/channel.v1.ChannelV1/SearchChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChannelV1Server is the server API for ChannelV1 service.
// All implementations must embed UnimplementedChannelV1Server
// for forward compatibility
type ChannelV1Server interface {
	GetShowLabels(context.Context, *EmptyReq) (*ShowLabelsReply, error)
	GetChannelsByShowLabel(context.Context, *ShowLabelReq) (*ChannelsReply, error)
	GetsByChannelIDs(context.Context, *ChannelIDsReq) (*ChannelsReply, error)
	SearchChannel(context.Context, *SearchChannelReq) (*SearchChannelsReply, error)
	mustEmbedUnimplementedChannelV1Server()
}

// UnimplementedChannelV1Server must be embedded to have forward compatible implementations.
type UnimplementedChannelV1Server struct {
}

func (UnimplementedChannelV1Server) GetShowLabels(context.Context, *EmptyReq) (*ShowLabelsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShowLabels not implemented")
}
func (UnimplementedChannelV1Server) GetChannelsByShowLabel(context.Context, *ShowLabelReq) (*ChannelsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannelsByShowLabel not implemented")
}
func (UnimplementedChannelV1Server) GetsByChannelIDs(context.Context, *ChannelIDsReq) (*ChannelsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetsByChannelIDs not implemented")
}
func (UnimplementedChannelV1Server) SearchChannel(context.Context, *SearchChannelReq) (*SearchChannelsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchChannel not implemented")
}
func (UnimplementedChannelV1Server) mustEmbedUnimplementedChannelV1Server() {}

// UnsafeChannelV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChannelV1Server will
// result in compilation errors.
type UnsafeChannelV1Server interface {
	mustEmbedUnimplementedChannelV1Server()
}

func RegisterChannelV1Server(s grpc.ServiceRegistrar, srv ChannelV1Server) {
	s.RegisterService(&ChannelV1_ServiceDesc, srv)
}

func _ChannelV1_GetShowLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelV1Server).GetShowLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/channel.v1.ChannelV1/GetShowLabels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelV1Server).GetShowLabels(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelV1_GetChannelsByShowLabel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShowLabelReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelV1Server).GetChannelsByShowLabel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/channel.v1.ChannelV1/GetChannelsByShowLabel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelV1Server).GetChannelsByShowLabel(ctx, req.(*ShowLabelReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelV1_GetsByChannelIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChannelIDsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelV1Server).GetsByChannelIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/channel.v1.ChannelV1/GetsByChannelIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelV1Server).GetsByChannelIDs(ctx, req.(*ChannelIDsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelV1_SearchChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchChannelReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelV1Server).SearchChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/channel.v1.ChannelV1/SearchChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelV1Server).SearchChannel(ctx, req.(*SearchChannelReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ChannelV1_ServiceDesc is the grpc.ServiceDesc for ChannelV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChannelV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "channel.v1.ChannelV1",
	HandlerType: (*ChannelV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetShowLabels",
			Handler:    _ChannelV1_GetShowLabels_Handler,
		},
		{
			MethodName: "GetChannelsByShowLabel",
			Handler:    _ChannelV1_GetChannelsByShowLabel_Handler,
		},
		{
			MethodName: "GetsByChannelIDs",
			Handler:    _ChannelV1_GetsByChannelIDs_Handler,
		},
		{
			MethodName: "SearchChannel",
			Handler:    _ChannelV1_SearchChannel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "channel.proto",
}
