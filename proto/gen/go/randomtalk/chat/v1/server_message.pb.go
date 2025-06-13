// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: randomtalk/chat/v1/server_message.proto

package chatpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Kind int32

const (
	Kind_KIND_UNSPECIFIED Kind = 0
	Kind_KIND_SYSTEM      Kind = 1 // System message
	Kind_KIND_USER        Kind = 2 // User message
)

// Enum value maps for Kind.
var (
	Kind_name = map[int32]string{
		0: "KIND_UNSPECIFIED",
		1: "KIND_SYSTEM",
		2: "KIND_USER",
	}
	Kind_value = map[string]int32{
		"KIND_UNSPECIFIED": 0,
		"KIND_SYSTEM":      1,
		"KIND_USER":        2,
	}
)

func (x Kind) Enum() *Kind {
	p := new(Kind)
	*p = x
	return p
}

func (x Kind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Kind) Descriptor() protoreflect.EnumDescriptor {
	return file_randomtalk_chat_v1_server_message_proto_enumTypes[0].Descriptor()
}

func (Kind) Type() protoreflect.EnumType {
	return &file_randomtalk_chat_v1_server_message_proto_enumTypes[0]
}

func (x Kind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Kind.Descriptor instead.
func (Kind) EnumDescriptor() ([]byte, []int) {
	return file_randomtalk_chat_v1_server_message_proto_rawDescGZIP(), []int{0}
}

// ServerMessage represents a message sent from the server to the client.
type ServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind Kind `protobuf:"varint,1,opt,name=kind,proto3,enum=randomtalk.chat.v1.Kind" json:"kind,omitempty"`
	// Types that are assignable to Data:
	//
	//	*ServerMessage_Command
	//	*ServerMessage_Error
	//	*ServerMessage_Info
	//	*ServerMessage_Notification
	Data isServerMessage_Data `protobuf_oneof:"data"`
}

func (x *ServerMessage) Reset() {
	*x = ServerMessage{}
	mi := &file_randomtalk_chat_v1_server_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage) ProtoMessage() {}

func (x *ServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_randomtalk_chat_v1_server_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage.ProtoReflect.Descriptor instead.
func (*ServerMessage) Descriptor() ([]byte, []int) {
	return file_randomtalk_chat_v1_server_message_proto_rawDescGZIP(), []int{0}
}

func (x *ServerMessage) GetKind() Kind {
	if x != nil {
		return x.Kind
	}
	return Kind_KIND_UNSPECIFIED
}

func (m *ServerMessage) GetData() isServerMessage_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *ServerMessage) GetCommand() *Command {
	if x, ok := x.GetData().(*ServerMessage_Command); ok {
		return x.Command
	}
	return nil
}

func (x *ServerMessage) GetError() *ErrorMessage {
	if x, ok := x.GetData().(*ServerMessage_Error); ok {
		return x.Error
	}
	return nil
}

func (x *ServerMessage) GetInfo() *InfoMessage {
	if x, ok := x.GetData().(*ServerMessage_Info); ok {
		return x.Info
	}
	return nil
}

func (x *ServerMessage) GetNotification() *NotificationMessage {
	if x, ok := x.GetData().(*ServerMessage_Notification); ok {
		return x.Notification
	}
	return nil
}

type isServerMessage_Data interface {
	isServerMessage_Data()
}

type ServerMessage_Command struct {
	Command *Command `protobuf:"bytes,2,opt,name=command,proto3,oneof"`
}

type ServerMessage_Error struct {
	Error *ErrorMessage `protobuf:"bytes,3,opt,name=error,proto3,oneof"`
}

type ServerMessage_Info struct {
	Info *InfoMessage `protobuf:"bytes,4,opt,name=info,proto3,oneof"`
}

type ServerMessage_Notification struct {
	Notification *NotificationMessage `protobuf:"bytes,5,opt,name=notification,proto3,oneof"`
}

func (*ServerMessage_Command) isServerMessage_Data() {}

func (*ServerMessage_Error) isServerMessage_Data() {}

func (*ServerMessage_Info) isServerMessage_Data() {}

func (*ServerMessage_Notification) isServerMessage_Data() {}

var File_randomtalk_chat_v1_server_message_proto protoreflect.FileDescriptor

var file_randomtalk_chat_v1_server_message_proto_rawDesc = []byte{
	0x0a, 0x27, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x63, 0x68, 0x61,
	0x74, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x72, 0x61, 0x6e, 0x64, 0x6f,
	0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20,
	0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x26, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x63, 0x68, 0x61,
	0x74, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x66,
	0x6f, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2d, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x63, 0x68, 0x61, 0x74,
	0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbe,
	0x02, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x2c, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18,
	0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x37,
	0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x07,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x38, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74,
	0x61, 0x6c, 0x6b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x12, 0x35, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1f, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x48, 0x00, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x4d, 0x0a, 0x0c, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27,
	0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a,
	0x3c, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x10, 0x4b, 0x49, 0x4e, 0x44, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a,
	0x0b, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x10, 0x02, 0x42, 0x2c, 0x5a,
	0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x66, 0x72, 0x72,
	0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_randomtalk_chat_v1_server_message_proto_rawDescOnce sync.Once
	file_randomtalk_chat_v1_server_message_proto_rawDescData = file_randomtalk_chat_v1_server_message_proto_rawDesc
)

func file_randomtalk_chat_v1_server_message_proto_rawDescGZIP() []byte {
	file_randomtalk_chat_v1_server_message_proto_rawDescOnce.Do(func() {
		file_randomtalk_chat_v1_server_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_randomtalk_chat_v1_server_message_proto_rawDescData)
	})
	return file_randomtalk_chat_v1_server_message_proto_rawDescData
}

var file_randomtalk_chat_v1_server_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_randomtalk_chat_v1_server_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_randomtalk_chat_v1_server_message_proto_goTypes = []any{
	(Kind)(0),                   // 0: randomtalk.chat.v1.Kind
	(*ServerMessage)(nil),       // 1: randomtalk.chat.v1.ServerMessage
	(*Command)(nil),             // 2: randomtalk.chat.v1.Command
	(*ErrorMessage)(nil),        // 3: randomtalk.chat.v1.ErrorMessage
	(*InfoMessage)(nil),         // 4: randomtalk.chat.v1.InfoMessage
	(*NotificationMessage)(nil), // 5: randomtalk.chat.v1.NotificationMessage
}
var file_randomtalk_chat_v1_server_message_proto_depIdxs = []int32{
	0, // 0: randomtalk.chat.v1.ServerMessage.kind:type_name -> randomtalk.chat.v1.Kind
	2, // 1: randomtalk.chat.v1.ServerMessage.command:type_name -> randomtalk.chat.v1.Command
	3, // 2: randomtalk.chat.v1.ServerMessage.error:type_name -> randomtalk.chat.v1.ErrorMessage
	4, // 3: randomtalk.chat.v1.ServerMessage.info:type_name -> randomtalk.chat.v1.InfoMessage
	5, // 4: randomtalk.chat.v1.ServerMessage.notification:type_name -> randomtalk.chat.v1.NotificationMessage
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_randomtalk_chat_v1_server_message_proto_init() }
func file_randomtalk_chat_v1_server_message_proto_init() {
	if File_randomtalk_chat_v1_server_message_proto != nil {
		return
	}
	file_randomtalk_chat_v1_command_proto_init()
	file_randomtalk_chat_v1_error_message_proto_init()
	file_randomtalk_chat_v1_info_message_proto_init()
	file_randomtalk_chat_v1_notification_message_proto_init()
	file_randomtalk_chat_v1_server_message_proto_msgTypes[0].OneofWrappers = []any{
		(*ServerMessage_Command)(nil),
		(*ServerMessage_Error)(nil),
		(*ServerMessage_Info)(nil),
		(*ServerMessage_Notification)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_randomtalk_chat_v1_server_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_randomtalk_chat_v1_server_message_proto_goTypes,
		DependencyIndexes: file_randomtalk_chat_v1_server_message_proto_depIdxs,
		EnumInfos:         file_randomtalk_chat_v1_server_message_proto_enumTypes,
		MessageInfos:      file_randomtalk_chat_v1_server_message_proto_msgTypes,
	}.Build()
	File_randomtalk_chat_v1_server_message_proto = out.File
	file_randomtalk_chat_v1_server_message_proto_rawDesc = nil
	file_randomtalk_chat_v1_server_message_proto_goTypes = nil
	file_randomtalk_chat_v1_server_message_proto_depIdxs = nil
}
