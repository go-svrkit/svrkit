package factory

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// proto message for test
type PrebuildReq struct {
	Type int32 `protobuf:"varint,1,opt,name=Type,proto3" json:"Type,omitempty"`
	PosX int32 `protobuf:"varint,2,opt,name=PosX,proto3" json:"PosX,omitempty"`
	PosZ int32 `protobuf:"varint,3,opt,name=PosZ,proto3" json:"PosZ,omitempty"`
}

func (m *PrebuildReq) Reset()         { *m = PrebuildReq{} }
func (m *PrebuildReq) String() string { return proto.CompactTextString(m) }
func (*PrebuildReq) ProtoMessage()    {}

type PrebuildAck struct {
	Code uint32 `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Id   int32  `protobuf:"varint,2,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (m *PrebuildAck) Reset()         { *m = PrebuildAck{} }
func (m *PrebuildAck) String() string { return proto.CompactTextString(m) }
func (*PrebuildAck) ProtoMessage()    {}

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

func TestIsReqMessage(t *testing.T) {
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
		if got := IsReqMessage(tt.input); got != tt.want {
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

func TestGetPairingAckName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Req", ""},
		{"FooReq", "FooAck"},
		{"FooReqq", ""},
	}
	for i, tt := range tests {
		if got := GetPairingAckName(tt.input); got != tt.want {
			t.Errorf("case %d GetPairingAckName() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestGetPairingAckNameOf(t *testing.T) {
	defer Clear()
	assert.Nil(t, Register(reflect.TypeOf((*PrebuildReq)(nil))))
	assert.Nil(t, Register(reflect.TypeOf((*PrebuildAck)(nil))))

	tests := []struct {
		input uint32
		want  string
	}{
		{0, ""},
		{GetMessageId("factory.PrebuildReq"), "factory.PrebuildAck"},
	}
	for i, tt := range tests {
		if got := GetPairingAckNameOf(tt.input); got != tt.want {
			t.Errorf("case %d GetPairingAckNameOf() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestRegister(t *testing.T) {
	defer Clear()

	assert.Nil(t, Register(reflect.TypeOf((*PrebuildReq)(nil))))

	var name = "factory.PrebuildReq"
	var hash = GetMessageId(name)
	assert.Greater(t, hash, uint32(0))
	assert.Equal(t, GetMessageFullName(hash), name)
	var req PrebuildReq
	assert.Equal(t, GetMessageType(hash).String(), reflect.TypeOf(req).String())
	assert.Equal(t, hash, GetMessageIdOf(&req))
}

func TestCreateMessage(t *testing.T) {
	defer Clear()

	var name = "factory.PrebuildReq"
	assert.Nil(t, Register(reflect.TypeOf((*PrebuildReq)(nil))))

	var hash = GetMessageId(name)
	var msg = CreateMessageByID(hash)
	assert.NotNil(t, msg)
	req, ok := msg.(*PrebuildReq)
	assert.True(t, ok)
	assert.NotNil(t, req)

	var msg2 = CreateMessageByName(name)
	assert.NotNil(t, msg2)
	req, ok = msg.(*PrebuildReq)
	assert.True(t, ok)
	assert.NotNil(t, req)
}

func TestCreatePairingAck(t *testing.T) {
	defer Clear()
	assert.Nil(t, Register(reflect.TypeOf((*PrebuildReq)(nil))))
	assert.Nil(t, Register(reflect.TypeOf((*PrebuildAck)(nil))))

	var name = "factory.PrebuildReq"
	var msg = CreatePairingAck(name)
	assert.NotNil(t, msg)
	ack, ok := msg.(*PrebuildAck)
	assert.True(t, ok)
	assert.NotNil(t, ack)
}
