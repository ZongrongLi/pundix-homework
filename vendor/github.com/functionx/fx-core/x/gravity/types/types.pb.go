// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gravity/v1/types.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

// BridgeValidator represents a validator's ETH address and its power
type BridgeValidator struct {
	Power      uint64 `protobuf:"varint,1,opt,name=power,proto3" json:"power,omitempty"`
	EthAddress string `protobuf:"bytes,2,opt,name=eth_address,json=ethAddress,proto3" json:"eth_address,omitempty"`
}

func (m *BridgeValidator) Reset()         { *m = BridgeValidator{} }
func (m *BridgeValidator) String() string { return proto.CompactTextString(m) }
func (*BridgeValidator) ProtoMessage()    {}
func (*BridgeValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_163831c23fcc179f, []int{0}
}
func (m *BridgeValidator) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BridgeValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BridgeValidator.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BridgeValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BridgeValidator.Merge(m, src)
}
func (m *BridgeValidator) XXX_Size() int {
	return m.Size()
}
func (m *BridgeValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_BridgeValidator.DiscardUnknown(m)
}

var xxx_messageInfo_BridgeValidator proto.InternalMessageInfo

func (m *BridgeValidator) GetPower() uint64 {
	if m != nil {
		return m.Power
	}
	return 0
}

func (m *BridgeValidator) GetEthAddress() string {
	if m != nil {
		return m.EthAddress
	}
	return ""
}

// Valset is the Ethereum Bridge Multsig Set, each gravity validator also
// maintains an ETH key to sign messages, these are used to check signatures on
// ETH because of the significant gas savings
type Valset struct {
	Nonce   uint64             `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Members []*BridgeValidator `protobuf:"bytes,2,rep,name=members,proto3" json:"members,omitempty"`
	Height  uint64             `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *Valset) Reset()         { *m = Valset{} }
func (m *Valset) String() string { return proto.CompactTextString(m) }
func (*Valset) ProtoMessage()    {}
func (*Valset) Descriptor() ([]byte, []int) {
	return fileDescriptor_163831c23fcc179f, []int{1}
}
func (m *Valset) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Valset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Valset.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Valset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Valset.Merge(m, src)
}
func (m *Valset) XXX_Size() int {
	return m.Size()
}
func (m *Valset) XXX_DiscardUnknown() {
	xxx_messageInfo_Valset.DiscardUnknown(m)
}

var xxx_messageInfo_Valset proto.InternalMessageInfo

func (m *Valset) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *Valset) GetMembers() []*BridgeValidator {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *Valset) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

// LastObservedEthereumBlockHeight stores the last observed
// Ethereum block height along with the fx block height that
// it was observed at. These two numbers can be used to project
// outward and always produce batches with timeouts in the future
// even if no Ethereum block height has been relayed for a long time
type LastObservedEthereumBlockHeight struct {
	FxBlockHeight  uint64 `protobuf:"varint,1,opt,name=fx_block_height,json=fxBlockHeight,proto3" json:"fx_block_height,omitempty"`
	EthBlockHeight uint64 `protobuf:"varint,2,opt,name=eth_block_height,json=ethBlockHeight,proto3" json:"eth_block_height,omitempty"`
}

func (m *LastObservedEthereumBlockHeight) Reset()         { *m = LastObservedEthereumBlockHeight{} }
func (m *LastObservedEthereumBlockHeight) String() string { return proto.CompactTextString(m) }
func (*LastObservedEthereumBlockHeight) ProtoMessage()    {}
func (*LastObservedEthereumBlockHeight) Descriptor() ([]byte, []int) {
	return fileDescriptor_163831c23fcc179f, []int{2}
}
func (m *LastObservedEthereumBlockHeight) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LastObservedEthereumBlockHeight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LastObservedEthereumBlockHeight.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LastObservedEthereumBlockHeight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LastObservedEthereumBlockHeight.Merge(m, src)
}
func (m *LastObservedEthereumBlockHeight) XXX_Size() int {
	return m.Size()
}
func (m *LastObservedEthereumBlockHeight) XXX_DiscardUnknown() {
	xxx_messageInfo_LastObservedEthereumBlockHeight.DiscardUnknown(m)
}

var xxx_messageInfo_LastObservedEthereumBlockHeight proto.InternalMessageInfo

func (m *LastObservedEthereumBlockHeight) GetFxBlockHeight() uint64 {
	if m != nil {
		return m.FxBlockHeight
	}
	return 0
}

func (m *LastObservedEthereumBlockHeight) GetEthBlockHeight() uint64 {
	if m != nil {
		return m.EthBlockHeight
	}
	return 0
}

// This records the relationship between an ERC20 token and the denom
// of the corresponding fx originated asset
type ERC20ToDenom struct {
	Erc20 string `protobuf:"bytes,1,opt,name=erc20,proto3" json:"erc20,omitempty"`
	Denom string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *ERC20ToDenom) Reset()         { *m = ERC20ToDenom{} }
func (m *ERC20ToDenom) String() string { return proto.CompactTextString(m) }
func (*ERC20ToDenom) ProtoMessage()    {}
func (*ERC20ToDenom) Descriptor() ([]byte, []int) {
	return fileDescriptor_163831c23fcc179f, []int{3}
}
func (m *ERC20ToDenom) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ERC20ToDenom) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ERC20ToDenom.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ERC20ToDenom) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ERC20ToDenom.Merge(m, src)
}
func (m *ERC20ToDenom) XXX_Size() int {
	return m.Size()
}
func (m *ERC20ToDenom) XXX_DiscardUnknown() {
	xxx_messageInfo_ERC20ToDenom.DiscardUnknown(m)
}

var xxx_messageInfo_ERC20ToDenom proto.InternalMessageInfo

func (m *ERC20ToDenom) GetErc20() string {
	if m != nil {
		return m.Erc20
	}
	return ""
}

func (m *ERC20ToDenom) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func init() {
	proto.RegisterType((*BridgeValidator)(nil), "fx.gravity.v1.BridgeValidator")
	proto.RegisterType((*Valset)(nil), "fx.gravity.v1.Valset")
	proto.RegisterType((*LastObservedEthereumBlockHeight)(nil), "fx.gravity.v1.LastObservedEthereumBlockHeight")
	proto.RegisterType((*ERC20ToDenom)(nil), "fx.gravity.v1.ERC20ToDenom")
}

func init() { proto.RegisterFile("gravity/v1/types.proto", fileDescriptor_163831c23fcc179f) }

var fileDescriptor_163831c23fcc179f = []byte{
	// 353 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0x4f, 0x4f, 0xea, 0x40,
	0x14, 0xc5, 0x29, 0xbc, 0xc7, 0x0b, 0xc3, 0xe3, 0xf1, 0xd2, 0x18, 0xc2, 0x6a, 0x20, 0x5d, 0x98,
	0x2e, 0xb4, 0x05, 0xdc, 0x18, 0x77, 0xa2, 0x18, 0x16, 0x26, 0x26, 0x8d, 0x61, 0xe1, 0x86, 0xf4,
	0xcf, 0x6d, 0xa7, 0x91, 0x76, 0x9a, 0x99, 0xa1, 0x96, 0x6f, 0xe1, 0xc7, 0x72, 0xc9, 0xd2, 0xa5,
	0x81, 0x2f, 0x62, 0x66, 0x5a, 0x12, 0x70, 0x79, 0xee, 0xfd, 0xcd, 0xc9, 0xb9, 0x73, 0x50, 0x2f,
	0x62, 0x6e, 0x1e, 0x8b, 0x8d, 0x9d, 0x8f, 0x6d, 0xb1, 0xc9, 0x80, 0x5b, 0x19, 0xa3, 0x82, 0xea,
	0x9d, 0xb0, 0xb0, 0xaa, 0x95, 0x95, 0x8f, 0x8d, 0x39, 0xea, 0x4e, 0x59, 0x1c, 0x44, 0xb0, 0x70,
	0x57, 0x71, 0xe0, 0x0a, 0xca, 0xf4, 0x33, 0xf4, 0x3b, 0xa3, 0x6f, 0xc0, 0xfa, 0xda, 0x50, 0x33,
	0x7f, 0x39, 0xa5, 0xd0, 0x07, 0xa8, 0x0d, 0x82, 0x2c, 0xdd, 0x20, 0x60, 0xc0, 0x79, 0xbf, 0x3e,
	0xd4, 0xcc, 0x96, 0x83, 0x40, 0x90, 0xdb, 0x72, 0x62, 0x64, 0xa8, 0xb9, 0x70, 0x57, 0x1c, 0x84,
	0x34, 0x48, 0x69, 0xea, 0xc3, 0xc1, 0x40, 0x09, 0xfd, 0x1a, 0xfd, 0x49, 0x20, 0xf1, 0x80, 0xc9,
	0xc7, 0x0d, 0xb3, 0x3d, 0xc1, 0xd6, 0x49, 0x14, 0xeb, 0x47, 0x0e, 0xe7, 0x80, 0xeb, 0x3d, 0xd4,
	0x24, 0x10, 0x47, 0x44, 0xf4, 0x1b, 0xca, 0xb0, 0x52, 0x06, 0x47, 0x83, 0x47, 0x97, 0x8b, 0x27,
	0x8f, 0x03, 0xcb, 0x21, 0x98, 0x09, 0x02, 0x0c, 0xd6, 0xc9, 0x74, 0x45, 0xfd, 0xd7, 0xb9, 0x42,
	0xf4, 0x73, 0xd4, 0x0d, 0x8b, 0xa5, 0x27, 0x27, 0xcb, 0xca, 0xa3, 0x0c, 0xd5, 0x09, 0x8b, 0x63,
	0xce, 0x44, 0xff, 0xe5, 0x75, 0x27, 0x60, 0x5d, 0x81, 0xff, 0x40, 0x90, 0x23, 0xd2, 0xb8, 0x41,
	0x7f, 0x67, 0xce, 0xdd, 0x64, 0xf4, 0x4c, 0xef, 0x21, 0xa5, 0x89, 0x3c, 0x16, 0x98, 0x3f, 0x19,
	0x29, 0xdf, 0x96, 0x53, 0x0a, 0x39, 0x0d, 0xe4, 0xba, 0xfa, 0xa7, 0x52, 0x4c, 0x1f, 0x3e, 0x76,
	0x58, 0xdb, 0xee, 0xb0, 0xf6, 0xb5, 0xc3, 0xda, 0xfb, 0x1e, 0xd7, 0xb6, 0x7b, 0x5c, 0xfb, 0xdc,
	0xe3, 0xda, 0xcb, 0x45, 0x14, 0x0b, 0xb2, 0xf6, 0x2c, 0x9f, 0x26, 0x76, 0xb8, 0x4e, 0x7d, 0x11,
	0xd3, 0xb4, 0xb0, 0xc3, 0xe2, 0xd2, 0xa7, 0x0c, 0xec, 0xc2, 0x3e, 0x94, 0xa9, 0x9a, 0xf4, 0x9a,
	0xaa, 0xca, 0xab, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x43, 0x38, 0x2e, 0x0f, 0xe4, 0x01, 0x00,
	0x00,
}

func (m *BridgeValidator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BridgeValidator) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BridgeValidator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EthAddress) > 0 {
		i -= len(m.EthAddress)
		copy(dAtA[i:], m.EthAddress)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.EthAddress)))
		i--
		dAtA[i] = 0x12
	}
	if m.Power != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Valset) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Valset) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Valset) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Height != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Members) > 0 {
		for iNdEx := len(m.Members) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Members[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Nonce != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LastObservedEthereumBlockHeight) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LastObservedEthereumBlockHeight) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LastObservedEthereumBlockHeight) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.EthBlockHeight != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.EthBlockHeight))
		i--
		dAtA[i] = 0x10
	}
	if m.FxBlockHeight != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.FxBlockHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ERC20ToDenom) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ERC20ToDenom) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ERC20ToDenom) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Erc20) > 0 {
		i -= len(m.Erc20)
		copy(dAtA[i:], m.Erc20)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Erc20)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BridgeValidator) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Power != 0 {
		n += 1 + sovTypes(uint64(m.Power))
	}
	l = len(m.EthAddress)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *Valset) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Nonce != 0 {
		n += 1 + sovTypes(uint64(m.Nonce))
	}
	if len(m.Members) > 0 {
		for _, e := range m.Members {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if m.Height != 0 {
		n += 1 + sovTypes(uint64(m.Height))
	}
	return n
}

func (m *LastObservedEthereumBlockHeight) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.FxBlockHeight != 0 {
		n += 1 + sovTypes(uint64(m.FxBlockHeight))
	}
	if m.EthBlockHeight != 0 {
		n += 1 + sovTypes(uint64(m.EthBlockHeight))
	}
	return n
}

func (m *ERC20ToDenom) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Erc20)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BridgeValidator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: BridgeValidator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BridgeValidator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EthAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *Valset) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: Valset: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Valset: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Members", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Members = append(m.Members, &BridgeValidator{})
			if err := m.Members[len(m.Members)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *LastObservedEthereumBlockHeight) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: LastObservedEthereumBlockHeight: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LastObservedEthereumBlockHeight: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FxBlockHeight", wireType)
			}
			m.FxBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FxBlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthBlockHeight", wireType)
			}
			m.EthBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EthBlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *ERC20ToDenom) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: ERC20ToDenom: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ERC20ToDenom: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Erc20", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Erc20 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
