// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"fmt"
	"reflect"
	"time"

	"gopkg.in/svrkit.v1/codec"
	"gopkg.in/svrkit.v1/pool"
	"gopkg.in/svrkit.v1/reflext"
)

const ErrCodeField = "Code"

type SessionMessage struct {
	Session Endpoint
	MsgId   uint32
	MsgBody codec.Message
}

type NetMessage struct {
	CreatedAt int64         `json:"created_at,omitempty"` // microseconds
	Command   uint32        `json:"cmd"`
	Seq       uint32        `json:"seq,omitempty"`
	Data      []byte        `json:"data,omitempty"`
	Body      codec.Message `json:"body,omitempty"`
	Session   Endpoint      `json:"-"`
}

func NewNetMessage(cmd, seq uint32, data []byte) *NetMessage {
	var msg = AllocNetMessage()
	msg.Command = cmd
	msg.Seq = seq
	msg.Data = data
	return msg
}

func CreateNetMessage(cmd, seq uint32, body codec.Message) *NetMessage {
	var msg = AllocNetMessage()
	msg.Command = cmd
	msg.Seq = seq
	msg.Body = body
	return msg
}

func CreateNetMessageWith(body codec.Message) *NetMessage {
	var msg = AllocNetMessage()
	msg.Command = DefaultMsgIDReflector(body)
	msg.Body = body
	return msg
}

func (m *NetMessage) Reset() {
	m.CreatedAt = 0
	m.Command = 0
	m.Seq = 0
	m.Data = nil
	m.Body = nil
	m.Session = nil
}

func (m *NetMessage) Clone() *NetMessage {
	return &NetMessage{
		CreatedAt: time.Now().UnixMicro(),
		Seq:       m.Seq,
		Command:   m.Command,
		Body:      m.Body,
	}
}

// Encode encode `Body` to `Data`
func (m *NetMessage) Encode() error {
	if m.Data != nil {
		return nil
	}
	if m.Body != nil {
		data, err := codec.Marshal(m.Body)
		if err != nil {
			return err
		}
		m.Data = data
		return nil
	}
	return nil
}

// DecodeTo decode `Data` to `msg`
func (m *NetMessage) DecodeTo(msg codec.Message) error {
	if err := codec.Unmarshal(m.Data, msg); err != nil {
		return err
	}
	m.Data = nil
	return nil
}

func (m *NetMessage) Reply(cmd uint32, data []byte) error {
	var netMsg = AllocNetMessage()
	netMsg.Seq = m.Seq
	netMsg.Command = cmd
	netMsg.Data = data
	return m.Session.SendMsg(netMsg, SendNonblock)
}

func (m *NetMessage) Ack(ack codec.Message) error {
	var netMsg = CreateNetMessageWith(ack)
	netMsg.Seq = m.Seq
	return m.Session.SendMsg(netMsg, SendNonblock)
}

// Refuse 返回一个带错误码的Ack
func (m *NetMessage) Refuse(ec int32) error {
	var fullName = codec.GetMessageFullName(m.Command)
	var ackName = codec.GetPairingAckName(fullName)
	if ackName == "" {
		return fmt.Errorf("%s(%d) not req message", fullName, m.Command)
	}
	var ack = codec.CreateMessageByName(ackName)
	if ack == nil {
		return fmt.Errorf("cannot create message %s", ackName)
	}
	var rval = reflect.ValueOf(ack)
	var field = rval.Elem().FieldByName(ErrCodeField)
	if field.IsValid() && reflext.IsSignedInteger(field.Kind()) {
		field.SetInt(int64(ec))
	} else {
		return fmt.Errorf("message %s has no field named `%s`", ackName, ErrCodeField)
	}
	return m.Ack(ack)
}

var msgPool = pool.NewObjectPool[NetMessage]()

func AllocNetMessage() *NetMessage {
	return msgPool.Get()
}

func FreeNetMessage(netMsg *NetMessage) {
	netMsg.Reset()
	msgPool.Put(netMsg)
}

// DefaultMsgIDReflector get message ID by reflection
var DefaultMsgIDReflector = func(msg codec.Message) uint32 {
	var fullname string
	var rType = reflect.TypeOf(msg)
	if rType.Kind() == reflect.Ptr {
		fullname = rType.Elem().String()
	} else {
		fullname = rType.String()
	}
	return codec.NameHash(fullname)
}

// TryEnqueueMsg 尝试将消息放入队列，如果队列已满返回false
func TryEnqueueMsg(queue chan<- *NetMessage, msg *NetMessage) bool {
	select {
	case queue <- msg:
		return true
	default:
		return false // queue is full
	}
}

// TryDequeueMsg 尝试从队列中取出消息，如果队列为空返回nil
func TryDequeueMsg(queue <-chan *NetMessage) *NetMessage {
	select {
	case msg := <-queue:
		return msg
	default:
		return nil
	}
}
