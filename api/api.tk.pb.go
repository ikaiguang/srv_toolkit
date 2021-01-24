// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: api.tk.proto

package tkpb

import (
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
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

// Platform 平台标识
type Platform int32

const (
	Platform_platform_unknown Platform = 0 // 未知
	// 平台
	Platform_platform_computer Platform = 1 // 电脑端
	Platform_platform_mobile   Platform = 2 // 移动端
	// 设备
	Platform_platform_desktop Platform = 20 // 桌面
	Platform_platform_android Platform = 21 // 安卓
	Platform_platform_iphone  Platform = 22 // 苹果
	Platform_platform_pad     Platform = 23 // 平板
)

// Enum value maps for Platform.
var (
	Platform_name = map[int32]string{
		0:  "platform_unknown",
		1:  "platform_computer",
		2:  "platform_mobile",
		20: "platform_desktop",
		21: "platform_android",
		22: "platform_iphone",
		23: "platform_pad",
	}
	Platform_value = map[string]int32{
		"platform_unknown":  0,
		"platform_computer": 1,
		"platform_mobile":   2,
		"platform_desktop":  20,
		"platform_android":  21,
		"platform_iphone":   22,
		"platform_pad":      23,
	}
)

func (x Platform) Enum() *Platform {
	p := new(Platform)
	*p = x
	return p
}

func (x Platform) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Platform) Descriptor() protoreflect.EnumDescriptor {
	return file_api_tk_proto_enumTypes[0].Descriptor()
}

func (Platform) Type() protoreflect.EnumType {
	return &file_api_tk_proto_enumTypes[0]
}

func (x Platform) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Platform.Descriptor instead.
func (Platform) EnumDescriptor() ([]byte, []int) {
	return file_api_tk_proto_rawDescGZIP(), []int{0}
}

// Response resp with google/protobuf/any.proto
type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`    // code
	Msg    string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`       // 响应描述
	Detail string   `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"` // 消息详情
	Data   *any.Any `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`     // data
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_tk_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_api_tk_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_api_tk_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Response) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Response) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *Response) GetData() *any.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_api_tk_proto protoreflect.FileDescriptor

var file_api_tk_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04,
	0x74, 0x6b, 0x70, 0x62, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x72, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x2a, 0x9f, 0x01, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x12, 0x14, 0x0a, 0x10, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x75, 0x6e, 0x6b,
	0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x72, 0x10, 0x01, 0x12, 0x13, 0x0a,
	0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65,
	0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x64,
	0x65, 0x73, 0x6b, 0x74, 0x6f, 0x70, 0x10, 0x14, 0x12, 0x14, 0x0a, 0x10, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x61, 0x6e, 0x64, 0x72, 0x6f, 0x69, 0x64, 0x10, 0x15, 0x12, 0x13,
	0x0a, 0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x69, 0x70, 0x68, 0x6f, 0x6e,
	0x65, 0x10, 0x16, 0x12, 0x10, 0x0a, 0x0c, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f,
	0x70, 0x61, 0x64, 0x10, 0x17, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6b, 0x61, 0x69, 0x67, 0x75, 0x61, 0x6e, 0x67, 0x2f, 0x73, 0x72,
	0x76, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x6b, 0x69, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x74, 0x6b,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_tk_proto_rawDescOnce sync.Once
	file_api_tk_proto_rawDescData = file_api_tk_proto_rawDesc
)

func file_api_tk_proto_rawDescGZIP() []byte {
	file_api_tk_proto_rawDescOnce.Do(func() {
		file_api_tk_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_tk_proto_rawDescData)
	})
	return file_api_tk_proto_rawDescData
}

var file_api_tk_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_tk_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_tk_proto_goTypes = []interface{}{
	(Platform)(0),    // 0: tkpb.Platform
	(*Response)(nil), // 1: tkpb.Response
	(*any.Any)(nil),  // 2: google.protobuf.Any
}
var file_api_tk_proto_depIdxs = []int32{
	2, // 0: tkpb.Response.data:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_tk_proto_init() }
func file_api_tk_proto_init() {
	if File_api_tk_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_tk_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_api_tk_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_tk_proto_goTypes,
		DependencyIndexes: file_api_tk_proto_depIdxs,
		EnumInfos:         file_api_tk_proto_enumTypes,
		MessageInfos:      file_api_tk_proto_msgTypes,
	}.Build()
	File_api_tk_proto = out.File
	file_api_tk_proto_rawDesc = nil
	file_api_tk_proto_goTypes = nil
	file_api_tk_proto_depIdxs = nil
}
