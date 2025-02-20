// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: card.proto

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

// CardV1Client is the client API for CardV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CardV1Client interface {
	GetCardsByBankID(ctx context.Context, in *CardsByBankIDReq, opts ...grpc.CallOption) (*CardsReply, error)
	GetLatestCards(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*CardsReply, error)
	GetCardByID(ctx context.Context, in *CardIDReq, opts ...grpc.CallOption) (*CardReply, error)
	SearchCard(ctx context.Context, in *SearchCardReq, opts ...grpc.CallOption) (*CardsReply, error)
}

type cardV1Client struct {
	cc grpc.ClientConnInterface
}

func NewCardV1Client(cc grpc.ClientConnInterface) CardV1Client {
	return &cardV1Client{cc}
}

func (c *cardV1Client) GetCardsByBankID(ctx context.Context, in *CardsByBankIDReq, opts ...grpc.CallOption) (*CardsReply, error) {
	out := new(CardsReply)
	err := c.cc.Invoke(ctx, "/card.v1.CardV1/GetCardsByBankID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardV1Client) GetLatestCards(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*CardsReply, error) {
	out := new(CardsReply)
	err := c.cc.Invoke(ctx, "/card.v1.CardV1/GetLatestCards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardV1Client) GetCardByID(ctx context.Context, in *CardIDReq, opts ...grpc.CallOption) (*CardReply, error) {
	out := new(CardReply)
	err := c.cc.Invoke(ctx, "/card.v1.CardV1/GetCardByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardV1Client) SearchCard(ctx context.Context, in *SearchCardReq, opts ...grpc.CallOption) (*CardsReply, error) {
	out := new(CardsReply)
	err := c.cc.Invoke(ctx, "/card.v1.CardV1/SearchCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CardV1Server is the server API for CardV1 service.
// All implementations must embed UnimplementedCardV1Server
// for forward compatibility
type CardV1Server interface {
	GetCardsByBankID(context.Context, *CardsByBankIDReq) (*CardsReply, error)
	GetLatestCards(context.Context, *EmptyReq) (*CardsReply, error)
	GetCardByID(context.Context, *CardIDReq) (*CardReply, error)
	SearchCard(context.Context, *SearchCardReq) (*CardsReply, error)
	mustEmbedUnimplementedCardV1Server()
}

// UnimplementedCardV1Server must be embedded to have forward compatible implementations.
type UnimplementedCardV1Server struct {
}

func (UnimplementedCardV1Server) GetCardsByBankID(context.Context, *CardsByBankIDReq) (*CardsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCardsByBankID not implemented")
}
func (UnimplementedCardV1Server) GetLatestCards(context.Context, *EmptyReq) (*CardsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLatestCards not implemented")
}
func (UnimplementedCardV1Server) GetCardByID(context.Context, *CardIDReq) (*CardReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCardByID not implemented")
}
func (UnimplementedCardV1Server) SearchCard(context.Context, *SearchCardReq) (*CardsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchCard not implemented")
}
func (UnimplementedCardV1Server) mustEmbedUnimplementedCardV1Server() {}

// UnsafeCardV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CardV1Server will
// result in compilation errors.
type UnsafeCardV1Server interface {
	mustEmbedUnimplementedCardV1Server()
}

func RegisterCardV1Server(s grpc.ServiceRegistrar, srv CardV1Server) {
	s.RegisterService(&CardV1_ServiceDesc, srv)
}

func _CardV1_GetCardsByBankID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CardsByBankIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardV1Server).GetCardsByBankID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/card.v1.CardV1/GetCardsByBankID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardV1Server).GetCardsByBankID(ctx, req.(*CardsByBankIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardV1_GetLatestCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardV1Server).GetLatestCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/card.v1.CardV1/GetLatestCards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardV1Server).GetLatestCards(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardV1_GetCardByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CardIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardV1Server).GetCardByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/card.v1.CardV1/GetCardByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardV1Server).GetCardByID(ctx, req.(*CardIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardV1_SearchCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchCardReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardV1Server).SearchCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/card.v1.CardV1/SearchCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardV1Server).SearchCard(ctx, req.(*SearchCardReq))
	}
	return interceptor(ctx, in, info, handler)
}

// CardV1_ServiceDesc is the grpc.ServiceDesc for CardV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CardV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "card.v1.CardV1",
	HandlerType: (*CardV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCardsByBankID",
			Handler:    _CardV1_GetCardsByBankID_Handler,
		},
		{
			MethodName: "GetLatestCards",
			Handler:    _CardV1_GetLatestCards_Handler,
		},
		{
			MethodName: "GetCardByID",
			Handler:    _CardV1_GetCardByID_Handler,
		},
		{
			MethodName: "SearchCard",
			Handler:    _CardV1_SearchCard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "card.proto",
}
