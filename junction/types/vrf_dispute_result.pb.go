// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: junction/junction/vrf_dispute_result.proto

package types

import (
	encoding_binary "encoding/binary"
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type VrfDisputeResult struct {
	Votes               []bool   `protobuf:"varint,1,rep,packed,name=votes,proto3" json:"votes,omitempty"`
	AddressList         []string `protobuf:"bytes,2,rep,name=addressList,proto3" json:"addressList,omitempty"`
	ConsentVote         uint64   `protobuf:"varint,3,opt,name=consentVote,proto3" json:"consentVote,omitempty"`
	DissentVote         uint64   `protobuf:"varint,4,opt,name=dissentVote,proto3" json:"dissentVote,omitempty"`
	AgreementPercentage float32  `protobuf:"fixed32,5,opt,name=agreementPercentage,proto3" json:"agreementPercentage,omitempty"`
	Result              bool     `protobuf:"varint,6,opt,name=result,proto3" json:"result,omitempty"`
	Message             string   `protobuf:"bytes,7,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *VrfDisputeResult) Reset()         { *m = VrfDisputeResult{} }
func (m *VrfDisputeResult) String() string { return proto.CompactTextString(m) }
func (*VrfDisputeResult) ProtoMessage()    {}
func (*VrfDisputeResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_294207027caba5f4, []int{0}
}
func (m *VrfDisputeResult) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *VrfDisputeResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_VrfDisputeResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *VrfDisputeResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VrfDisputeResult.Merge(m, src)
}
func (m *VrfDisputeResult) XXX_Size() int {
	return m.Size()
}
func (m *VrfDisputeResult) XXX_DiscardUnknown() {
	xxx_messageInfo_VrfDisputeResult.DiscardUnknown(m)
}

var xxx_messageInfo_VrfDisputeResult proto.InternalMessageInfo

func (m *VrfDisputeResult) GetVotes() []bool {
	if m != nil {
		return m.Votes
	}
	return nil
}

func (m *VrfDisputeResult) GetAddressList() []string {
	if m != nil {
		return m.AddressList
	}
	return nil
}

func (m *VrfDisputeResult) GetConsentVote() uint64 {
	if m != nil {
		return m.ConsentVote
	}
	return 0
}

func (m *VrfDisputeResult) GetDissentVote() uint64 {
	if m != nil {
		return m.DissentVote
	}
	return 0
}

func (m *VrfDisputeResult) GetAgreementPercentage() float32 {
	if m != nil {
		return m.AgreementPercentage
	}
	return 0
}

func (m *VrfDisputeResult) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

func (m *VrfDisputeResult) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*VrfDisputeResult)(nil), "junction.junction.VrfDisputeResult")
}

func init() {
	proto.RegisterFile("junction/junction/vrf_dispute_result.proto", fileDescriptor_294207027caba5f4)
}

var fileDescriptor_294207027caba5f4 = []byte{
	// 281 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x85, 0xeb, 0xfe, 0xd7, 0x2c, 0x60, 0x10, 0xf2, 0x64, 0x59, 0x4c, 0x16, 0x43, 0x8b, 0xe0,
	0x0d, 0x80, 0x8d, 0x0e, 0xc8, 0x43, 0x07, 0x96, 0x2a, 0x4d, 0x6e, 0x43, 0x10, 0xb1, 0x23, 0xfb,
	0xa6, 0x82, 0xb7, 0xe0, 0xb1, 0x18, 0x3b, 0x32, 0xa2, 0xe4, 0x0d, 0x78, 0x02, 0x94, 0x86, 0xb4,
	0x19, 0xd8, 0xce, 0x39, 0xfe, 0x3c, 0xdc, 0x8f, 0x5e, 0xbe, 0xe4, 0x26, 0xc4, 0xc4, 0x9a, 0xd9,
	0x3e, 0x6c, 0xdc, 0x7a, 0x19, 0x25, 0x3e, 0xcb, 0x11, 0x96, 0x0e, 0x7c, 0xfe, 0x8a, 0xd3, 0xcc,
	0x59, 0xb4, 0xec, 0xa4, 0x41, 0xa6, 0x4d, 0xb8, 0xf8, 0x21, 0xf4, 0x78, 0xe1, 0xd6, 0xf7, 0x35,
	0xae, 0x77, 0x34, 0x3b, 0xa3, 0x83, 0x8d, 0x45, 0xf0, 0x9c, 0xc8, 0x9e, 0x1a, 0xeb, 0xba, 0x30,
	0x49, 0x8f, 0x82, 0x28, 0x72, 0xe0, 0xfd, 0x3c, 0xf1, 0xc8, 0xbb, 0xb2, 0xa7, 0x26, 0xba, 0x3d,
	0x55, 0x44, 0x68, 0x8d, 0x07, 0x83, 0x0b, 0x8b, 0xc0, 0x7b, 0x92, 0xa8, 0xbe, 0x6e, 0x4f, 0x15,
	0x11, 0x25, 0x7e, 0x4f, 0xf4, 0x6b, 0xa2, 0x35, 0xb1, 0x2b, 0x7a, 0x1a, 0xc4, 0x0e, 0x20, 0x05,
	0x83, 0x8f, 0xe0, 0x42, 0x30, 0x18, 0xc4, 0xc0, 0x07, 0x92, 0xa8, 0xae, 0xfe, 0xef, 0x89, 0x9d,
	0xd3, 0x61, 0x7d, 0x25, 0x1f, 0x4a, 0xa2, 0xc6, 0xfa, 0xaf, 0x31, 0x4e, 0x47, 0x29, 0x78, 0x5f,
	0xfd, 0x1e, 0x49, 0xa2, 0x26, 0xba, 0xa9, 0xb7, 0xf3, 0xcf, 0x42, 0x90, 0x6d, 0x21, 0xc8, 0x77,
	0x21, 0xc8, 0x47, 0x29, 0x3a, 0xdb, 0x52, 0x74, 0xbe, 0x4a, 0xd1, 0x79, 0xba, 0x8e, 0x13, 0x7c,
	0xce, 0x57, 0xd3, 0xd0, 0xa6, 0xb3, 0x3b, 0x9b, 0x56, 0x4e, 0xdc, 0x03, 0x40, 0x14, 0x1c, 0xec,
	0xbe, 0x1d, 0x22, 0xbe, 0x67, 0xe0, 0x57, 0xc3, 0x9d, 0xdc, 0x9b, 0xdf, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xea, 0x46, 0xa7, 0x76, 0x8a, 0x01, 0x00, 0x00,
}

func (m *VrfDisputeResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VrfDisputeResult) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *VrfDisputeResult) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintVrfDisputeResult(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x3a
	}
	if m.Result {
		i--
		if m.Result {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if m.AgreementPercentage != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.AgreementPercentage))))
		i--
		dAtA[i] = 0x2d
	}
	if m.DissentVote != 0 {
		i = encodeVarintVrfDisputeResult(dAtA, i, uint64(m.DissentVote))
		i--
		dAtA[i] = 0x20
	}
	if m.ConsentVote != 0 {
		i = encodeVarintVrfDisputeResult(dAtA, i, uint64(m.ConsentVote))
		i--
		dAtA[i] = 0x18
	}
	if len(m.AddressList) > 0 {
		for iNdEx := len(m.AddressList) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AddressList[iNdEx])
			copy(dAtA[i:], m.AddressList[iNdEx])
			i = encodeVarintVrfDisputeResult(dAtA, i, uint64(len(m.AddressList[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Votes) > 0 {
		for iNdEx := len(m.Votes) - 1; iNdEx >= 0; iNdEx-- {
			i--
			if m.Votes[iNdEx] {
				dAtA[i] = 1
			} else {
				dAtA[i] = 0
			}
		}
		i = encodeVarintVrfDisputeResult(dAtA, i, uint64(len(m.Votes)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintVrfDisputeResult(dAtA []byte, offset int, v uint64) int {
	offset -= sovVrfDisputeResult(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *VrfDisputeResult) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Votes) > 0 {
		n += 1 + sovVrfDisputeResult(uint64(len(m.Votes))) + len(m.Votes)*1
	}
	if len(m.AddressList) > 0 {
		for _, s := range m.AddressList {
			l = len(s)
			n += 1 + l + sovVrfDisputeResult(uint64(l))
		}
	}
	if m.ConsentVote != 0 {
		n += 1 + sovVrfDisputeResult(uint64(m.ConsentVote))
	}
	if m.DissentVote != 0 {
		n += 1 + sovVrfDisputeResult(uint64(m.DissentVote))
	}
	if m.AgreementPercentage != 0 {
		n += 5
	}
	if m.Result {
		n += 2
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovVrfDisputeResult(uint64(l))
	}
	return n
}

func sovVrfDisputeResult(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVrfDisputeResult(x uint64) (n int) {
	return sovVrfDisputeResult(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *VrfDisputeResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVrfDisputeResult
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
			return fmt.Errorf("proto: VrfDisputeResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VrfDisputeResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowVrfDisputeResult
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Votes = append(m.Votes, bool(v != 0))
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowVrfDisputeResult
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthVrfDisputeResult
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthVrfDisputeResult
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen
				if elementCount != 0 && len(m.Votes) == 0 {
					m.Votes = make([]bool, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowVrfDisputeResult
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Votes = append(m.Votes, bool(v != 0))
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Votes", wireType)
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressList", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrfDisputeResult
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
				return ErrInvalidLengthVrfDisputeResult
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVrfDisputeResult
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddressList = append(m.AddressList, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsentVote", wireType)
			}
			m.ConsentVote = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrfDisputeResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ConsentVote |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DissentVote", wireType)
			}
			m.DissentVote = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrfDisputeResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DissentVote |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field AgreementPercentage", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.AgreementPercentage = float32(math.Float32frombits(v))
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrfDisputeResult
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Result = bool(v != 0)
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVrfDisputeResult
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
				return ErrInvalidLengthVrfDisputeResult
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVrfDisputeResult
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVrfDisputeResult(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVrfDisputeResult
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipVrfDisputeResult(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVrfDisputeResult
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
					return 0, ErrIntOverflowVrfDisputeResult
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
					return 0, ErrIntOverflowVrfDisputeResult
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
				return 0, ErrInvalidLengthVrfDisputeResult
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVrfDisputeResult
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVrfDisputeResult
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVrfDisputeResult        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVrfDisputeResult          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVrfDisputeResult = fmt.Errorf("proto: unexpected end of group")
)
