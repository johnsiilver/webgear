// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: snippets.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type SaveReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnixNano int64  `protobuf:"varint,1,opt,name=unixNano,proto3" json:"unixNano,omitempty"`
	Content  string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *SaveReq) Reset() {
	*x = SaveReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snippets_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveReq) ProtoMessage() {}

func (x *SaveReq) ProtoReflect() protoreflect.Message {
	mi := &file_snippets_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveReq.ProtoReflect.Descriptor instead.
func (*SaveReq) Descriptor() ([]byte, []int) {
	return file_snippets_proto_rawDescGZIP(), []int{0}
}

func (x *SaveReq) GetUnixNano() int64 {
	if x != nil {
		return x.UnixNano
	}
	return 0
}

func (x *SaveReq) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type SaveResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SaveResp) Reset() {
	*x = SaveResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snippets_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveResp) ProtoMessage() {}

func (x *SaveResp) ProtoReflect() protoreflect.Message {
	mi := &file_snippets_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveResp.ProtoReflect.Descriptor instead.
func (*SaveResp) Descriptor() ([]byte, []int) {
	return file_snippets_proto_rawDescGZIP(), []int{1}
}

type GetReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnixNano int64 `protobuf:"varint,1,opt,name=unixNano,proto3" json:"unixNano,omitempty"`
}

func (x *GetReq) Reset() {
	*x = GetReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snippets_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReq) ProtoMessage() {}

func (x *GetReq) ProtoReflect() protoreflect.Message {
	mi := &file_snippets_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReq.ProtoReflect.Descriptor instead.
func (*GetReq) Descriptor() ([]byte, []int) {
	return file_snippets_proto_rawDescGZIP(), []int{2}
}

func (x *GetReq) GetUnixNano() int64 {
	if x != nil {
		return x.UnixNano
	}
	return 0
}

type GetResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnixNano int64  `protobuf:"varint,1,opt,name=unixNano,proto3" json:"unixNano,omitempty"`
	Content  string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *GetResp) Reset() {
	*x = GetResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snippets_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResp) ProtoMessage() {}

func (x *GetResp) ProtoReflect() protoreflect.Message {
	mi := &file_snippets_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResp.ProtoReflect.Descriptor instead.
func (*GetResp) Descriptor() ([]byte, []int) {
	return file_snippets_proto_rawDescGZIP(), []int{3}
}

func (x *GetResp) GetUnixNano() int64 {
	if x != nil {
		return x.UnixNano
	}
	return 0
}

func (x *GetResp) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_snippets_proto protoreflect.FileDescriptor

var file_snippets_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x07, 0x53, 0x61, 0x76, 0x65,
	0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x6e, 0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x75, 0x6e, 0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x0a, 0x0a, 0x08, 0x53, 0x61, 0x76,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x24, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x12,
	0x1a, 0x0a, 0x08, 0x75, 0x6e, 0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x75, 0x6e, 0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x22, 0x3f, 0x0a, 0x07, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x6e, 0x69, 0x78, 0x4e, 0x61,
	0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x75, 0x6e, 0x69, 0x78, 0x4e, 0x61,
	0x6e, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0xae, 0x01, 0x0a,
	0x08, 0x53, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x12, 0x52, 0x0a, 0x04, 0x53, 0x61, 0x76,
	0x65, 0x12, 0x11, 0x2e, 0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x2e, 0x53, 0x61, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x2e,
	0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d,
	0x22, 0x18, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x73, 0x61, 0x76, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x4e, 0x0a,
	0x03, 0x47, 0x65, 0x74, 0x12, 0x10, 0x2e, 0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e, 0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74,
	0x73, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22, 0x22, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1c, 0x22, 0x17, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x67, 0x65, 0x74, 0x3a, 0x01, 0x2a, 0x42, 0x42, 0x5a,
	0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x68, 0x6e,
	0x73, 0x69, 0x69, 0x6c, 0x76, 0x65, 0x72, 0x2f, 0x77, 0x65, 0x62, 0x67, 0x65, 0x61, 0x72, 0x2f,
	0x77, 0x61, 0x73, 0x6d, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x73, 0x6e,
	0x69, 0x70, 0x70, 0x65, 0x74, 0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_snippets_proto_rawDescOnce sync.Once
	file_snippets_proto_rawDescData = file_snippets_proto_rawDesc
)

func file_snippets_proto_rawDescGZIP() []byte {
	file_snippets_proto_rawDescOnce.Do(func() {
		file_snippets_proto_rawDescData = protoimpl.X.CompressGZIP(file_snippets_proto_rawDescData)
	})
	return file_snippets_proto_rawDescData
}

var file_snippets_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_snippets_proto_goTypes = []interface{}{
	(*SaveReq)(nil),  // 0: snippets.SaveReq
	(*SaveResp)(nil), // 1: snippets.SaveResp
	(*GetReq)(nil),   // 2: snippets.GetReq
	(*GetResp)(nil),  // 3: snippets.GetResp
}
var file_snippets_proto_depIdxs = []int32{
	0, // 0: snippets.Snippets.Save:input_type -> snippets.SaveReq
	2, // 1: snippets.Snippets.Get:input_type -> snippets.GetReq
	1, // 2: snippets.Snippets.Save:output_type -> snippets.SaveResp
	3, // 3: snippets.Snippets.Get:output_type -> snippets.GetResp
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_snippets_proto_init() }
func file_snippets_proto_init() {
	if File_snippets_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_snippets_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveReq); i {
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
		file_snippets_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveResp); i {
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
		file_snippets_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReq); i {
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
		file_snippets_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResp); i {
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
			RawDescriptor: file_snippets_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_snippets_proto_goTypes,
		DependencyIndexes: file_snippets_proto_depIdxs,
		MessageInfos:      file_snippets_proto_msgTypes,
	}.Build()
	File_snippets_proto = out.File
	file_snippets_proto_rawDesc = nil
	file_snippets_proto_goTypes = nil
	file_snippets_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SnippetsClient is the client API for Snippets service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SnippetsClient interface {
	Save(ctx context.Context, in *SaveReq, opts ...grpc.CallOption) (*SaveResp, error)
	Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error)
}

type snippetsClient struct {
	cc grpc.ClientConnInterface
}

func NewSnippetsClient(cc grpc.ClientConnInterface) SnippetsClient {
	return &snippetsClient{cc}
}

func (c *snippetsClient) Save(ctx context.Context, in *SaveReq, opts ...grpc.CallOption) (*SaveResp, error) {
	out := new(SaveResp)
	err := c.cc.Invoke(ctx, "/snippets.Snippets/Save", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *snippetsClient) Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error) {
	out := new(GetResp)
	err := c.cc.Invoke(ctx, "/snippets.Snippets/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SnippetsServer is the server API for Snippets service.
type SnippetsServer interface {
	Save(context.Context, *SaveReq) (*SaveResp, error)
	Get(context.Context, *GetReq) (*GetResp, error)
}

// UnimplementedSnippetsServer can be embedded to have forward compatible implementations.
type UnimplementedSnippetsServer struct {
}

func (*UnimplementedSnippetsServer) Save(context.Context, *SaveReq) (*SaveResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (*UnimplementedSnippetsServer) Get(context.Context, *GetReq) (*GetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func RegisterSnippetsServer(s *grpc.Server, srv SnippetsServer) {
	s.RegisterService(&_Snippets_serviceDesc, srv)
}

func _Snippets_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnippetsServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/snippets.Snippets/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnippetsServer).Save(ctx, req.(*SaveReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Snippets_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnippetsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/snippets.Snippets/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnippetsServer).Get(ctx, req.(*GetReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Snippets_serviceDesc = grpc.ServiceDesc{
	ServiceName: "snippets.Snippets",
	HandlerType: (*SnippetsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Save",
			Handler:    _Snippets_Save_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Snippets_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "snippets.proto",
}
