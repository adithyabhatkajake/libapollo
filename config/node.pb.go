// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: config/node.proto

package config

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

// All configuration aggregated here
type NodeDataConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Protocol specific configuration
	ProtConfig *ProtocolConfig `protobuf:"bytes,1,opt,name=ProtConfig,proto3" json:"ProtConfig,omitempty"`
	// Network configuration for nodes
	NetConfig *NetConfig `protobuf:"bytes,2,opt,name=NetConfig,proto3" json:"NetConfig,omitempty"`
	// Network configuration for clients
	ClientPort string `protobuf:"bytes,3,opt,name=ClientPort,proto3" json:"ClientPort,omitempty"`
	// Cryptographic configuration
	CryptoCon *CryptoConfig `protobuf:"bytes,4,opt,name=CryptoCon,proto3" json:"CryptoCon,omitempty"`
}

func (x *NodeDataConfig) Reset() {
	*x = NodeDataConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeDataConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeDataConfig) ProtoMessage() {}

func (x *NodeDataConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeDataConfig.ProtoReflect.Descriptor instead.
func (*NodeDataConfig) Descriptor() ([]byte, []int) {
	return file_config_node_proto_rawDescGZIP(), []int{0}
}

func (x *NodeDataConfig) GetProtConfig() *ProtocolConfig {
	if x != nil {
		return x.ProtConfig
	}
	return nil
}

func (x *NodeDataConfig) GetNetConfig() *NetConfig {
	if x != nil {
		return x.NetConfig
	}
	return nil
}

func (x *NodeDataConfig) GetClientPort() string {
	if x != nil {
		return x.ClientPort
	}
	return ""
}

func (x *NodeDataConfig) GetCryptoCon() *CryptoConfig {
	if x != nil {
		return x.CryptoCon
	}
	return nil
}

var File_config_node_proto protoreflect.FileDescriptor

var file_config_node_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x10, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xcd, 0x01, 0x0a, 0x0e, 0x4e, 0x6f, 0x64, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x36, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x52, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x2f,
	0x0a, 0x09, 0x4e, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4e, 0x65, 0x74, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x09, 0x4e, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x1e, 0x0a, 0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12,
	0x32, 0x0a, 0x09, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x43, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x72, 0x79, 0x70,
	0x74, 0x6f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x09, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f,
	0x43, 0x6f, 0x6e, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x64, 0x69, 0x74, 0x68, 0x79, 0x61, 0x62, 0x68, 0x61, 0x74, 0x6b, 0x61, 0x6a,
	0x61, 0x6b, 0x65, 0x2f, 0x6c, 0x69, 0x62, 0x61, 0x70, 0x6f, 0x6c, 0x6c, 0x6f, 0x2f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_node_proto_rawDescOnce sync.Once
	file_config_node_proto_rawDescData = file_config_node_proto_rawDesc
)

func file_config_node_proto_rawDescGZIP() []byte {
	file_config_node_proto_rawDescOnce.Do(func() {
		file_config_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_node_proto_rawDescData)
	})
	return file_config_node_proto_rawDescData
}

var file_config_node_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_config_node_proto_goTypes = []interface{}{
	(*NodeDataConfig)(nil), // 0: config.NodeDataConfig
	(*ProtocolConfig)(nil), // 1: config.ProtocolConfig
	(*NetConfig)(nil),      // 2: config.NetConfig
	(*CryptoConfig)(nil),   // 3: config.CryptoConfig
}
var file_config_node_proto_depIdxs = []int32{
	1, // 0: config.NodeDataConfig.ProtConfig:type_name -> config.ProtocolConfig
	2, // 1: config.NodeDataConfig.NetConfig:type_name -> config.NetConfig
	3, // 2: config.NodeDataConfig.CryptoCon:type_name -> config.CryptoConfig
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_config_node_proto_init() }
func file_config_node_proto_init() {
	if File_config_node_proto != nil {
		return
	}
	file_cryptoconf_proto_init()
	file_config_protocol_proto_init()
	file_network_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_config_node_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeDataConfig); i {
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
			RawDescriptor: file_config_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_node_proto_goTypes,
		DependencyIndexes: file_config_node_proto_depIdxs,
		MessageInfos:      file_config_node_proto_msgTypes,
	}.Build()
	File_config_node_proto = out.File
	file_config_node_proto_rawDesc = nil
	file_config_node_proto_goTypes = nil
	file_config_node_proto_depIdxs = nil
}
