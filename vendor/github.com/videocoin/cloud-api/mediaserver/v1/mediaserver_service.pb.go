// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mediaserver/v1/mediaserver_service.proto

package v1

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type StreamRequest struct {
	StreamId             string   `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty" validate:"required"`
	Sdp                  string   `protobuf:"bytes,2,opt,name=sdp,proto3" json:"sdp,omitempty" validate:"required"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamRequest) Reset()         { *m = StreamRequest{} }
func (m *StreamRequest) String() string { return proto.CompactTextString(m) }
func (*StreamRequest) ProtoMessage()    {}
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_87d4b62d2d585b2e, []int{0}
}
func (m *StreamRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StreamRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StreamRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StreamRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamRequest.Merge(m, src)
}
func (m *StreamRequest) XXX_Size() int {
	return m.Size()
}
func (m *StreamRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StreamRequest proto.InternalMessageInfo

func (m *StreamRequest) GetStreamId() string {
	if m != nil {
		return m.StreamId
	}
	return ""
}

func (m *StreamRequest) GetSdp() string {
	if m != nil {
		return m.Sdp
	}
	return ""
}

func (*StreamRequest) XXX_MessageName() string {
	return "cloud.api.mediaserver.v1.StreamRequest"
}

type WebRTCStreamResponse struct {
	Sdp                  string   `protobuf:"bytes,1,opt,name=sdp,proto3" json:"sdp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WebRTCStreamResponse) Reset()         { *m = WebRTCStreamResponse{} }
func (m *WebRTCStreamResponse) String() string { return proto.CompactTextString(m) }
func (*WebRTCStreamResponse) ProtoMessage()    {}
func (*WebRTCStreamResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_87d4b62d2d585b2e, []int{1}
}
func (m *WebRTCStreamResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WebRTCStreamResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WebRTCStreamResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WebRTCStreamResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebRTCStreamResponse.Merge(m, src)
}
func (m *WebRTCStreamResponse) XXX_Size() int {
	return m.Size()
}
func (m *WebRTCStreamResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WebRTCStreamResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WebRTCStreamResponse proto.InternalMessageInfo

func (m *WebRTCStreamResponse) GetSdp() string {
	if m != nil {
		return m.Sdp
	}
	return ""
}

func (*WebRTCStreamResponse) XXX_MessageName() string {
	return "cloud.api.mediaserver.v1.WebRTCStreamResponse"
}
func init() {
	proto.RegisterType((*StreamRequest)(nil), "cloud.api.mediaserver.v1.StreamRequest")
	golang_proto.RegisterType((*StreamRequest)(nil), "cloud.api.mediaserver.v1.StreamRequest")
	proto.RegisterType((*WebRTCStreamResponse)(nil), "cloud.api.mediaserver.v1.WebRTCStreamResponse")
	golang_proto.RegisterType((*WebRTCStreamResponse)(nil), "cloud.api.mediaserver.v1.WebRTCStreamResponse")
}

func init() {
	proto.RegisterFile("mediaserver/v1/mediaserver_service.proto", fileDescriptor_87d4b62d2d585b2e)
}
func init() {
	golang_proto.RegisterFile("mediaserver/v1/mediaserver_service.proto", fileDescriptor_87d4b62d2d585b2e)
}

var fileDescriptor_87d4b62d2d585b2e = []byte{
	// 344 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xc1, 0x4a, 0x33, 0x31,
	0x10, 0x80, 0x49, 0x7f, 0xf8, 0xb1, 0x0b, 0x82, 0x44, 0xc1, 0xb5, 0xca, 0x2a, 0x8b, 0x60, 0x15,
	0x4c, 0xa8, 0x7a, 0xea, 0xb1, 0x3d, 0x79, 0xf0, 0xd2, 0x0a, 0x82, 0x97, 0x92, 0xdd, 0x8c, 0x6b,
	0xa0, 0xbb, 0x49, 0x93, 0xec, 0x8a, 0x57, 0x5f, 0xa1, 0x6f, 0xe1, 0x53, 0x78, 0xec, 0x51, 0xf0,
	0x2e, 0xd2, 0xfa, 0x04, 0x3e, 0x81, 0x6c, 0xb6, 0x62, 0x05, 0xeb, 0x29, 0x99, 0xcc, 0x97, 0x99,
	0xcc, 0x17, 0xaf, 0x99, 0x02, 0x17, 0xcc, 0x80, 0x2e, 0x40, 0xd3, 0xa2, 0x45, 0x17, 0xc2, 0x41,
	0xb9, 0x88, 0x18, 0x88, 0xd2, 0xd2, 0x4a, 0xec, 0xc7, 0x43, 0x99, 0x73, 0xc2, 0x94, 0x20, 0x0b,
	0x10, 0x29, 0x5a, 0x8d, 0xed, 0x44, 0xca, 0x64, 0x08, 0xd4, 0x71, 0x51, 0x7e, 0x43, 0x21, 0x55,
	0xf6, 0xbe, 0xba, 0xd6, 0xd8, 0x99, 0x27, 0x99, 0x12, 0x94, 0x65, 0x99, 0xb4, 0xcc, 0x0a, 0x99,
	0x99, 0x79, 0xf6, 0x38, 0x11, 0xf6, 0x36, 0x8f, 0x48, 0x2c, 0x53, 0x9a, 0xc8, 0x44, 0x7e, 0xd7,
	0x28, 0x23, 0x17, 0xb8, 0x5d, 0x85, 0x87, 0xca, 0x5b, 0xed, 0x5b, 0x0d, 0x2c, 0xed, 0xc1, 0x28,
	0x07, 0x63, 0xf1, 0x99, 0x57, 0x37, 0xee, 0x60, 0x20, 0xb8, 0x8f, 0xf6, 0x50, 0xb3, 0xde, 0xd9,
	0xfc, 0x78, 0xdd, 0x5d, 0x2f, 0xd8, 0x50, 0x70, 0x66, 0xa1, 0x1d, 0x6a, 0x18, 0xe5, 0x42, 0x03,
	0x0f, 0x7b, 0x2b, 0x15, 0x79, 0xce, 0xf1, 0xa1, 0xf7, 0xcf, 0x70, 0xe5, 0xd7, 0xfe, 0xe6, 0x4b,
	0x26, 0x6c, 0x7a, 0x1b, 0x57, 0x10, 0xf5, 0x2e, 0xbb, 0x5f, 0x7d, 0x8d, 0x92, 0x99, 0x01, 0xbc,
	0x56, 0x95, 0x70, 0x2d, 0x1d, 0x79, 0xf2, 0x88, 0x3c, 0x7c, 0x51, 0x8a, 0xe9, 0x3b, 0x31, 0xfd,
	0x4a, 0x1e, 0x1e, 0x23, 0x0f, 0x77, 0x35, 0x30, 0x0b, 0x8b, 0x75, 0xf0, 0x01, 0x59, 0xa6, 0x93,
	0xfc, 0x98, 0xb0, 0x41, 0x96, 0x83, 0xbf, 0x3d, 0x2c, 0xdc, 0x7f, 0x78, 0x79, 0x1f, 0xd7, 0x82,
	0x70, 0xcb, 0x19, 0x2f, 0x7f, 0xd4, 0xd0, 0x6a, 0x70, 0x43, 0xef, 0x20, 0xd2, 0x36, 0x6e, 0xa3,
	0xa3, 0x8e, 0x3f, 0x99, 0x06, 0xe8, 0x79, 0x1a, 0xa0, 0xb7, 0x69, 0x80, 0x9e, 0x66, 0x01, 0x9a,
	0xcc, 0x02, 0x74, 0x5d, 0x2b, 0x5a, 0xd1, 0x7f, 0x67, 0xfa, 0xf4, 0x33, 0x00, 0x00, 0xff, 0xff,
	0x0e, 0x26, 0x7f, 0xc9, 0x19, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MediaServerServiceClient is the client API for MediaServerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MediaServerServiceClient interface {
	CreateWebRTCStream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (*WebRTCStreamResponse, error)
}

type mediaServerServiceClient struct {
	cc *grpc.ClientConn
}

func NewMediaServerServiceClient(cc *grpc.ClientConn) MediaServerServiceClient {
	return &mediaServerServiceClient{cc}
}

func (c *mediaServerServiceClient) CreateWebRTCStream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (*WebRTCStreamResponse, error) {
	out := new(WebRTCStreamResponse)
	err := c.cc.Invoke(ctx, "/cloud.api.mediaserver.v1.MediaServerService/CreateWebRTCStream", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MediaServerServiceServer is the server API for MediaServerService service.
type MediaServerServiceServer interface {
	CreateWebRTCStream(context.Context, *StreamRequest) (*WebRTCStreamResponse, error)
}

// UnimplementedMediaServerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMediaServerServiceServer struct {
}

func (*UnimplementedMediaServerServiceServer) CreateWebRTCStream(ctx context.Context, req *StreamRequest) (*WebRTCStreamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWebRTCStream not implemented")
}

func RegisterMediaServerServiceServer(s *grpc.Server, srv MediaServerServiceServer) {
	s.RegisterService(&_MediaServerService_serviceDesc, srv)
}

func _MediaServerService_CreateWebRTCStream_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServerServiceServer).CreateWebRTCStream(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.api.mediaserver.v1.MediaServerService/CreateWebRTCStream",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServerServiceServer).CreateWebRTCStream(ctx, req.(*StreamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MediaServerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cloud.api.mediaserver.v1.MediaServerService",
	HandlerType: (*MediaServerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateWebRTCStream",
			Handler:    _MediaServerService_CreateWebRTCStream_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mediaserver/v1/mediaserver_service.proto",
}

func (m *StreamRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StreamRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StreamRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Sdp) > 0 {
		i -= len(m.Sdp)
		copy(dAtA[i:], m.Sdp)
		i = encodeVarintMediaserverService(dAtA, i, uint64(len(m.Sdp)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.StreamId) > 0 {
		i -= len(m.StreamId)
		copy(dAtA[i:], m.StreamId)
		i = encodeVarintMediaserverService(dAtA, i, uint64(len(m.StreamId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *WebRTCStreamResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WebRTCStreamResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WebRTCStreamResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Sdp) > 0 {
		i -= len(m.Sdp)
		copy(dAtA[i:], m.Sdp)
		i = encodeVarintMediaserverService(dAtA, i, uint64(len(m.Sdp)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMediaserverService(dAtA []byte, offset int, v uint64) int {
	offset -= sovMediaserverService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StreamRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.StreamId)
	if l > 0 {
		n += 1 + l + sovMediaserverService(uint64(l))
	}
	l = len(m.Sdp)
	if l > 0 {
		n += 1 + l + sovMediaserverService(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *WebRTCStreamResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sdp)
	if l > 0 {
		n += 1 + l + sovMediaserverService(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMediaserverService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMediaserverService(x uint64) (n int) {
	return sovMediaserverService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StreamRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMediaserverService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: StreamRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StreamId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMediaserverService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMediaserverService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMediaserverService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StreamId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sdp", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMediaserverService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMediaserverService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMediaserverService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sdp = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMediaserverService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMediaserverService
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMediaserverService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *WebRTCStreamResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMediaserverService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WebRTCStreamResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WebRTCStreamResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sdp", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMediaserverService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMediaserverService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMediaserverService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sdp = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMediaserverService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMediaserverService
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMediaserverService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMediaserverService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMediaserverService
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMediaserverService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMediaserverService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMediaserverService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMediaserverService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMediaserverService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMediaserverService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMediaserverService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMediaserverService = fmt.Errorf("proto: unexpected end of group")
)
