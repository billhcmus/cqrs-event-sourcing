// Code generated by protoc-gen-go. DO NOT EDIT.
// source: aggregate.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	aggregate.proto
	api.proto
	command.proto
	common.proto
	event.proto

It has these top-level messages:
	BaseAggregate
	BaseCommand
	Status
	BaseEvent
	SaveEventsRequest
	SaveEventsResponse
	LoadEventsRequest
	LoadEventsResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type BaseAggregate struct {
	AggregateId   string `protobuf:"bytes,1,opt,name=aggregate_id,json=aggregateId" json:"aggregate_id,omitempty"`
	AggregateType string `protobuf:"bytes,2,opt,name=aggregate_type,json=aggregateType" json:"aggregate_type,omitempty"`
	Version       uint64 `protobuf:"varint,3,opt,name=version" json:"version,omitempty"`
}

func (m *BaseAggregate) Reset()                    { *m = BaseAggregate{} }
func (m *BaseAggregate) String() string            { return proto.CompactTextString(m) }
func (*BaseAggregate) ProtoMessage()               {}
func (*BaseAggregate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *BaseAggregate) GetAggregateId() string {
	if m != nil {
		return m.AggregateId
	}
	return ""
}

func (m *BaseAggregate) GetAggregateType() string {
	if m != nil {
		return m.AggregateType
	}
	return ""
}

func (m *BaseAggregate) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func init() {
	proto.RegisterType((*BaseAggregate)(nil), "pb.BaseAggregate")
}

func init() { proto.RegisterFile("aggregate.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 123 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x4c, 0x4f, 0x2f,
	0x4a, 0x4d, 0x4f, 0x2c, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52,
	0x2a, 0xe6, 0xe2, 0x75, 0x4a, 0x2c, 0x4e, 0x75, 0x84, 0x49, 0x09, 0x29, 0x72, 0xf1, 0xc0, 0xd5,
	0xc5, 0x67, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x71, 0xc3, 0xc5, 0x3c, 0x53, 0x84,
	0x54, 0xb9, 0xf8, 0x10, 0x4a, 0x4a, 0x2a, 0x0b, 0x52, 0x25, 0x98, 0xc0, 0x8a, 0x78, 0xe1, 0xa2,
	0x21, 0x95, 0x05, 0xa9, 0x42, 0x12, 0x5c, 0xec, 0x65, 0xa9, 0x45, 0xc5, 0x99, 0xf9, 0x79, 0x12,
	0xcc, 0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x30, 0x6e, 0x12, 0x1b, 0xd8, 0x7e, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x80, 0x3a, 0xa4, 0x7f, 0x92, 0x00, 0x00, 0x00,
}