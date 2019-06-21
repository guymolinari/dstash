// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dstash.proto

package dstash

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type QueryFragment_OpType int32

const (
	QueryFragment_INTERSECT  QueryFragment_OpType = 0
	QueryFragment_UNION      QueryFragment_OpType = 1
	QueryFragment_DIFFERENCE QueryFragment_OpType = 2
)

var QueryFragment_OpType_name = map[int32]string{
	0: "INTERSECT",
	1: "UNION",
	2: "DIFFERENCE",
}

var QueryFragment_OpType_value = map[string]int32{
	"INTERSECT":  0,
	"UNION":      1,
	"DIFFERENCE": 2,
}

func (x QueryFragment_OpType) String() string {
	return proto.EnumName(QueryFragment_OpType_name, int32(x))
}

func (QueryFragment_OpType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0c6a1428cb258f96, []int{4, 0}
}

type Success struct {
	Ok                   bool     `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Success) Reset()         { *m = Success{} }
func (m *Success) String() string { return proto.CompactTextString(m) }
func (*Success) ProtoMessage()    {}
func (*Success) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c6a1428cb258f96, []int{0}
}

func (m *Success) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Success.Unmarshal(m, b)
}
func (m *Success) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Success.Marshal(b, m, deterministic)
}
func (m *Success) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Success.Merge(m, src)
}
func (m *Success) XXX_Size() int {
	return xxx_messageInfo_Success.Size(m)
}
func (m *Success) XXX_DiscardUnknown() {
	xxx_messageInfo_Success.DiscardUnknown(m)
}

var xxx_messageInfo_Success proto.InternalMessageInfo

func (m *Success) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

type StatusMessage struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusMessage) Reset()         { *m = StatusMessage{} }
func (m *StatusMessage) String() string { return proto.CompactTextString(m) }
func (*StatusMessage) ProtoMessage()    {}
func (*StatusMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c6a1428cb258f96, []int{1}
}

func (m *StatusMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusMessage.Unmarshal(m, b)
}
func (m *StatusMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusMessage.Marshal(b, m, deterministic)
}
func (m *StatusMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusMessage.Merge(m, src)
}
func (m *StatusMessage) XXX_Size() int {
	return xxx_messageInfo_StatusMessage.Size(m)
}
func (m *StatusMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusMessage.DiscardUnknown(m)
}

var xxx_messageInfo_StatusMessage proto.InternalMessageInfo

func (m *StatusMessage) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type IndexKVPair struct {
	IndexPath            string   `protobuf:"bytes,1,opt,name=indexPath,proto3" json:"indexPath,omitempty"`
	Key                  []byte   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IndexKVPair) Reset()         { *m = IndexKVPair{} }
func (m *IndexKVPair) String() string { return proto.CompactTextString(m) }
func (*IndexKVPair) ProtoMessage()    {}
func (*IndexKVPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c6a1428cb258f96, []int{2}
}

func (m *IndexKVPair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IndexKVPair.Unmarshal(m, b)
}
func (m *IndexKVPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IndexKVPair.Marshal(b, m, deterministic)
}
func (m *IndexKVPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexKVPair.Merge(m, src)
}
func (m *IndexKVPair) XXX_Size() int {
	return xxx_messageInfo_IndexKVPair.Size(m)
}
func (m *IndexKVPair) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexKVPair.DiscardUnknown(m)
}

var xxx_messageInfo_IndexKVPair proto.InternalMessageInfo

func (m *IndexKVPair) GetIndexPath() string {
	if m != nil {
		return m.IndexPath
	}
	return ""
}

func (m *IndexKVPair) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *IndexKVPair) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type BitmapQuery struct {
	Query                []*QueryFragment `protobuf:"bytes,1,rep,name=query,proto3" json:"query,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *BitmapQuery) Reset()         { *m = BitmapQuery{} }
func (m *BitmapQuery) String() string { return proto.CompactTextString(m) }
func (*BitmapQuery) ProtoMessage()    {}
func (*BitmapQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c6a1428cb258f96, []int{3}
}

func (m *BitmapQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BitmapQuery.Unmarshal(m, b)
}
func (m *BitmapQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BitmapQuery.Marshal(b, m, deterministic)
}
func (m *BitmapQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BitmapQuery.Merge(m, src)
}
func (m *BitmapQuery) XXX_Size() int {
	return xxx_messageInfo_BitmapQuery.Size(m)
}
func (m *BitmapQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_BitmapQuery.DiscardUnknown(m)
}

var xxx_messageInfo_BitmapQuery proto.InternalMessageInfo

func (m *BitmapQuery) GetQuery() []*QueryFragment {
	if m != nil {
		return m.Query
	}
	return nil
}

type QueryFragment struct {
	Index                string               `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	Field                string               `protobuf:"bytes,2,opt,name=field,proto3" json:"field,omitempty"`
	RowID                uint64               `protobuf:"varint,3,opt,name=rowID,proto3" json:"rowID,omitempty"`
	Value                int64                `protobuf:"varint,4,opt,name=value,proto3" json:"value,omitempty"`
	Operation            QueryFragment_OpType `protobuf:"varint,5,opt,name=operation,proto3,enum=dstash.QueryFragment_OpType" json:"operation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *QueryFragment) Reset()         { *m = QueryFragment{} }
func (m *QueryFragment) String() string { return proto.CompactTextString(m) }
func (*QueryFragment) ProtoMessage()    {}
func (*QueryFragment) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c6a1428cb258f96, []int{4}
}

func (m *QueryFragment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryFragment.Unmarshal(m, b)
}
func (m *QueryFragment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryFragment.Marshal(b, m, deterministic)
}
func (m *QueryFragment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryFragment.Merge(m, src)
}
func (m *QueryFragment) XXX_Size() int {
	return xxx_messageInfo_QueryFragment.Size(m)
}
func (m *QueryFragment) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryFragment.DiscardUnknown(m)
}

var xxx_messageInfo_QueryFragment proto.InternalMessageInfo

func (m *QueryFragment) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *QueryFragment) GetField() string {
	if m != nil {
		return m.Field
	}
	return ""
}

func (m *QueryFragment) GetRowID() uint64 {
	if m != nil {
		return m.RowID
	}
	return 0
}

func (m *QueryFragment) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *QueryFragment) GetOperation() QueryFragment_OpType {
	if m != nil {
		return m.Operation
	}
	return QueryFragment_INTERSECT
}

type KVPair struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KVPair) Reset()         { *m = KVPair{} }
func (m *KVPair) String() string { return proto.CompactTextString(m) }
func (*KVPair) ProtoMessage()    {}
func (*KVPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c6a1428cb258f96, []int{5}
}

func (m *KVPair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVPair.Unmarshal(m, b)
}
func (m *KVPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVPair.Marshal(b, m, deterministic)
}
func (m *KVPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVPair.Merge(m, src)
}
func (m *KVPair) XXX_Size() int {
	return xxx_messageInfo_KVPair.Size(m)
}
func (m *KVPair) XXX_DiscardUnknown() {
	xxx_messageInfo_KVPair.DiscardUnknown(m)
}

var xxx_messageInfo_KVPair proto.InternalMessageInfo

func (m *KVPair) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *KVPair) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterEnum("dstash.QueryFragment_OpType", QueryFragment_OpType_name, QueryFragment_OpType_value)
	proto.RegisterType((*Success)(nil), "dstash.Success")
	proto.RegisterType((*StatusMessage)(nil), "dstash.StatusMessage")
	proto.RegisterType((*IndexKVPair)(nil), "dstash.IndexKVPair")
	proto.RegisterType((*BitmapQuery)(nil), "dstash.BitmapQuery")
	proto.RegisterType((*QueryFragment)(nil), "dstash.QueryFragment")
	proto.RegisterType((*KVPair)(nil), "dstash.KVPair")
}

func init() { proto.RegisterFile("dstash.proto", fileDescriptor_0c6a1428cb258f96) }

var fileDescriptor_0c6a1428cb258f96 = []byte{
	// 610 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0xed, 0x4e, 0xdb, 0x30,
	0x14, 0xad, 0x5b, 0x1a, 0xe8, 0x6d, 0xa9, 0x2a, 0x6f, 0x43, 0x5d, 0x41, 0x53, 0x95, 0x3f, 0x8b,
	0x86, 0x14, 0x50, 0x98, 0x90, 0xc6, 0xa4, 0x49, 0x0b, 0xa4, 0x5a, 0xc4, 0x56, 0xba, 0x04, 0xf8,
	0x6f, 0x5a, 0x93, 0x46, 0x6d, 0xe2, 0xcc, 0x71, 0xc6, 0xf2, 0x06, 0xd3, 0x9e, 0x62, 0x6f, 0xb5,
	0xd7, 0x99, 0x62, 0xa7, 0x7c, 0x15, 0xd0, 0xf6, 0x2b, 0x39, 0xf7, 0xde, 0x73, 0x7d, 0x7c, 0x7c,
	0x2f, 0xb4, 0x26, 0xa9, 0x20, 0xe9, 0xd4, 0x4c, 0x38, 0x13, 0x0c, 0x6b, 0x0a, 0xf5, 0x36, 0x03,
	0xc6, 0x82, 0x39, 0xdd, 0x91, 0xd1, 0x8b, 0xec, 0x72, 0x87, 0x46, 0x89, 0xc8, 0x55, 0x51, 0xef,
	0xd5, 0xfd, 0xe4, 0x15, 0x27, 0x49, 0x42, 0x79, 0xaa, 0xf2, 0xfa, 0x4b, 0x58, 0xf5, 0xb3, 0xf1,
	0x98, 0xa6, 0x29, 0x6e, 0x43, 0x95, 0xcd, 0xba, 0xa8, 0x8f, 0x8c, 0x35, 0xaf, 0xca, 0x66, 0xfa,
	0x6b, 0x58, 0xf7, 0x05, 0x11, 0x59, 0xfa, 0x85, 0xa6, 0x29, 0x09, 0x28, 0xde, 0x00, 0x2d, 0x95,
	0x01, 0x59, 0xd4, 0xf0, 0x4a, 0xa4, 0xfb, 0xd0, 0x74, 0xe3, 0x09, 0xfd, 0x71, 0x7c, 0x3e, 0x22,
	0x21, 0xc7, 0x5b, 0xd0, 0x08, 0x0b, 0x38, 0x22, 0x62, 0x5a, 0x56, 0xde, 0x04, 0x70, 0x07, 0x6a,
	0x33, 0x9a, 0x77, 0xab, 0x7d, 0x64, 0xb4, 0xbc, 0xe2, 0x17, 0x3f, 0x87, 0xfa, 0x77, 0x32, 0xcf,
	0x68, 0xb7, 0x26, 0x63, 0x0a, 0xe8, 0x07, 0xd0, 0xb4, 0x43, 0x11, 0x91, 0xe4, 0x6b, 0x46, 0x79,
	0x8e, 0xb7, 0xa1, 0xfe, 0xad, 0xf8, 0xe9, 0xa2, 0x7e, 0xcd, 0x68, 0x5a, 0x2f, 0xcc, 0xd2, 0x0a,
	0x99, 0x1d, 0x70, 0x12, 0x44, 0x34, 0x16, 0x9e, 0xaa, 0xd1, 0xff, 0x20, 0x58, 0xbf, 0x93, 0x28,
	0xce, 0x90, 0x12, 0x4a, 0x3d, 0x0a, 0x14, 0xd1, 0xcb, 0x90, 0xce, 0x27, 0x52, 0x4d, 0xc3, 0x53,
	0xa0, 0x88, 0x72, 0x76, 0xe5, 0x1e, 0x49, 0x3d, 0x2b, 0x9e, 0x02, 0x37, 0x2a, 0x57, 0xfa, 0xc8,
	0xa8, 0x95, 0x2a, 0xf1, 0x01, 0x34, 0x58, 0x42, 0x39, 0x11, 0x21, 0x8b, 0xbb, 0xf5, 0x3e, 0x32,
	0xda, 0xd6, 0xd6, 0x83, 0xd2, 0xcc, 0x93, 0xe4, 0x34, 0x4f, 0xa8, 0x77, 0x53, 0xae, 0x5b, 0xa0,
	0xa9, 0x20, 0x5e, 0x87, 0x86, 0x3b, 0x3c, 0x75, 0x3c, 0xdf, 0x39, 0x3c, 0xed, 0x54, 0x70, 0x03,
	0xea, 0x67, 0x43, 0xf7, 0x64, 0xd8, 0x41, 0xb8, 0x0d, 0x70, 0xe4, 0x0e, 0x06, 0x8e, 0xe7, 0x0c,
	0x0f, 0x9d, 0x4e, 0x55, 0xdf, 0x05, 0xad, 0x74, 0xb9, 0xf4, 0x11, 0x3d, 0xe0, 0x63, 0xf5, 0x96,
	0x8f, 0x96, 0x0b, 0xad, 0xc3, 0x79, 0x96, 0x0a, 0xca, 0x3f, 0x4e, 0xa2, 0x30, 0xc6, 0xef, 0x40,
	0x53, 0xaf, 0x8a, 0x37, 0x4c, 0x35, 0x1b, 0xe6, 0x62, 0x36, 0x4c, 0xa7, 0x18, 0x9c, 0xde, 0xb5,
	0xb7, 0x77, 0x5e, 0x5f, 0xaf, 0x58, 0x3f, 0xab, 0xb0, 0x7a, 0x7c, 0xee, 0x0b, 0xc6, 0x29, 0xde,
	0x81, 0xda, 0x28, 0x13, 0xb8, 0xbd, 0xa8, 0x55, 0xaa, 0x7a, 0x8f, 0xf4, 0xd4, 0x2b, 0x78, 0x1f,
	0xd6, 0x6c, 0x22, 0xc6, 0xd3, 0xff, 0x62, 0x19, 0x08, 0xbf, 0x01, 0xed, 0x33, 0x63, 0xb3, 0x2c,
	0x59, 0x62, 0xdd, 0xc3, 0x7a, 0x05, 0xef, 0x41, 0x53, 0x9e, 0xf1, 0xaf, 0x04, 0x03, 0xed, 0x22,
	0xbc, 0x07, 0x75, 0x57, 0xd0, 0xe8, 0x71, 0x3f, 0x96, 0x68, 0xbb, 0xc8, 0xfa, 0x8d, 0xa0, 0xe5,
	0x0b, 0x1e, 0xc6, 0x81, 0x4f, 0x09, 0x1f, 0x4f, 0xf1, 0x00, 0x40, 0x1e, 0x2d, 0x17, 0x01, 0x6f,
	0x2d, 0xb5, 0x52, 0xc5, 0xe7, 0xc5, 0x93, 0x3c, 0x79, 0xdd, 0x4f, 0xa0, 0x95, 0x1d, 0x9f, 0xee,
	0xb1, 0x9c, 0x3d, 0x73, 0x63, 0xb1, 0xff, 0x56, 0x66, 0xa5, 0xc4, 0x5f, 0x68, 0xb1, 0x41, 0x4a,
	0xd3, 0x87, 0xd2, 0x1c, 0x9f, 0x0a, 0x3b, 0x14, 0xf8, 0xd9, 0xe2, 0x56, 0xb7, 0x56, 0xf7, 0x49,
	0x65, 0xef, 0xa1, 0xae, 0x56, 0xf1, 0x9a, 0x79, 0x6b, 0x3f, 0x7b, 0x9b, 0x4b, 0x4c, 0x3b, 0x17,
	0x34, 0x2d, 0xe5, 0xd8, 0xdb, 0xd0, 0x0b, 0x99, 0x19, 0xf0, 0x64, 0x6c, 0x06, 0x59, 0x1e, 0xb1,
	0x79, 0x18, 0x13, 0x1e, 0x96, 0x8d, 0xec, 0xe6, 0x91, 0x5f, 0x7c, 0x47, 0x05, 0x75, 0x84, 0x2e,
	0x34, 0xd9, 0x63, 0xef, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x32, 0x4a, 0xa2, 0x22, 0xef, 0x04,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClusterAdminClient is the client API for ClusterAdmin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterAdminClient interface {
	Status(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*StatusMessage, error)
}

type clusterAdminClient struct {
	cc *grpc.ClientConn
}

func NewClusterAdminClient(cc *grpc.ClientConn) ClusterAdminClient {
	return &clusterAdminClient{cc}
}

func (c *clusterAdminClient) Status(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*StatusMessage, error) {
	out := new(StatusMessage)
	err := c.cc.Invoke(ctx, "/dstash.ClusterAdmin/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterAdminServer is the server API for ClusterAdmin service.
type ClusterAdminServer interface {
	Status(context.Context, *empty.Empty) (*StatusMessage, error)
}

// UnimplementedClusterAdminServer can be embedded to have forward compatible implementations.
type UnimplementedClusterAdminServer struct {
}

func (*UnimplementedClusterAdminServer) Status(ctx context.Context, req *empty.Empty) (*StatusMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}

func RegisterClusterAdminServer(s *grpc.Server, srv ClusterAdminServer) {
	s.RegisterService(&_ClusterAdmin_serviceDesc, srv)
}

func _ClusterAdmin_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterAdminServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dstash.ClusterAdmin/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterAdminServer).Status(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClusterAdmin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dstash.ClusterAdmin",
	HandlerType: (*ClusterAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _ClusterAdmin_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dstash.proto",
}

// KVStoreClient is the client API for KVStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type KVStoreClient interface {
	Put(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*empty.Empty, error)
	BatchPut(ctx context.Context, opts ...grpc.CallOption) (KVStore_BatchPutClient, error)
	Lookup(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*KVPair, error)
	BatchLookup(ctx context.Context, opts ...grpc.CallOption) (KVStore_BatchLookupClient, error)
	Items(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (KVStore_ItemsClient, error)
}

type kVStoreClient struct {
	cc *grpc.ClientConn
}

func NewKVStoreClient(cc *grpc.ClientConn) KVStoreClient {
	return &kVStoreClient{cc}
}

func (c *kVStoreClient) Put(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/dstash.KVStore/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVStoreClient) BatchPut(ctx context.Context, opts ...grpc.CallOption) (KVStore_BatchPutClient, error) {
	stream, err := c.cc.NewStream(ctx, &_KVStore_serviceDesc.Streams[0], "/dstash.KVStore/BatchPut", opts...)
	if err != nil {
		return nil, err
	}
	x := &kVStoreBatchPutClient{stream}
	return x, nil
}

type KVStore_BatchPutClient interface {
	Send(*KVPair) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type kVStoreBatchPutClient struct {
	grpc.ClientStream
}

func (x *kVStoreBatchPutClient) Send(m *KVPair) error {
	return x.ClientStream.SendMsg(m)
}

func (x *kVStoreBatchPutClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *kVStoreClient) Lookup(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*KVPair, error) {
	out := new(KVPair)
	err := c.cc.Invoke(ctx, "/dstash.KVStore/Lookup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVStoreClient) BatchLookup(ctx context.Context, opts ...grpc.CallOption) (KVStore_BatchLookupClient, error) {
	stream, err := c.cc.NewStream(ctx, &_KVStore_serviceDesc.Streams[1], "/dstash.KVStore/BatchLookup", opts...)
	if err != nil {
		return nil, err
	}
	x := &kVStoreBatchLookupClient{stream}
	return x, nil
}

type KVStore_BatchLookupClient interface {
	Send(*KVPair) error
	Recv() (*KVPair, error)
	grpc.ClientStream
}

type kVStoreBatchLookupClient struct {
	grpc.ClientStream
}

func (x *kVStoreBatchLookupClient) Send(m *KVPair) error {
	return x.ClientStream.SendMsg(m)
}

func (x *kVStoreBatchLookupClient) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *kVStoreClient) Items(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (KVStore_ItemsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_KVStore_serviceDesc.Streams[2], "/dstash.KVStore/Items", opts...)
	if err != nil {
		return nil, err
	}
	x := &kVStoreItemsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KVStore_ItemsClient interface {
	Recv() (*KVPair, error)
	grpc.ClientStream
}

type kVStoreItemsClient struct {
	grpc.ClientStream
}

func (x *kVStoreItemsClient) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// KVStoreServer is the server API for KVStore service.
type KVStoreServer interface {
	Put(context.Context, *KVPair) (*empty.Empty, error)
	BatchPut(KVStore_BatchPutServer) error
	Lookup(context.Context, *KVPair) (*KVPair, error)
	BatchLookup(KVStore_BatchLookupServer) error
	Items(*empty.Empty, KVStore_ItemsServer) error
}

// UnimplementedKVStoreServer can be embedded to have forward compatible implementations.
type UnimplementedKVStoreServer struct {
}

func (*UnimplementedKVStoreServer) Put(ctx context.Context, req *KVPair) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (*UnimplementedKVStoreServer) BatchPut(srv KVStore_BatchPutServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchPut not implemented")
}
func (*UnimplementedKVStoreServer) Lookup(ctx context.Context, req *KVPair) (*KVPair, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lookup not implemented")
}
func (*UnimplementedKVStoreServer) BatchLookup(srv KVStore_BatchLookupServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchLookup not implemented")
}
func (*UnimplementedKVStoreServer) Items(req *empty.Empty, srv KVStore_ItemsServer) error {
	return status.Errorf(codes.Unimplemented, "method Items not implemented")
}

func RegisterKVStoreServer(s *grpc.Server, srv KVStoreServer) {
	s.RegisterService(&_KVStore_serviceDesc, srv)
}

func _KVStore_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KVPair)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVStoreServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dstash.KVStore/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVStoreServer).Put(ctx, req.(*KVPair))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVStore_BatchPut_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(KVStoreServer).BatchPut(&kVStoreBatchPutServer{stream})
}

type KVStore_BatchPutServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*KVPair, error)
	grpc.ServerStream
}

type kVStoreBatchPutServer struct {
	grpc.ServerStream
}

func (x *kVStoreBatchPutServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *kVStoreBatchPutServer) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _KVStore_Lookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KVPair)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVStoreServer).Lookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dstash.KVStore/Lookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVStoreServer).Lookup(ctx, req.(*KVPair))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVStore_BatchLookup_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(KVStoreServer).BatchLookup(&kVStoreBatchLookupServer{stream})
}

type KVStore_BatchLookupServer interface {
	Send(*KVPair) error
	Recv() (*KVPair, error)
	grpc.ServerStream
}

type kVStoreBatchLookupServer struct {
	grpc.ServerStream
}

func (x *kVStoreBatchLookupServer) Send(m *KVPair) error {
	return x.ServerStream.SendMsg(m)
}

func (x *kVStoreBatchLookupServer) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _KVStore_Items_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KVStoreServer).Items(m, &kVStoreItemsServer{stream})
}

type KVStore_ItemsServer interface {
	Send(*KVPair) error
	grpc.ServerStream
}

type kVStoreItemsServer struct {
	grpc.ServerStream
}

func (x *kVStoreItemsServer) Send(m *KVPair) error {
	return x.ServerStream.SendMsg(m)
}

var _KVStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dstash.KVStore",
	HandlerType: (*KVStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _KVStore_Put_Handler,
		},
		{
			MethodName: "Lookup",
			Handler:    _KVStore_Lookup_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BatchPut",
			Handler:       _KVStore_BatchPut_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "BatchLookup",
			Handler:       _KVStore_BatchLookup_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Items",
			Handler:       _KVStore_Items_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "dstash.proto",
}

// StringSearchClient is the client API for StringSearch service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StringSearchClient interface {
	BatchIndex(ctx context.Context, opts ...grpc.CallOption) (StringSearch_BatchIndexClient, error)
	Search(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (StringSearch_SearchClient, error)
}

type stringSearchClient struct {
	cc *grpc.ClientConn
}

func NewStringSearchClient(cc *grpc.ClientConn) StringSearchClient {
	return &stringSearchClient{cc}
}

func (c *stringSearchClient) BatchIndex(ctx context.Context, opts ...grpc.CallOption) (StringSearch_BatchIndexClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StringSearch_serviceDesc.Streams[0], "/dstash.StringSearch/BatchIndex", opts...)
	if err != nil {
		return nil, err
	}
	x := &stringSearchBatchIndexClient{stream}
	return x, nil
}

type StringSearch_BatchIndexClient interface {
	Send(*wrappers.StringValue) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type stringSearchBatchIndexClient struct {
	grpc.ClientStream
}

func (x *stringSearchBatchIndexClient) Send(m *wrappers.StringValue) error {
	return x.ClientStream.SendMsg(m)
}

func (x *stringSearchBatchIndexClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *stringSearchClient) Search(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (StringSearch_SearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StringSearch_serviceDesc.Streams[1], "/dstash.StringSearch/Search", opts...)
	if err != nil {
		return nil, err
	}
	x := &stringSearchSearchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StringSearch_SearchClient interface {
	Recv() (*wrappers.UInt64Value, error)
	grpc.ClientStream
}

type stringSearchSearchClient struct {
	grpc.ClientStream
}

func (x *stringSearchSearchClient) Recv() (*wrappers.UInt64Value, error) {
	m := new(wrappers.UInt64Value)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StringSearchServer is the server API for StringSearch service.
type StringSearchServer interface {
	BatchIndex(StringSearch_BatchIndexServer) error
	Search(*wrappers.StringValue, StringSearch_SearchServer) error
}

// UnimplementedStringSearchServer can be embedded to have forward compatible implementations.
type UnimplementedStringSearchServer struct {
}

func (*UnimplementedStringSearchServer) BatchIndex(srv StringSearch_BatchIndexServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchIndex not implemented")
}
func (*UnimplementedStringSearchServer) Search(req *wrappers.StringValue, srv StringSearch_SearchServer) error {
	return status.Errorf(codes.Unimplemented, "method Search not implemented")
}

func RegisterStringSearchServer(s *grpc.Server, srv StringSearchServer) {
	s.RegisterService(&_StringSearch_serviceDesc, srv)
}

func _StringSearch_BatchIndex_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StringSearchServer).BatchIndex(&stringSearchBatchIndexServer{stream})
}

type StringSearch_BatchIndexServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*wrappers.StringValue, error)
	grpc.ServerStream
}

type stringSearchBatchIndexServer struct {
	grpc.ServerStream
}

func (x *stringSearchBatchIndexServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *stringSearchBatchIndexServer) Recv() (*wrappers.StringValue, error) {
	m := new(wrappers.StringValue)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _StringSearch_Search_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(wrappers.StringValue)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StringSearchServer).Search(m, &stringSearchSearchServer{stream})
}

type StringSearch_SearchServer interface {
	Send(*wrappers.UInt64Value) error
	grpc.ServerStream
}

type stringSearchSearchServer struct {
	grpc.ServerStream
}

func (x *stringSearchSearchServer) Send(m *wrappers.UInt64Value) error {
	return x.ServerStream.SendMsg(m)
}

var _StringSearch_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dstash.StringSearch",
	HandlerType: (*StringSearchServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BatchIndex",
			Handler:       _StringSearch_BatchIndex_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Search",
			Handler:       _StringSearch_Search_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "dstash.proto",
}

// BitmapIndexClient is the client API for BitmapIndex service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BitmapIndexClient interface {
	BatchSetBit(ctx context.Context, opts ...grpc.CallOption) (BitmapIndex_BatchSetBitClient, error)
	Query(ctx context.Context, in *BitmapQuery, opts ...grpc.CallOption) (*wrappers.BytesValue, error)
}

type bitmapIndexClient struct {
	cc *grpc.ClientConn
}

func NewBitmapIndexClient(cc *grpc.ClientConn) BitmapIndexClient {
	return &bitmapIndexClient{cc}
}

func (c *bitmapIndexClient) BatchSetBit(ctx context.Context, opts ...grpc.CallOption) (BitmapIndex_BatchSetBitClient, error) {
	stream, err := c.cc.NewStream(ctx, &_BitmapIndex_serviceDesc.Streams[0], "/dstash.BitmapIndex/BatchSetBit", opts...)
	if err != nil {
		return nil, err
	}
	x := &bitmapIndexBatchSetBitClient{stream}
	return x, nil
}

type BitmapIndex_BatchSetBitClient interface {
	Send(*IndexKVPair) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type bitmapIndexBatchSetBitClient struct {
	grpc.ClientStream
}

func (x *bitmapIndexBatchSetBitClient) Send(m *IndexKVPair) error {
	return x.ClientStream.SendMsg(m)
}

func (x *bitmapIndexBatchSetBitClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *bitmapIndexClient) Query(ctx context.Context, in *BitmapQuery, opts ...grpc.CallOption) (*wrappers.BytesValue, error) {
	out := new(wrappers.BytesValue)
	err := c.cc.Invoke(ctx, "/dstash.BitmapIndex/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BitmapIndexServer is the server API for BitmapIndex service.
type BitmapIndexServer interface {
	BatchSetBit(BitmapIndex_BatchSetBitServer) error
	Query(context.Context, *BitmapQuery) (*wrappers.BytesValue, error)
}

// UnimplementedBitmapIndexServer can be embedded to have forward compatible implementations.
type UnimplementedBitmapIndexServer struct {
}

func (*UnimplementedBitmapIndexServer) BatchSetBit(srv BitmapIndex_BatchSetBitServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchSetBit not implemented")
}
func (*UnimplementedBitmapIndexServer) Query(ctx context.Context, req *BitmapQuery) (*wrappers.BytesValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}

func RegisterBitmapIndexServer(s *grpc.Server, srv BitmapIndexServer) {
	s.RegisterService(&_BitmapIndex_serviceDesc, srv)
}

func _BitmapIndex_BatchSetBit_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BitmapIndexServer).BatchSetBit(&bitmapIndexBatchSetBitServer{stream})
}

type BitmapIndex_BatchSetBitServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*IndexKVPair, error)
	grpc.ServerStream
}

type bitmapIndexBatchSetBitServer struct {
	grpc.ServerStream
}

func (x *bitmapIndexBatchSetBitServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *bitmapIndexBatchSetBitServer) Recv() (*IndexKVPair, error) {
	m := new(IndexKVPair)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _BitmapIndex_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BitmapQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BitmapIndexServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dstash.BitmapIndex/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BitmapIndexServer).Query(ctx, req.(*BitmapQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var _BitmapIndex_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dstash.BitmapIndex",
	HandlerType: (*BitmapIndexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Query",
			Handler:    _BitmapIndex_Query_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BatchSetBit",
			Handler:       _BitmapIndex_BatchSetBit_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "dstash.proto",
}
