// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.17.3
// source: protos/message.proto

package protos

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Enumeration of message types
type MessageMsgtype int32

const (
	Message_QUERY    MessageMsgtype = 0
	Message_ENTITY   MessageMsgtype = 1
	Message_RESPONSE MessageMsgtype = 2
)

// Enum value maps for MessageMsgtype.
var (
	MessageMsgtype_name = map[int32]string{
		0: "QUERY",
		1: "ENTITY",
		2: "RESPONSE",
	}
	MessageMsgtype_value = map[string]int32{
		"QUERY":    0,
		"ENTITY":   1,
		"RESPONSE": 2,
	}
)

func (x MessageMsgtype) Enum() *MessageMsgtype {
	p := new(MessageMsgtype)
	*p = x
	return p
}

func (x MessageMsgtype) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageMsgtype) Descriptor() protoreflect.EnumDescriptor {
	return file_protos_message_proto_enumTypes[0].Descriptor()
}

func (MessageMsgtype) Type() protoreflect.EnumType {
	return &file_protos_message_proto_enumTypes[0]
}

func (x MessageMsgtype) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageMsgtype.Descriptor instead.
func (MessageMsgtype) EnumDescriptor() ([]byte, []int) {
	return file_protos_message_proto_rawDescGZIP(), []int{0, 0}
}

// An arbitary message sent between two peers on the network.
// This message wraps any possible underlying message.
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Type of the message
	Type MessageMsgtype `protobuf:"varint,1,opt,name=type,proto3,enum=MessageMsgtype" json:"type,omitempty"`
	// PeerID of the sender
	Peerid string `protobuf:"bytes,2,opt,name=peerid,proto3" json:"peerid,omitempty"`
	// Body of the message
	//
	// Types that are assignable to Message:
	//	*Message_Query
	//	*Message_Entity
	//	*Message_Response
	Message isMessage_Message `protobuf_oneof:"message"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_protos_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_protos_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetType() MessageMsgtype {
	if x != nil {
		return x.Type
	}
	return Message_QUERY
}

func (x *Message) GetPeerid() string {
	if x != nil {
		return x.Peerid
	}
	return ""
}

func (m *Message) GetMessage() isMessage_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Message) GetQuery() *Query {
	if x, ok := x.GetMessage().(*Message_Query); ok {
		return x.Query
	}
	return nil
}

func (x *Message) GetEntity() *Entity {
	if x, ok := x.GetMessage().(*Message_Entity); ok {
		return x.Entity
	}
	return nil
}

func (x *Message) GetResponse() *Response {
	if x, ok := x.GetMessage().(*Message_Response); ok {
		return x.Response
	}
	return nil
}

type isMessage_Message interface {
	isMessage_Message()
}

type Message_Query struct {
	// Type must be QUERY
	Query *Query `protobuf:"bytes,3,opt,name=query,proto3,oneof"`
}

type Message_Entity struct {
	// Type must be ENTITY
	Entity *Entity `protobuf:"bytes,4,opt,name=entity,proto3,oneof"`
}

type Message_Response struct {
	// Type must be RESPONSE
	Response *Response `protobuf:"bytes,5,opt,name=response,proto3,oneof"`
}

func (*Message_Query) isMessage_Message() {}

func (*Message_Entity) isMessage_Message() {}

func (*Message_Response) isMessage_Message() {}

var File_protos_message_proto protoreflect.FileDescriptor

var file_protos_message_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xee, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x10, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x6d, 0x73, 0x67, 0x74, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x65, 0x72,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x69, 0x64,
	0x12, 0x1e, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x06, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48, 0x00, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x12, 0x21, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x07, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x48, 0x00, 0x52, 0x06, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x12, 0x27, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2e, 0x0a, 0x07,
	0x6d, 0x73, 0x67, 0x74, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x51, 0x55, 0x45, 0x52, 0x59,
	0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x45, 0x4e, 0x54, 0x49, 0x54, 0x59, 0x10, 0x01, 0x12, 0x0c,
	0x0a, 0x08, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x02, 0x42, 0x09, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_message_proto_rawDescOnce sync.Once
	file_protos_message_proto_rawDescData = file_protos_message_proto_rawDesc
)

func file_protos_message_proto_rawDescGZIP() []byte {
	file_protos_message_proto_rawDescOnce.Do(func() {
		file_protos_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_message_proto_rawDescData)
	})
	return file_protos_message_proto_rawDescData
}

var file_protos_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protos_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_protos_message_proto_goTypes = []interface{}{
	(MessageMsgtype)(0), // 0: Message.msgtype
	(*Message)(nil),     // 1: Message
	(*Query)(nil),       // 2: Query
	(*Entity)(nil),      // 3: Entity
	(*Response)(nil),    // 4: Response
}
var file_protos_message_proto_depIdxs = []int32{
	0, // 0: Message.type:type_name -> Message.msgtype
	2, // 1: Message.query:type_name -> Query
	3, // 2: Message.entity:type_name -> Entity
	4, // 3: Message.response:type_name -> Response
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_protos_message_proto_init() }
func file_protos_message_proto_init() {
	if File_protos_message_proto != nil {
		return
	}
	file_protos_entity_proto_init()
	file_protos_query_proto_init()
	file_protos_response_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_protos_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
	file_protos_message_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_Query)(nil),
		(*Message_Entity)(nil),
		(*Message_Response)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_message_proto_goTypes,
		DependencyIndexes: file_protos_message_proto_depIdxs,
		EnumInfos:         file_protos_message_proto_enumTypes,
		MessageInfos:      file_protos_message_proto_msgTypes,
	}.Build()
	File_protos_message_proto = out.File
	file_protos_message_proto_rawDesc = nil
	file_protos_message_proto_goTypes = nil
	file_protos_message_proto_depIdxs = nil
}
