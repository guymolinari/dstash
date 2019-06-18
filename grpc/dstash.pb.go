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
	// 574 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x51, 0x4f, 0xdb, 0x30,
	0x18, 0xac, 0x5b, 0x92, 0x91, 0xaf, 0x50, 0x55, 0xde, 0x86, 0xba, 0x82, 0xa6, 0x2a, 0x2f, 0x8b,
	0x86, 0x14, 0x50, 0x99, 0x90, 0xc6, 0xa4, 0x3d, 0x04, 0x52, 0x2d, 0xda, 0x56, 0xba, 0xa4, 0xf0,
	0x6e, 0x5a, 0x93, 0x46, 0x6d, 0xe3, 0xcc, 0x71, 0xc6, 0xf2, 0x33, 0xf7, 0xb4, 0xbf, 0x33, 0xc5,
	0x4e, 0x69, 0xa1, 0xc0, 0xb6, 0xb7, 0xdc, 0x7d, 0xdf, 0xd9, 0x97, 0xd3, 0x19, 0xb6, 0xc6, 0xa9,
	0x20, 0xe9, 0xc4, 0x4e, 0x38, 0x13, 0x0c, 0xeb, 0x0a, 0xb5, 0x77, 0x43, 0xc6, 0xc2, 0x19, 0x3d,
	0x90, 0xec, 0x55, 0x76, 0x7d, 0x40, 0xe7, 0x89, 0xc8, 0xd5, 0x52, 0xfb, 0xf5, 0xfd, 0xe1, 0x0d,
	0x27, 0x49, 0x42, 0x79, 0xaa, 0xe6, 0xe6, 0x2b, 0x78, 0x16, 0x64, 0xa3, 0x11, 0x4d, 0x53, 0xdc,
	0x80, 0x2a, 0x9b, 0xb6, 0x50, 0x07, 0x59, 0x9b, 0x7e, 0x95, 0x4d, 0xcd, 0x37, 0xb0, 0x1d, 0x08,
	0x22, 0xb2, 0xf4, 0x2b, 0x4d, 0x53, 0x12, 0x52, 0xbc, 0x03, 0x7a, 0x2a, 0x09, 0xb9, 0x64, 0xf8,
	0x25, 0x32, 0x03, 0xa8, 0x7b, 0xf1, 0x98, 0xfe, 0xfc, 0x7c, 0x39, 0x20, 0x11, 0xc7, 0x7b, 0x60,
	0x44, 0x05, 0x1c, 0x10, 0x31, 0x29, 0x37, 0x97, 0x04, 0x6e, 0x42, 0x6d, 0x4a, 0xf3, 0x56, 0xb5,
	0x83, 0xac, 0x2d, 0xbf, 0xf8, 0xc4, 0x2f, 0x40, 0xfb, 0x41, 0x66, 0x19, 0x6d, 0xd5, 0x24, 0xa7,
	0x80, 0x79, 0x02, 0x75, 0x27, 0x12, 0x73, 0x92, 0x7c, 0xcb, 0x28, 0xcf, 0xf1, 0x3e, 0x68, 0xdf,
	0x8b, 0x8f, 0x16, 0xea, 0xd4, 0xac, 0x7a, 0xf7, 0xa5, 0x5d, 0x46, 0x21, 0xa7, 0x3d, 0x4e, 0xc2,
	0x39, 0x8d, 0x85, 0xaf, 0x76, 0xcc, 0xdf, 0x08, 0xb6, 0xef, 0x0c, 0x8a, 0x3b, 0xa4, 0x85, 0xd2,
	0x8f, 0x02, 0x05, 0x7b, 0x1d, 0xd1, 0xd9, 0x58, 0xba, 0x31, 0x7c, 0x05, 0x0a, 0x96, 0xb3, 0x1b,
	0xef, 0x4c, 0xfa, 0xd9, 0xf0, 0x15, 0x58, 0xba, 0xdc, 0xe8, 0x20, 0xab, 0x56, 0xba, 0xc4, 0x27,
	0x60, 0xb0, 0x84, 0x72, 0x22, 0x22, 0x16, 0xb7, 0xb4, 0x0e, 0xb2, 0x1a, 0xdd, 0xbd, 0x07, 0xad,
	0xd9, 0xe7, 0xc9, 0x30, 0x4f, 0xa8, 0xbf, 0x5c, 0x37, 0xbb, 0xa0, 0x2b, 0x12, 0x6f, 0x83, 0xe1,
	0xf5, 0x87, 0xae, 0x1f, 0xb8, 0xa7, 0xc3, 0x66, 0x05, 0x1b, 0xa0, 0x5d, 0xf4, 0xbd, 0xf3, 0x7e,
	0x13, 0xe1, 0x06, 0xc0, 0x99, 0xd7, 0xeb, 0xb9, 0xbe, 0xdb, 0x3f, 0x75, 0x9b, 0x55, 0xf3, 0x10,
	0xf4, 0x32, 0xe5, 0x32, 0x47, 0xf4, 0x40, 0x8e, 0xd5, 0x95, 0x1c, 0xbb, 0xbf, 0x36, 0x40, 0x3f,
	0x0b, 0x0a, 0x43, 0xf8, 0x00, 0x6a, 0x83, 0x4c, 0xe0, 0xc6, 0xc2, 0xa0, 0x3a, 0xa9, 0xbd, 0x63,
	0xab, 0x8e, 0xd8, 0x8b, 0x8e, 0xd8, 0x6e, 0x51, 0x20, 0xb3, 0x82, 0x8f, 0x61, 0xd3, 0x21, 0x62,
	0x34, 0xf9, 0x2f, 0x95, 0x85, 0xf0, 0x5b, 0xd0, 0xbf, 0x30, 0x36, 0xcd, 0x92, 0x35, 0xd5, 0x3d,
	0x6c, 0x56, 0xf0, 0x11, 0xd4, 0xe5, 0x1d, 0xff, 0x2a, 0xb0, 0xd0, 0x21, 0xc2, 0x47, 0xa0, 0x79,
	0x82, 0xce, 0x53, 0xfc, 0x88, 0x8b, 0x75, 0xd9, 0x21, 0xc2, 0xef, 0x41, 0x57, 0x7d, 0x7e, 0x54,
	0x75, 0xdb, 0xaa, 0x3b, 0xbd, 0x37, 0x2b, 0xb8, 0x07, 0x20, 0x4d, 0xca, 0x9a, 0xe3, 0xbd, 0x35,
	0x79, 0x20, 0x78, 0x14, 0x87, 0x97, 0x45, 0xe0, 0x4f, 0x06, 0xf3, 0x09, 0xf4, 0x80, 0x12, 0x3e,
	0x9a, 0xfc, 0xe5, 0x8c, 0xf5, 0xe9, 0x85, 0x17, 0x8b, 0xe3, 0x77, 0x72, 0x2a, 0x7f, 0xe6, 0x63,
	0x19, 0x5b, 0x40, 0x85, 0x13, 0x09, 0xfc, 0x7c, 0xe1, 0x7c, 0xe5, 0x21, 0x3e, 0xe9, 0xe4, 0x03,
	0x68, 0xea, 0x61, 0xdd, 0x2a, 0x57, 0x5e, 0x5b, 0x7b, 0x77, 0x4d, 0xe9, 0xe4, 0x82, 0xa6, 0xe5,
	0xf5, 0xce, 0x3e, 0xb4, 0x23, 0x66, 0x87, 0x3c, 0x19, 0xd9, 0x61, 0x96, 0xcf, 0xd9, 0x2c, 0x8a,
	0x09, 0x8f, 0xca, 0x83, 0x9c, 0xba, 0xaa, 0xdb, 0xa0, 0x90, 0x0e, 0xd0, 0x95, 0x2e, 0xcf, 0x38,
	0xfa, 0x13, 0x00, 0x00, 0xff, 0xff, 0x21, 0x1c, 0xb3, 0xf6, 0xbd, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DStashClient is the client API for DStash service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DStashClient interface {
	Put(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*empty.Empty, error)
	BatchPut(ctx context.Context, opts ...grpc.CallOption) (DStash_BatchPutClient, error)
	Lookup(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*KVPair, error)
	BatchLookup(ctx context.Context, opts ...grpc.CallOption) (DStash_BatchLookupClient, error)
	Items(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (DStash_ItemsClient, error)
	Status(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*StatusMessage, error)
	BatchIndex(ctx context.Context, opts ...grpc.CallOption) (DStash_BatchIndexClient, error)
	Search(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (DStash_SearchClient, error)
	BatchSetBit(ctx context.Context, opts ...grpc.CallOption) (DStash_BatchSetBitClient, error)
	Query(ctx context.Context, in *BitmapQuery, opts ...grpc.CallOption) (*wrappers.BytesValue, error)
}

type dStashClient struct {
	cc *grpc.ClientConn
}

func NewDStashClient(cc *grpc.ClientConn) DStashClient {
	return &dStashClient{cc}
}

func (c *dStashClient) Put(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/dstash.DStash/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dStashClient) BatchPut(ctx context.Context, opts ...grpc.CallOption) (DStash_BatchPutClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DStash_serviceDesc.Streams[0], "/dstash.DStash/BatchPut", opts...)
	if err != nil {
		return nil, err
	}
	x := &dStashBatchPutClient{stream}
	return x, nil
}

type DStash_BatchPutClient interface {
	Send(*KVPair) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type dStashBatchPutClient struct {
	grpc.ClientStream
}

func (x *dStashBatchPutClient) Send(m *KVPair) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dStashBatchPutClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dStashClient) Lookup(ctx context.Context, in *KVPair, opts ...grpc.CallOption) (*KVPair, error) {
	out := new(KVPair)
	err := c.cc.Invoke(ctx, "/dstash.DStash/Lookup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dStashClient) BatchLookup(ctx context.Context, opts ...grpc.CallOption) (DStash_BatchLookupClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DStash_serviceDesc.Streams[1], "/dstash.DStash/BatchLookup", opts...)
	if err != nil {
		return nil, err
	}
	x := &dStashBatchLookupClient{stream}
	return x, nil
}

type DStash_BatchLookupClient interface {
	Send(*KVPair) error
	Recv() (*KVPair, error)
	grpc.ClientStream
}

type dStashBatchLookupClient struct {
	grpc.ClientStream
}

func (x *dStashBatchLookupClient) Send(m *KVPair) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dStashBatchLookupClient) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dStashClient) Items(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (DStash_ItemsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DStash_serviceDesc.Streams[2], "/dstash.DStash/Items", opts...)
	if err != nil {
		return nil, err
	}
	x := &dStashItemsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DStash_ItemsClient interface {
	Recv() (*KVPair, error)
	grpc.ClientStream
}

type dStashItemsClient struct {
	grpc.ClientStream
}

func (x *dStashItemsClient) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dStashClient) Status(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*StatusMessage, error) {
	out := new(StatusMessage)
	err := c.cc.Invoke(ctx, "/dstash.DStash/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dStashClient) BatchIndex(ctx context.Context, opts ...grpc.CallOption) (DStash_BatchIndexClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DStash_serviceDesc.Streams[3], "/dstash.DStash/BatchIndex", opts...)
	if err != nil {
		return nil, err
	}
	x := &dStashBatchIndexClient{stream}
	return x, nil
}

type DStash_BatchIndexClient interface {
	Send(*wrappers.StringValue) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type dStashBatchIndexClient struct {
	grpc.ClientStream
}

func (x *dStashBatchIndexClient) Send(m *wrappers.StringValue) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dStashBatchIndexClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dStashClient) Search(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (DStash_SearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DStash_serviceDesc.Streams[4], "/dstash.DStash/Search", opts...)
	if err != nil {
		return nil, err
	}
	x := &dStashSearchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DStash_SearchClient interface {
	Recv() (*wrappers.UInt64Value, error)
	grpc.ClientStream
}

type dStashSearchClient struct {
	grpc.ClientStream
}

func (x *dStashSearchClient) Recv() (*wrappers.UInt64Value, error) {
	m := new(wrappers.UInt64Value)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dStashClient) BatchSetBit(ctx context.Context, opts ...grpc.CallOption) (DStash_BatchSetBitClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DStash_serviceDesc.Streams[5], "/dstash.DStash/BatchSetBit", opts...)
	if err != nil {
		return nil, err
	}
	x := &dStashBatchSetBitClient{stream}
	return x, nil
}

type DStash_BatchSetBitClient interface {
	Send(*IndexKVPair) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type dStashBatchSetBitClient struct {
	grpc.ClientStream
}

func (x *dStashBatchSetBitClient) Send(m *IndexKVPair) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dStashBatchSetBitClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dStashClient) Query(ctx context.Context, in *BitmapQuery, opts ...grpc.CallOption) (*wrappers.BytesValue, error) {
	out := new(wrappers.BytesValue)
	err := c.cc.Invoke(ctx, "/dstash.DStash/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DStashServer is the server API for DStash service.
type DStashServer interface {
	Put(context.Context, *KVPair) (*empty.Empty, error)
	BatchPut(DStash_BatchPutServer) error
	Lookup(context.Context, *KVPair) (*KVPair, error)
	BatchLookup(DStash_BatchLookupServer) error
	Items(*empty.Empty, DStash_ItemsServer) error
	Status(context.Context, *empty.Empty) (*StatusMessage, error)
	BatchIndex(DStash_BatchIndexServer) error
	Search(*wrappers.StringValue, DStash_SearchServer) error
	BatchSetBit(DStash_BatchSetBitServer) error
	Query(context.Context, *BitmapQuery) (*wrappers.BytesValue, error)
}

// UnimplementedDStashServer can be embedded to have forward compatible implementations.
type UnimplementedDStashServer struct {
}

func (*UnimplementedDStashServer) Put(ctx context.Context, req *KVPair) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (*UnimplementedDStashServer) BatchPut(srv DStash_BatchPutServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchPut not implemented")
}
func (*UnimplementedDStashServer) Lookup(ctx context.Context, req *KVPair) (*KVPair, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lookup not implemented")
}
func (*UnimplementedDStashServer) BatchLookup(srv DStash_BatchLookupServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchLookup not implemented")
}
func (*UnimplementedDStashServer) Items(req *empty.Empty, srv DStash_ItemsServer) error {
	return status.Errorf(codes.Unimplemented, "method Items not implemented")
}
func (*UnimplementedDStashServer) Status(ctx context.Context, req *empty.Empty) (*StatusMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (*UnimplementedDStashServer) BatchIndex(srv DStash_BatchIndexServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchIndex not implemented")
}
func (*UnimplementedDStashServer) Search(req *wrappers.StringValue, srv DStash_SearchServer) error {
	return status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (*UnimplementedDStashServer) BatchSetBit(srv DStash_BatchSetBitServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchSetBit not implemented")
}
func (*UnimplementedDStashServer) Query(ctx context.Context, req *BitmapQuery) (*wrappers.BytesValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}

func RegisterDStashServer(s *grpc.Server, srv DStashServer) {
	s.RegisterService(&_DStash_serviceDesc, srv)
}

func _DStash_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KVPair)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DStashServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dstash.DStash/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DStashServer).Put(ctx, req.(*KVPair))
	}
	return interceptor(ctx, in, info, handler)
}

func _DStash_BatchPut_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DStashServer).BatchPut(&dStashBatchPutServer{stream})
}

type DStash_BatchPutServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*KVPair, error)
	grpc.ServerStream
}

type dStashBatchPutServer struct {
	grpc.ServerStream
}

func (x *dStashBatchPutServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dStashBatchPutServer) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DStash_Lookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KVPair)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DStashServer).Lookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dstash.DStash/Lookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DStashServer).Lookup(ctx, req.(*KVPair))
	}
	return interceptor(ctx, in, info, handler)
}

func _DStash_BatchLookup_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DStashServer).BatchLookup(&dStashBatchLookupServer{stream})
}

type DStash_BatchLookupServer interface {
	Send(*KVPair) error
	Recv() (*KVPair, error)
	grpc.ServerStream
}

type dStashBatchLookupServer struct {
	grpc.ServerStream
}

func (x *dStashBatchLookupServer) Send(m *KVPair) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dStashBatchLookupServer) Recv() (*KVPair, error) {
	m := new(KVPair)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DStash_Items_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DStashServer).Items(m, &dStashItemsServer{stream})
}

type DStash_ItemsServer interface {
	Send(*KVPair) error
	grpc.ServerStream
}

type dStashItemsServer struct {
	grpc.ServerStream
}

func (x *dStashItemsServer) Send(m *KVPair) error {
	return x.ServerStream.SendMsg(m)
}

func _DStash_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DStashServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dstash.DStash/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DStashServer).Status(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DStash_BatchIndex_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DStashServer).BatchIndex(&dStashBatchIndexServer{stream})
}

type DStash_BatchIndexServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*wrappers.StringValue, error)
	grpc.ServerStream
}

type dStashBatchIndexServer struct {
	grpc.ServerStream
}

func (x *dStashBatchIndexServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dStashBatchIndexServer) Recv() (*wrappers.StringValue, error) {
	m := new(wrappers.StringValue)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DStash_Search_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(wrappers.StringValue)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DStashServer).Search(m, &dStashSearchServer{stream})
}

type DStash_SearchServer interface {
	Send(*wrappers.UInt64Value) error
	grpc.ServerStream
}

type dStashSearchServer struct {
	grpc.ServerStream
}

func (x *dStashSearchServer) Send(m *wrappers.UInt64Value) error {
	return x.ServerStream.SendMsg(m)
}

func _DStash_BatchSetBit_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DStashServer).BatchSetBit(&dStashBatchSetBitServer{stream})
}

type DStash_BatchSetBitServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*IndexKVPair, error)
	grpc.ServerStream
}

type dStashBatchSetBitServer struct {
	grpc.ServerStream
}

func (x *dStashBatchSetBitServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dStashBatchSetBitServer) Recv() (*IndexKVPair, error) {
	m := new(IndexKVPair)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DStash_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BitmapQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DStashServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dstash.DStash/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DStashServer).Query(ctx, req.(*BitmapQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var _DStash_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dstash.DStash",
	HandlerType: (*DStashServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _DStash_Put_Handler,
		},
		{
			MethodName: "Lookup",
			Handler:    _DStash_Lookup_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _DStash_Status_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _DStash_Query_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BatchPut",
			Handler:       _DStash_BatchPut_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "BatchLookup",
			Handler:       _DStash_BatchLookup_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Items",
			Handler:       _DStash_Items_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BatchIndex",
			Handler:       _DStash_BatchIndex_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Search",
			Handler:       _DStash_Search_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BatchSetBit",
			Handler:       _DStash_BatchSetBit_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "dstash.proto",
}
