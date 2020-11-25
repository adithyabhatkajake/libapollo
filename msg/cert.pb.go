// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: msg/cert.proto

package msg

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

type Certificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The data on which we need signatures
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	// The raw signatures
	Signatures [][]byte `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
	// The IDs of the signers
	Ids []uint64 `protobuf:"varint,3,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

func (x *Certificate) Reset() {
	*x = Certificate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_cert_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Certificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Certificate) ProtoMessage() {}

func (x *Certificate) ProtoReflect() protoreflect.Message {
	mi := &file_msg_cert_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Certificate.ProtoReflect.Descriptor instead.
func (*Certificate) Descriptor() ([]byte, []int) {
	return file_msg_cert_proto_rawDescGZIP(), []int{0}
}

func (x *Certificate) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Certificate) GetSignatures() [][]byte {
	if x != nil {
		return x.Signatures
	}
	return nil
}

func (x *Certificate) GetIds() []uint64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

var File_msg_cert_proto protoreflect.FileDescriptor

var file_msg_cert_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x73, 0x67, 0x2f, 0x63, 0x65, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x53, 0x0a, 0x0b, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0a, 0x73, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x04, 0x52, 0x03, 0x69, 0x64, 0x73, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x64, 0x69, 0x74, 0x68, 0x79, 0x61,
	0x62, 0x68, 0x61, 0x74, 0x6b, 0x61, 0x6a, 0x61, 0x6b, 0x65, 0x2f, 0x6c, 0x69, 0x62, 0x61, 0x70,
	0x6f, 0x6c, 0x6c, 0x6f, 0x2f, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_cert_proto_rawDescOnce sync.Once
	file_msg_cert_proto_rawDescData = file_msg_cert_proto_rawDesc
)

func file_msg_cert_proto_rawDescGZIP() []byte {
	file_msg_cert_proto_rawDescOnce.Do(func() {
		file_msg_cert_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_cert_proto_rawDescData)
	})
	return file_msg_cert_proto_rawDescData
}

var file_msg_cert_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_msg_cert_proto_goTypes = []interface{}{
	(*Certificate)(nil), // 0: msg.Certificate
}
var file_msg_cert_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_msg_cert_proto_init() }
func file_msg_cert_proto_init() {
	if File_msg_cert_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_cert_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Certificate); i {
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
			RawDescriptor: file_msg_cert_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_cert_proto_goTypes,
		DependencyIndexes: file_msg_cert_proto_depIdxs,
		MessageInfos:      file_msg_cert_proto_msgTypes,
	}.Build()
	File_msg_cert_proto = out.File
	file_msg_cert_proto_rawDesc = nil
	file_msg_cert_proto_goTypes = nil
	file_msg_cert_proto_depIdxs = nil
}