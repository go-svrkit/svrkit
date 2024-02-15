package qnet

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/factory"
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
	Code int32 `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Id   int32 `protobuf:"varint,2,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (m *PrebuildAck) Reset()         { *m = PrebuildAck{} }
func (m *PrebuildAck) String() string { return proto.CompactTextString(m) }
func (*PrebuildAck) ProtoMessage()    {}

func TestTryEnqueueMsg(t *testing.T) {
	var queue = make(chan *NetMessage, 1)
	var msg = AllocNetMessage()
	assert.True(t, TryEnqueueMsg(queue, msg))
	assert.False(t, TryEnqueueMsg(queue, msg))
	assert.NotNil(t, TryDequeueMsg(queue))
	assert.Nil(t, TryDequeueMsg(queue))
	assert.True(t, TryEnqueueMsg(queue, msg))
}

func TestAllocNetMessage(t *testing.T) {
	var msg = AllocNetMessage()
	assert.NotNil(t, msg)
	FreeNetMessage(msg)
}

func TestNewNetMessage(t *testing.T) {
	var msg = NewNetMessage(1234, 5678, []byte("hello"))
	assert.Equal(t, msg.Command, uint32(1234))
	assert.Equal(t, msg.Seq, uint32(5678))
	assert.Equal(t, "hello", string(msg.Data))
}

func TestCreateNetMessage(t *testing.T) {
	var req = &PrebuildReq{PosX: 11, PosZ: 22}
	var msg = CreateNetMessage(1234, 5678, req)
	assert.Equal(t, msg.Command, uint32(1234))
	assert.Equal(t, msg.Seq, uint32(5678))
	dt, ok := msg.Body.(*PrebuildReq)
	assert.True(t, ok)
	assert.Equal(t, dt, req)
}

func TestCreateNetMessageWith(t *testing.T) {
	var req = &PrebuildReq{PosX: 11, PosZ: 22}
	var msg = CreateNetMessageWith(req)
	assert.Equal(t, msg.Command, uint32(2523815879))
	dt, ok := msg.Body.(*PrebuildReq)
	assert.True(t, ok)
	assert.Equal(t, dt, req)
}

func TestNetMessage_Reset(t *testing.T) {
	var req = &PrebuildReq{PosX: 11, PosZ: 22}
	var msg = CreateNetMessage(1234, 5678, req)
	assert.NotNil(t, msg)
	assert.Equal(t, msg.Command, uint32(1234))
	assert.Equal(t, msg.Seq, uint32(5678))
	assert.NotNil(t, msg.Body)
	msg.Reset()
	assert.Equal(t, msg.Command, uint32(0))
	assert.Equal(t, msg.Seq, uint32(0))
	assert.Nil(t, msg.Body)
}

func TestNetMessage_Clone(t *testing.T) {
	var req = &PrebuildReq{PosX: 11, PosZ: 22}
	var msg = CreateNetMessage(1234, 5678, req)
	assert.Equal(t, msg.Command, uint32(1234))
	assert.Equal(t, msg.Seq, uint32(5678))
	var clone = msg.Clone()
	assert.NotEqual(t, msg.CreatedAt, clone.CreatedAt)
	assert.Equal(t, clone.Command, uint32(1234))
	assert.Equal(t, clone.Seq, uint32(5678))
}

func TestNetMessage_Encode(t *testing.T) {
	var req = &PrebuildReq{PosX: 11, PosZ: 22}
	var msg = CreateNetMessage(1234, 5678, req)
	assert.Nil(t, msg.Data)
	assert.Nil(t, msg.Encode())
	assert.NotNil(t, msg.Data)

	var breq PrebuildReq
	assert.Nil(t, msg.DecodeTo(&breq))
	assert.Nil(t, msg.Data)
	assert.Equal(t, req.PosX, breq.PosX)
	assert.Equal(t, req.PosZ, breq.PosZ)
}

func TestNetMessage_Refuse(t *testing.T) {
	defer factory.Clear()
	assert.Nil(t, factory.Register(reflect.TypeOf((*PrebuildReq)(nil))))
	assert.Nil(t, factory.Register(reflect.TypeOf((*PrebuildAck)(nil))))

	var queue = make(chan *NetMessage, 1)
	var session = NewFakeSession(queue)
	session.Running.Store(true)

	var req = &PrebuildReq{PosX: 11, PosZ: 22}
	var msg = CreateNetMessageWith(req)
	msg.Session = session
	assert.Nil(t, msg.Refuse(10001))
	var ackMsg = <-queue
	assert.NotNil(t, ackMsg)
	assert.Equal(t, "qnet.PrebuildAck", factory.GetMessageFullName(ackMsg.Command))
	ack, ok := ackMsg.Body.(*PrebuildAck)
	assert.True(t, ok)
	assert.NotNil(t, ack)
	assert.Equal(t, int32(10001), ack.Code)
}
