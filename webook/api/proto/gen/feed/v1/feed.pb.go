// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: feed/v1/feed.proto

package feedv1

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

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feed_v1_feed_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_feed_v1_feed_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_feed_v1_feed_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Article struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Article) Reset() {
	*x = Article{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feed_v1_feed_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Article) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Article) ProtoMessage() {}

func (x *Article) ProtoReflect() protoreflect.Message {
	mi := &file_feed_v1_feed_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Article.ProtoReflect.Descriptor instead.
func (*Article) Descriptor() ([]byte, []int) {
	return file_feed_v1_feed_proto_rawDescGZIP(), []int{1}
}

func (x *Article) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type FeedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	User    *User  `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Type    string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Content string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Ctime   int64  `protobuf:"varint,5,opt,name=ctime,proto3" json:"ctime,omitempty"`
}

func (x *FeedEvent) Reset() {
	*x = FeedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feed_v1_feed_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FeedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FeedEvent) ProtoMessage() {}

func (x *FeedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_feed_v1_feed_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FeedEvent.ProtoReflect.Descriptor instead.
func (*FeedEvent) Descriptor() ([]byte, []int) {
	return file_feed_v1_feed_proto_rawDescGZIP(), []int{2}
}

func (x *FeedEvent) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FeedEvent) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *FeedEvent) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *FeedEvent) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *FeedEvent) GetCtime() int64 {
	if x != nil {
		return x.Ctime
	}
	return 0
}

type CreateFeedEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FeedEvent *FeedEvent `protobuf:"bytes,1,opt,name=feedEvent,proto3" json:"feedEvent,omitempty"`
}

func (x *CreateFeedEventRequest) Reset() {
	*x = CreateFeedEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feed_v1_feed_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFeedEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFeedEventRequest) ProtoMessage() {}

func (x *CreateFeedEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_feed_v1_feed_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFeedEventRequest.ProtoReflect.Descriptor instead.
func (*CreateFeedEventRequest) Descriptor() ([]byte, []int) {
	return file_feed_v1_feed_proto_rawDescGZIP(), []int{3}
}

func (x *CreateFeedEventRequest) GetFeedEvent() *FeedEvent {
	if x != nil {
		return x.FeedEvent
	}
	return nil
}

type CreateFeedEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateFeedEventResponse) Reset() {
	*x = CreateFeedEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feed_v1_feed_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFeedEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFeedEventResponse) ProtoMessage() {}

func (x *CreateFeedEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_feed_v1_feed_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFeedEventResponse.ProtoReflect.Descriptor instead.
func (*CreateFeedEventResponse) Descriptor() ([]byte, []int) {
	return file_feed_v1_feed_proto_rawDescGZIP(), []int{4}
}

type FindFeedEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid       int64 `protobuf:"varint,1,opt,name=Uid,proto3" json:"Uid,omitempty"`
	Limit     int64 `protobuf:"varint,2,opt,name=Limit,proto3" json:"Limit,omitempty"`
	Timestamp int64 `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *FindFeedEventsRequest) Reset() {
	*x = FindFeedEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feed_v1_feed_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindFeedEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindFeedEventsRequest) ProtoMessage() {}

func (x *FindFeedEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_feed_v1_feed_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindFeedEventsRequest.ProtoReflect.Descriptor instead.
func (*FindFeedEventsRequest) Descriptor() ([]byte, []int) {
	return file_feed_v1_feed_proto_rawDescGZIP(), []int{5}
}

func (x *FindFeedEventsRequest) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *FindFeedEventsRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *FindFeedEventsRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type FindFeedEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FeedEvents []*FeedEvent `protobuf:"bytes,1,rep,name=feedEvents,proto3" json:"feedEvents,omitempty"`
}

func (x *FindFeedEventsResponse) Reset() {
	*x = FindFeedEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feed_v1_feed_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindFeedEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindFeedEventsResponse) ProtoMessage() {}

func (x *FindFeedEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_feed_v1_feed_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindFeedEventsResponse.ProtoReflect.Descriptor instead.
func (*FindFeedEventsResponse) Descriptor() ([]byte, []int) {
	return file_feed_v1_feed_proto_rawDescGZIP(), []int{6}
}

func (x *FindFeedEventsResponse) GetFeedEvents() []*FeedEvent {
	if x != nil {
		return x.FeedEvents
	}
	return nil
}

var File_feed_v1_feed_proto protoreflect.FileDescriptor

var file_feed_v1_feed_proto_rawDesc = []byte{
	0x0a, 0x12, 0x66, 0x65, 0x65, 0x64, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x65, 0x65, 0x64, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x66, 0x65, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x22, 0x16, 0x0a,
	0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x19, 0x0a, 0x07, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x82, 0x01, 0x0a, 0x09, 0x46, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x21,
	0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x66,
	0x65, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x63, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x4a, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46,
	0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x30, 0x0a, 0x09, 0x66, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65, 0x65,
	0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x66, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x22, 0x19, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x5d, 0x0a, 0x15,
	0x46, 0x69, 0x6e, 0x64, 0x46, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x03, 0x55, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x4c, 0x0a, 0x16, 0x46,
	0x69, 0x6e, 0x64, 0x46, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x0a, 0x66, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x65, 0x65, 0x64,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x66,
	0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x32, 0xb2, 0x01, 0x0a, 0x07, 0x46, 0x65,
	0x65, 0x64, 0x53, 0x76, 0x63, 0x12, 0x54, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46,
	0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x66, 0x65, 0x65, 0x64,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0e, 0x46,
	0x69, 0x6e, 0x64, 0x46, 0x65, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1e, 0x2e,
	0x66, 0x65, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x46, 0x65, 0x65, 0x64,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x66, 0x65, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x46, 0x65, 0x65, 0x64,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x96,
	0x01, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x42, 0x09,
	0x46, 0x65, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3f, 0x67, 0x69, 0x74,
	0x65, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x65, 0x65, 0x6b, 0x62, 0x61, 0x6e, 0x67, 0x2f,
	0x62, 0x61, 0x73, 0x69, 0x63, 0x2d, 0x67, 0x6f, 0x2f, 0x77, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x66, 0x65,
	0x65, 0x64, 0x2f, 0x76, 0x31, 0x3b, 0x66, 0x65, 0x65, 0x64, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x46,
	0x58, 0x58, 0xaa, 0x02, 0x07, 0x46, 0x65, 0x65, 0x64, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x07, 0x46,
	0x65, 0x65, 0x64, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x13, 0x46, 0x65, 0x65, 0x64, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x08, 0x46,
	0x65, 0x65, 0x64, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_feed_v1_feed_proto_rawDescOnce sync.Once
	file_feed_v1_feed_proto_rawDescData = file_feed_v1_feed_proto_rawDesc
)

func file_feed_v1_feed_proto_rawDescGZIP() []byte {
	file_feed_v1_feed_proto_rawDescOnce.Do(func() {
		file_feed_v1_feed_proto_rawDescData = protoimpl.X.CompressGZIP(file_feed_v1_feed_proto_rawDescData)
	})
	return file_feed_v1_feed_proto_rawDescData
}

var file_feed_v1_feed_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_feed_v1_feed_proto_goTypes = []interface{}{
	(*User)(nil),                    // 0: feed.v1.User
	(*Article)(nil),                 // 1: feed.v1.Article
	(*FeedEvent)(nil),               // 2: feed.v1.FeedEvent
	(*CreateFeedEventRequest)(nil),  // 3: feed.v1.CreateFeedEventRequest
	(*CreateFeedEventResponse)(nil), // 4: feed.v1.CreateFeedEventResponse
	(*FindFeedEventsRequest)(nil),   // 5: feed.v1.FindFeedEventsRequest
	(*FindFeedEventsResponse)(nil),  // 6: feed.v1.FindFeedEventsResponse
}
var file_feed_v1_feed_proto_depIdxs = []int32{
	0, // 0: feed.v1.FeedEvent.user:type_name -> feed.v1.User
	2, // 1: feed.v1.CreateFeedEventRequest.feedEvent:type_name -> feed.v1.FeedEvent
	2, // 2: feed.v1.FindFeedEventsResponse.feedEvents:type_name -> feed.v1.FeedEvent
	3, // 3: feed.v1.FeedSvc.CreateFeedEvent:input_type -> feed.v1.CreateFeedEventRequest
	5, // 4: feed.v1.FeedSvc.FindFeedEvents:input_type -> feed.v1.FindFeedEventsRequest
	4, // 5: feed.v1.FeedSvc.CreateFeedEvent:output_type -> feed.v1.CreateFeedEventResponse
	6, // 6: feed.v1.FeedSvc.FindFeedEvents:output_type -> feed.v1.FindFeedEventsResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_feed_v1_feed_proto_init() }
func file_feed_v1_feed_proto_init() {
	if File_feed_v1_feed_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_feed_v1_feed_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_feed_v1_feed_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Article); i {
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
		file_feed_v1_feed_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FeedEvent); i {
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
		file_feed_v1_feed_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateFeedEventRequest); i {
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
		file_feed_v1_feed_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateFeedEventResponse); i {
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
		file_feed_v1_feed_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindFeedEventsRequest); i {
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
		file_feed_v1_feed_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindFeedEventsResponse); i {
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
			RawDescriptor: file_feed_v1_feed_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_feed_v1_feed_proto_goTypes,
		DependencyIndexes: file_feed_v1_feed_proto_depIdxs,
		MessageInfos:      file_feed_v1_feed_proto_msgTypes,
	}.Build()
	File_feed_v1_feed_proto = out.File
	file_feed_v1_feed_proto_rawDesc = nil
	file_feed_v1_feed_proto_goTypes = nil
	file_feed_v1_feed_proto_depIdxs = nil
}
