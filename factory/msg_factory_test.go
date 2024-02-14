package factory

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

// proto message for test
type PrebuildReq struct {
	Type   int32 `protobuf:"varint,1,opt,name=Type,proto3" json:"Type,omitempty"`
	PosX   int32 `protobuf:"varint,2,opt,name=PosX,proto3" json:"PosX,omitempty"`
	PosZ   int32 `protobuf:"varint,3,opt,name=PosZ,proto3" json:"PosZ,omitempty"`
	CityID int32 `protobuf:"varint,4,opt,name=CityID,proto3" json:"CityID,omitempty"`
}

func (m *PrebuildReq) Reset()         { *m = PrebuildReq{} }
func (m *PrebuildReq) String() string { return proto.CompactTextString(m) }
func (*PrebuildReq) ProtoMessage()    {}

func TestHasValidSuffix(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Req", true},
		{"FooReq", true},
		{"FooReq-", false},
		{"Ack", true},
		{"FooAck", true},
		{"FooAck-", false},
		{"Ntf", true},
		{"Ntf-", false},
	}
	for i, tt := range tests {
		if got := HasValidSuffix(tt.input); got != tt.want {
			t.Errorf("case %d IsRequestMessage() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestIsRequestMessage(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Req", true},
		{"FooReq", true},
		{"FooReq-", false},
	}
	for i, tt := range tests {
		if got := IsRequestMessage(tt.input); got != tt.want {
			t.Errorf("case %d IsRequestMessage() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestIsAckMessage(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Ack", true},
		{"FooAck", true},
		{"FooAck-", false},
	}
	for i, tt := range tests {
		if got := IsAckMessage(tt.input); got != tt.want {
			t.Errorf("case %d IsAckMessage() = %v, want %v", i+1, got, tt.want)
		}
	}
}
