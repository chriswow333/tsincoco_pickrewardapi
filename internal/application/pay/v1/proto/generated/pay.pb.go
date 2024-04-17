// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: pay.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Reply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status int32  `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Error  *Error `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Reply) Reset() {
	*x = Reply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reply) ProtoMessage() {}

func (x *Reply) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reply.ProtoReflect.Descriptor instead.
func (*Reply) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{0}
}

func (x *Reply) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Reply) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorCode    int32  `protobuf:"varint,1,opt,name=errorCode,proto3" json:"errorCode,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=errorMessage,proto3" json:"errorMessage,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{1}
}

func (x *Error) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *Error) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type PayIDReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PayIDReq) Reset() {
	*x = PayIDReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayIDReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayIDReq) ProtoMessage() {}

func (x *PayIDReq) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayIDReq.ProtoReflect.Descriptor instead.
func (*PayIDReq) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{2}
}

func (x *PayIDReq) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type PaysReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reply *Reply           `protobuf:"bytes,1,opt,name=reply,proto3" json:"reply,omitempty"`
	Pays  []*PaysReply_Pay `protobuf:"bytes,2,rep,name=pays,proto3" json:"pays,omitempty"`
}

func (x *PaysReply) Reset() {
	*x = PaysReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaysReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaysReply) ProtoMessage() {}

func (x *PaysReply) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaysReply.ProtoReflect.Descriptor instead.
func (*PaysReply) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{3}
}

func (x *PaysReply) GetReply() *Reply {
	if x != nil {
		return x.Reply
	}
	return nil
}

func (x *PaysReply) GetPays() []*PaysReply_Pay {
	if x != nil {
		return x.Pays
	}
	return nil
}

type PayReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reply *Reply        `protobuf:"bytes,1,opt,name=reply,proto3" json:"reply,omitempty"`
	Pay   *PayReply_Pay `protobuf:"bytes,2,opt,name=pay,proto3" json:"pay,omitempty"`
}

func (x *PayReply) Reset() {
	*x = PayReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayReply) ProtoMessage() {}

func (x *PayReply) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayReply.ProtoReflect.Descriptor instead.
func (*PayReply) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{4}
}

func (x *PayReply) GetReply() *Reply {
	if x != nil {
		return x.Reply
	}
	return nil
}

func (x *PayReply) GetPay() *PayReply_Pay {
	if x != nil {
		return x.Pay
	}
	return nil
}

type EmptyReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyReq) Reset() {
	*x = EmptyReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyReq) ProtoMessage() {}

func (x *EmptyReq) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyReq.ProtoReflect.Descriptor instead.
func (*EmptyReq) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{5}
}

type PaysReply_Pay struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Order      int32  `protobuf:"varint,3,opt,name=order,proto3" json:"order,omitempty"`
	PayStatus  int32  `protobuf:"varint,4,opt,name=payStatus,proto3" json:"payStatus,omitempty"`
	CreateDate int64  `protobuf:"varint,5,opt,name=createDate,proto3" json:"createDate,omitempty"`
	UpdateDate int64  `protobuf:"varint,6,opt,name=updateDate,proto3" json:"updateDate,omitempty"`
}

func (x *PaysReply_Pay) Reset() {
	*x = PaysReply_Pay{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaysReply_Pay) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaysReply_Pay) ProtoMessage() {}

func (x *PaysReply_Pay) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaysReply_Pay.ProtoReflect.Descriptor instead.
func (*PaysReply_Pay) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{3, 0}
}

func (x *PaysReply_Pay) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PaysReply_Pay) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PaysReply_Pay) GetOrder() int32 {
	if x != nil {
		return x.Order
	}
	return 0
}

func (x *PaysReply_Pay) GetPayStatus() int32 {
	if x != nil {
		return x.PayStatus
	}
	return 0
}

func (x *PaysReply_Pay) GetCreateDate() int64 {
	if x != nil {
		return x.CreateDate
	}
	return 0
}

func (x *PaysReply_Pay) GetUpdateDate() int64 {
	if x != nil {
		return x.UpdateDate
	}
	return 0
}

type PayReply_Pay struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Order      int32  `protobuf:"varint,3,opt,name=order,proto3" json:"order,omitempty"`
	PayStatus  int32  `protobuf:"varint,4,opt,name=payStatus,proto3" json:"payStatus,omitempty"`
	CreateDate int64  `protobuf:"varint,5,opt,name=createDate,proto3" json:"createDate,omitempty"`
	UpdateDate int64  `protobuf:"varint,6,opt,name=updateDate,proto3" json:"updateDate,omitempty"`
}

func (x *PayReply_Pay) Reset() {
	*x = PayReply_Pay{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayReply_Pay) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayReply_Pay) ProtoMessage() {}

func (x *PayReply_Pay) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayReply_Pay.ProtoReflect.Descriptor instead.
func (*PayReply_Pay) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{4, 0}
}

func (x *PayReply_Pay) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PayReply_Pay) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PayReply_Pay) GetOrder() int32 {
	if x != nil {
		return x.Order
	}
	return 0
}

func (x *PayReply_Pay) GetPayStatus() int32 {
	if x != nil {
		return x.PayStatus
	}
	return 0
}

func (x *PayReply_Pay) GetCreateDate() int64 {
	if x != nil {
		return x.CreateDate
	}
	return 0
}

func (x *PayReply_Pay) GetUpdateDate() int64 {
	if x != nil {
		return x.UpdateDate
	}
	return 0
}

var File_pay_proto protoreflect.FileDescriptor

var file_pay_proto_rawDesc = []byte{
	0x0a, 0x09, 0x70, 0x61, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x61, 0x79,
	0x2e, 0x76, 0x31, 0x22, 0x44, 0x0a, 0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x23, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x49, 0x0a, 0x05, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x22, 0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0xfb, 0x01, 0x0a, 0x09, 0x50, 0x61, 0x79, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x23,
	0x0a, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x70, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x52, 0x05, 0x72, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x29, 0x0a, 0x04, 0x70, 0x61, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x70, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x73, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x2e, 0x50, 0x61, 0x79, 0x52, 0x04, 0x70, 0x61, 0x79, 0x73, 0x1a, 0x9d,
	0x01, 0x0a, 0x03, 0x50, 0x61, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x70, 0x61, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x22, 0xf7,
	0x01, 0x0a, 0x08, 0x50, 0x61, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x23, 0x0a, 0x05, 0x72,
	0x65, 0x70, 0x6c, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x61, 0x79,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x52, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x26, 0x0a, 0x03, 0x70, 0x61, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x70, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e,
	0x50, 0x61, 0x79, 0x52, 0x03, 0x70, 0x61, 0x79, 0x1a, 0x9d, 0x01, 0x0a, 0x03, 0x50, 0x61, 0x79,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61,
	0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70,
	0x61, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x22, 0x0a, 0x0a, 0x08, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x52, 0x65, 0x71, 0x32, 0x70, 0x0a, 0x05, 0x50, 0x61, 0x79, 0x56, 0x31, 0x12, 0x32, 0x0a,
	0x0a, 0x47, 0x65, 0x74, 0x50, 0x61, 0x79, 0x42, 0x79, 0x49, 0x44, 0x12, 0x10, 0x2e, 0x70, 0x61,
	0x79, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e,
	0x70, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x12, 0x33, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x50, 0x61, 0x79, 0x73, 0x12,
	0x10, 0x2e, 0x70, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65,
	0x71, 0x1a, 0x11, 0x2e, 0x70, 0x61, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x73, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x70, 0x69, 0x63, 0x6b, 0x72, 0x65,
	0x77, 0x61, 0x72, 0x64, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x61, 0x79,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pay_proto_rawDescOnce sync.Once
	file_pay_proto_rawDescData = file_pay_proto_rawDesc
)

func file_pay_proto_rawDescGZIP() []byte {
	file_pay_proto_rawDescOnce.Do(func() {
		file_pay_proto_rawDescData = protoimpl.X.CompressGZIP(file_pay_proto_rawDescData)
	})
	return file_pay_proto_rawDescData
}

var file_pay_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pay_proto_goTypes = []interface{}{
	(*Reply)(nil),         // 0: pay.v1.Reply
	(*Error)(nil),         // 1: pay.v1.Error
	(*PayIDReq)(nil),      // 2: pay.v1.PayIDReq
	(*PaysReply)(nil),     // 3: pay.v1.PaysReply
	(*PayReply)(nil),      // 4: pay.v1.PayReply
	(*EmptyReq)(nil),      // 5: pay.v1.EmptyReq
	(*PaysReply_Pay)(nil), // 6: pay.v1.PaysReply.Pay
	(*PayReply_Pay)(nil),  // 7: pay.v1.PayReply.Pay
}
var file_pay_proto_depIdxs = []int32{
	1, // 0: pay.v1.Reply.error:type_name -> pay.v1.Error
	0, // 1: pay.v1.PaysReply.reply:type_name -> pay.v1.Reply
	6, // 2: pay.v1.PaysReply.pays:type_name -> pay.v1.PaysReply.Pay
	0, // 3: pay.v1.PayReply.reply:type_name -> pay.v1.Reply
	7, // 4: pay.v1.PayReply.pay:type_name -> pay.v1.PayReply.Pay
	2, // 5: pay.v1.PayV1.GetPayByID:input_type -> pay.v1.PayIDReq
	5, // 6: pay.v1.PayV1.GetAllPays:input_type -> pay.v1.EmptyReq
	4, // 7: pay.v1.PayV1.GetPayByID:output_type -> pay.v1.PayReply
	3, // 8: pay.v1.PayV1.GetAllPays:output_type -> pay.v1.PaysReply
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pay_proto_init() }
func file_pay_proto_init() {
	if File_pay_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pay_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayIDReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaysReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaysReply_Pay); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayReply_Pay); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pay_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pay_proto_goTypes,
		DependencyIndexes: file_pay_proto_depIdxs,
		MessageInfos:      file_pay_proto_msgTypes,
	}.Build()
	File_pay_proto = out.File
	file_pay_proto_rawDesc = nil
	file_pay_proto_goTypes = nil
	file_pay_proto_depIdxs = nil
}
