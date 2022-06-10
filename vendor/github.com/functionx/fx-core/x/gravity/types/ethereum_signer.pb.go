// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gravity/v1/ethereum_signer.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// SignType defines messages that have been signed by an orchestrator
type SignType int32

const (
	SIGN_TYPE_UNSPECIFIED                          SignType = 0
	SIGN_TYPE_ORCHESTRATOR_SIGNED_MULTI_SIG_UPDATE SignType = 1
	SIGN_TYPE_ORCHESTRATOR_SIGNED_WITHDRAW_BATCH   SignType = 2
)

var SignType_name = map[int32]string{
	0: "SIGN_TYPE_UNSPECIFIED",
	1: "SIGN_TYPE_ORCHESTRATOR_SIGNED_MULTI_SIG_UPDATE",
	2: "SIGN_TYPE_ORCHESTRATOR_SIGNED_WITHDRAW_BATCH",
}

var SignType_value = map[string]int32{
	"SIGN_TYPE_UNSPECIFIED":                          0,
	"SIGN_TYPE_ORCHESTRATOR_SIGNED_MULTI_SIG_UPDATE": 1,
	"SIGN_TYPE_ORCHESTRATOR_SIGNED_WITHDRAW_BATCH":   2,
}

func (SignType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_005a3d0c6f36c26c, []int{0}
}

func init() {
	proto.RegisterEnum("fx.gravity.v1.SignType", SignType_name, SignType_value)
}

func init() { proto.RegisterFile("gravity/v1/ethereum_signer.proto", fileDescriptor_005a3d0c6f36c26c) }

var fileDescriptor_005a3d0c6f36c26c = []byte{
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x48, 0x2f, 0x4a, 0x2c,
	0xcb, 0x2c, 0xa9, 0xd4, 0x2f, 0x33, 0xd4, 0x4f, 0x2d, 0xc9, 0x48, 0x2d, 0x4a, 0x2d, 0xcd, 0x8d,
	0x2f, 0xce, 0x4c, 0xcf, 0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4d, 0xab,
	0xd0, 0x83, 0x2a, 0xd2, 0x2b, 0x33, 0x94, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0xcb, 0xe8, 0x83,
	0x58, 0x10, 0x45, 0x5a, 0x53, 0x19, 0xb9, 0x38, 0x82, 0x33, 0xd3, 0xf3, 0x42, 0x2a, 0x0b, 0x52,
	0x85, 0x24, 0xb9, 0x44, 0x83, 0x3d, 0xdd, 0xfd, 0xe2, 0x43, 0x22, 0x03, 0x5c, 0xe3, 0x43, 0xfd,
	0x82, 0x03, 0x5c, 0x9d, 0x3d, 0xdd, 0x3c, 0x5d, 0x5d, 0x04, 0x18, 0x84, 0x8c, 0xb8, 0xf4, 0x10,
	0x52, 0xfe, 0x41, 0xce, 0x1e, 0xae, 0xc1, 0x21, 0x41, 0x8e, 0x21, 0xfe, 0x41, 0xf1, 0x20, 0x61,
	0x57, 0x97, 0x78, 0xdf, 0x50, 0x9f, 0x10, 0x4f, 0x10, 0x27, 0x3e, 0x34, 0xc0, 0xc5, 0x31, 0xc4,
	0x55, 0x80, 0x51, 0xc8, 0x80, 0x4b, 0x07, 0xbf, 0x9e, 0x70, 0xcf, 0x10, 0x0f, 0x97, 0x20, 0xc7,
	0xf0, 0x78, 0x27, 0xc7, 0x10, 0x67, 0x0f, 0x01, 0x26, 0x29, 0x8e, 0x8e, 0xc5, 0x72, 0x0c, 0x2b,
	0x96, 0xc8, 0x31, 0x38, 0xb9, 0x9d, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47,
	0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94,
	0x4e, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x5a, 0x69, 0x5e, 0x72,
	0x49, 0x66, 0x7e, 0x5e, 0x85, 0x7e, 0x5a, 0x85, 0x6e, 0x72, 0x7e, 0x51, 0xaa, 0x7e, 0x85, 0x3e,
	0x2c, 0x5c, 0x4a, 0x2a, 0x0b, 0x52, 0x8b, 0x93, 0xd8, 0xc0, 0xde, 0x34, 0x06, 0x04, 0x00, 0x00,
	0xff, 0xff, 0x77, 0xf6, 0x38, 0xf5, 0x2f, 0x01, 0x00, 0x00,
}
