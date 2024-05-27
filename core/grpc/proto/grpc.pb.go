// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.3
// source: grpc.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MetaData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Leader  string `protobuf:"bytes,1,opt,name=leader,proto3" json:"leader,omitempty"`
	SeqNum  int64  `protobuf:"varint,2,opt,name=seqNum,proto3" json:"seqNum,omitempty"`
	ReqID   string `protobuf:"bytes,3,opt,name=reqID,proto3" json:"reqID,omitempty"`
	Timeout []byte `protobuf:"bytes,4,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Sig     []byte `protobuf:"bytes,5,opt,name=sig,proto3" json:"sig,omitempty"` // leader.Sign(hash({leaderAddr, seqNum, reqID, timeout}))
}

func (x *MetaData) Reset() {
	*x = MetaData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaData) ProtoMessage() {}

func (x *MetaData) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaData.ProtoReflect.Descriptor instead.
func (*MetaData) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{0}
}

func (x *MetaData) GetLeader() string {
	if x != nil {
		return x.Leader
	}
	return ""
}

func (x *MetaData) GetSeqNum() int64 {
	if x != nil {
		return x.SeqNum
	}
	return 0
}

func (x *MetaData) GetReqID() string {
	if x != nil {
		return x.ReqID
	}
	return ""
}

func (x *MetaData) GetTimeout() []byte {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *MetaData) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

// client -> leader
type RequestMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SeqNum  int64  `protobuf:"varint,1,opt,name=seqNum,proto3" json:"seqNum,omitempty"`
	ReqID   string `protobuf:"bytes,2,opt,name=reqID,proto3" json:"reqID,omitempty"`
	Timeout []byte `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *RequestMessage) Reset() {
	*x = RequestMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestMessage) ProtoMessage() {}

func (x *RequestMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestMessage.ProtoReflect.Descriptor instead.
func (*RequestMessage) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{1}
}

func (x *RequestMessage) GetSeqNum() int64 {
	if x != nil {
		return x.SeqNum
	}
	return 0
}

func (x *RequestMessage) GetReqID() string {
	if x != nil {
		return x.ReqID
	}
	return ""
}

func (x *RequestMessage) GetTimeout() []byte {
	if x != nil {
		return x.Timeout
	}
	return nil
}

// leader->generator || generator->generator || generator->leader
type TransportMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	MsgType string `protobuf:"bytes,2,opt,name=msgType,proto3" json:"msgType,omitempty"` // Start/Single/Collection/Merged
	Data    []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`       // StartNewRoundMessage/SingleRandMessage/CollectionRandMessage/MergedRandMessage
}

func (x *TransportMessage) Reset() {
	*x = TransportMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransportMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransportMessage) ProtoMessage() {}

func (x *TransportMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransportMessage.ProtoReflect.Descriptor instead.
func (*TransportMessage) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{2}
}

func (x *TransportMessage) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *TransportMessage) GetMsgType() string {
	if x != nil {
		return x.MsgType
	}
	return ""
}

func (x *TransportMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// TransportMessage.data, leader -> generator
type StartNewRoundMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetaData *MetaData `protobuf:"bytes,1,opt,name=metaData,proto3" json:"metaData,omitempty"`
	From     string    `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
}

func (x *StartNewRoundMessage) Reset() {
	*x = StartNewRoundMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartNewRoundMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartNewRoundMessage) ProtoMessage() {}

func (x *StartNewRoundMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartNewRoundMessage.ProtoReflect.Descriptor instead.
func (*StartNewRoundMessage) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{3}
}

func (x *StartNewRoundMessage) GetMetaData() *MetaData {
	if x != nil {
		return x.MetaData
	}
	return nil
}

func (x *StartNewRoundMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

// TransportMessage.data, generator -> generator
type SingleRandMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetaData  *MetaData `protobuf:"bytes,1,opt,name=metaData,proto3" json:"metaData,omitempty"`
	From      string    `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	RandomNum uint64    `protobuf:"varint,3,opt,name=randomNum,proto3" json:"randomNum,omitempty"`
	Sig       []byte    `protobuf:"bytes,4,opt,name=sig,proto3" json:"sig,omitempty"`
}

func (x *SingleRandMessage) Reset() {
	*x = SingleRandMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingleRandMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingleRandMessage) ProtoMessage() {}

func (x *SingleRandMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingleRandMessage.ProtoReflect.Descriptor instead.
func (*SingleRandMessage) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{4}
}

func (x *SingleRandMessage) GetMetaData() *MetaData {
	if x != nil {
		return x.MetaData
	}
	return nil
}

func (x *SingleRandMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *SingleRandMessage) GetRandomNum() uint64 {
	if x != nil {
		return x.RandomNum
	}
	return 0
}

func (x *SingleRandMessage) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

// TransportMessage.data, generator -> leader
type CollectionRandMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetaData                *MetaData         `protobuf:"bytes,1,opt,name=metaData,proto3" json:"metaData,omitempty"`
	From                    string            `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	CollectionRandomNumbers map[string]uint64 `protobuf:"bytes,3,rep,name=collectionRandomNumbers,proto3" json:"collectionRandomNumbers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"` // node:randNum
	Sig                     []byte            `protobuf:"bytes,4,opt,name=sig,proto3" json:"sig,omitempty"`
	SigPlaintext            []byte            `protobuf:"bytes,5,opt,name=sigPlaintext,proto3" json:"sigPlaintext,omitempty"`
}

func (x *CollectionRandMessage) Reset() {
	*x = CollectionRandMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CollectionRandMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CollectionRandMessage) ProtoMessage() {}

func (x *CollectionRandMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CollectionRandMessage.ProtoReflect.Descriptor instead.
func (*CollectionRandMessage) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{5}
}

func (x *CollectionRandMessage) GetMetaData() *MetaData {
	if x != nil {
		return x.MetaData
	}
	return nil
}

func (x *CollectionRandMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *CollectionRandMessage) GetCollectionRandomNumbers() map[string]uint64 {
	if x != nil {
		return x.CollectionRandomNumbers
	}
	return nil
}

func (x *CollectionRandMessage) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

func (x *CollectionRandMessage) GetSigPlaintext() []byte {
	if x != nil {
		return x.SigPlaintext
	}
	return nil
}

// TransportMessage.data, leader -> contract
type MergedRandMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetaData            *MetaData `protobuf:"bytes,1,opt,name=metaData,proto3" json:"metaData,omitempty"`
	From                string    `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	MergedRandomNumbers []uint64  `protobuf:"varint,3,rep,packed,name=mergedRandomNumbers,proto3" json:"mergedRandomNumbers,omitempty"`
	Sig                 []byte    `protobuf:"bytes,4,opt,name=sig,proto3" json:"sig,omitempty"`
}

func (x *MergedRandMessage) Reset() {
	*x = MergedRandMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MergedRandMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MergedRandMessage) ProtoMessage() {}

func (x *MergedRandMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MergedRandMessage.ProtoReflect.Descriptor instead.
func (*MergedRandMessage) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{6}
}

func (x *MergedRandMessage) GetMetaData() *MetaData {
	if x != nil {
		return x.MetaData
	}
	return nil
}

func (x *MergedRandMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *MergedRandMessage) GetMergedRandomNumbers() []uint64 {
	if x != nil {
		return x.MergedRandomNumbers
	}
	return nil
}

func (x *MergedRandMessage) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

type UnaryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From    string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	ReqID   string `protobuf:"bytes,2,opt,name=reqID,proto3" json:"reqID,omitempty"`
	SeqNum  int64  `protobuf:"varint,3,opt,name=seqNum,proto3" json:"seqNum,omitempty"`
	Message string `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UnaryResponse) Reset() {
	*x = UnaryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnaryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnaryResponse) ProtoMessage() {}

func (x *UnaryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnaryResponse.ProtoReflect.Descriptor instead.
func (*UnaryResponse) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{7}
}

func (x *UnaryResponse) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *UnaryResponse) GetReqID() string {
	if x != nil {
		return x.ReqID
	}
	return ""
}

func (x *UnaryResponse) GetSeqNum() int64 {
	if x != nil {
		return x.SeqNum
	}
	return 0
}

func (x *UnaryResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_grpc_proto protoreflect.FileDescriptor

var file_grpc_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x7c, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06,
	0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x71, 0x4e, 0x75, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x65, 0x71, 0x4e, 0x75, 0x6d, 0x12, 0x14, 0x0a, 0x05,
	0x72, 0x65, 0x71, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x71,
	0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x73, 0x69, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x69, 0x67, 0x22, 0x58,
	0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x71, 0x4e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x73, 0x65, 0x71, 0x4e, 0x75, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x71, 0x49,
	0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x71, 0x49, 0x44, 0x12, 0x18,
	0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x5a, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x57, 0x0a, 0x14, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4e, 0x65, 0x77,
	0x52, 0x6f, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x22, 0x84, 0x01,
	0x0a, 0x11, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65,
	0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x66, 0x72, 0x6f, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75,
	0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e,
	0x75, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x03, 0x73, 0x69, 0x67, 0x22, 0xcf, 0x02, 0x0a, 0x15, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x61, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12,
	0x73, 0x0a, 0x17, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x6e,
	0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x39, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x61, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x43,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x17, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x03, 0x73, 0x69, 0x67, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x69, 0x67, 0x50, 0x6c, 0x61,
	0x69, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x73, 0x69,
	0x67, 0x50, 0x6c, 0x61, 0x69, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x1a, 0x4a, 0x0a, 0x1c, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x98, 0x01, 0x0a, 0x11, 0x4d, 0x65, 0x72, 0x67, 0x65,
	0x64, 0x52, 0x61, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x30, 0x0a,
	0x13, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x64, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x04, 0x52, 0x13, 0x6d, 0x65, 0x72, 0x67,
	0x65, 0x64, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x69,
	0x67, 0x22, 0x6b, 0x0a, 0x0d, 0x55, 0x6e, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x71, 0x49, 0x44, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x71, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x65, 0x71, 0x4e, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x65,
	0x71, 0x4e, 0x75, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xa2,
	0x01, 0x0a, 0x15, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x43, 0x0a, 0x13, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x55, 0x6e, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a,
	0x13, 0x53, 0x65, 0x6e, 0x64, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x14, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x6e, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_proto_rawDescOnce sync.Once
	file_grpc_proto_rawDescData = file_grpc_proto_rawDesc
)

func file_grpc_proto_rawDescGZIP() []byte {
	file_grpc_proto_rawDescOnce.Do(func() {
		file_grpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_proto_rawDescData)
	})
	return file_grpc_proto_rawDescData
}

var file_grpc_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_grpc_proto_goTypes = []interface{}{
	(*MetaData)(nil),              // 0: proto.MetaData
	(*RequestMessage)(nil),        // 1: proto.RequestMessage
	(*TransportMessage)(nil),      // 2: proto.TransportMessage
	(*StartNewRoundMessage)(nil),  // 3: proto.StartNewRoundMessage
	(*SingleRandMessage)(nil),     // 4: proto.SingleRandMessage
	(*CollectionRandMessage)(nil), // 5: proto.CollectionRandMessage
	(*MergedRandMessage)(nil),     // 6: proto.MergedRandMessage
	(*UnaryResponse)(nil),         // 7: proto.UnaryResponse
	nil,                           // 8: proto.CollectionRandMessage.CollectionRandomNumbersEntry
	(*emptypb.Empty)(nil),         // 9: google.protobuf.Empty
}
var file_grpc_proto_depIdxs = []int32{
	0, // 0: proto.StartNewRoundMessage.metaData:type_name -> proto.MetaData
	0, // 1: proto.SingleRandMessage.metaData:type_name -> proto.MetaData
	0, // 2: proto.CollectionRandMessage.metaData:type_name -> proto.MetaData
	8, // 3: proto.CollectionRandMessage.collectionRandomNumbers:type_name -> proto.CollectionRandMessage.CollectionRandomNumbersEntry
	0, // 4: proto.MergedRandMessage.metaData:type_name -> proto.MetaData
	9, // 5: proto.RandomNumberGenerator.RequestRandomNumber:input_type -> google.protobuf.Empty
	2, // 6: proto.RandomNumberGenerator.SendInternalMessage:input_type -> proto.TransportMessage
	7, // 7: proto.RandomNumberGenerator.RequestRandomNumber:output_type -> proto.UnaryResponse
	7, // 8: proto.RandomNumberGenerator.SendInternalMessage:output_type -> proto.UnaryResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_grpc_proto_init() }
func file_grpc_proto_init() {
	if File_grpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetaData); i {
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
		file_grpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestMessage); i {
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
		file_grpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransportMessage); i {
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
		file_grpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartNewRoundMessage); i {
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
		file_grpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SingleRandMessage); i {
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
		file_grpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CollectionRandMessage); i {
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
		file_grpc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MergedRandMessage); i {
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
		file_grpc_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnaryResponse); i {
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
			RawDescriptor: file_grpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_proto_goTypes,
		DependencyIndexes: file_grpc_proto_depIdxs,
		MessageInfos:      file_grpc_proto_msgTypes,
	}.Build()
	File_grpc_proto = out.File
	file_grpc_proto_rawDesc = nil
	file_grpc_proto_goTypes = nil
	file_grpc_proto_depIdxs = nil
}
