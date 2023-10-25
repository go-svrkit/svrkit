// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"time"

	"github.com/golang/protobuf/proto"
	"gopkg.in/svrkit.v1/pool"
)

const (
	MsgArenaPoolSize = 1024
)

var msgAlloc = pool.NewArenaAllocator[NetMessage](MsgArenaPoolSize)

func AllocNetMessage() *NetMessage {
	var message = msgAlloc.Alloc()
	message.CreatedAt = time.Now().UnixNano() / 1e6
	return message
}

type SessionMessage struct {
	Session Endpoint
	MsgId   uint32
	MsgBody proto.Message
}

// NetMessage 投递给业务层的网络消息
type NetMessage struct {
	MsgID     uint32        `json:"cmd"`            // 消息ID
	Seq       uint32        `json:"seq"`            // 序列号
	Errno     int32         `json:"errno"`          // 错误码
	CreatedAt int64         `json:"created_at"`     // 创建时间(毫秒)
	Body      proto.Message `json:"body,omitempty"` // pb结构体
	Data      []byte        `json:"data,omitempty"` // raw binary data
	Session   Endpoint      `json:"-"`              //
}

func NewNetMessage(seq uint32, body proto.Message) *NetMessage {
	var msg = AllocNetMessage()
	msg.MsgID = DefaultMsgIDReflector(body)
	msg.Body = body
	msg.Seq = seq
	return msg
}

func (m *NetMessage) Reset() {
	*m = NetMessage{}
}

func (m *NetMessage) Clone() *NetMessage {
	return &NetMessage{
		CreatedAt: time.Now().UnixNano() / 1e6,
		Seq:       m.Seq,
		Errno:     m.Errno,
		MsgID:     m.MsgID,
		Body:      m.Body,
	}
}

func (m *NetMessage) SetMsgID(msgId uint32) {
	m.MsgID = msgId
}

func (m *NetMessage) ErrCode() int32 {
	return m.Errno
}

// Encode encode `Body` to `Data`
func (m *NetMessage) Encode() error {
	if m.Data == nil && m.Body != nil {
		data, err := proto.Marshal(m.Body)
		if err != nil {
			return err
		}
		m.Data = data
		m.Body = nil
		return nil
	}
	return nil
}

// DecodeTo decode `Data` to `msg`
func (m *NetMessage) DecodeTo(msg proto.Message) error {
	if err := proto.Unmarshal(m.Data, msg); err != nil {
		return err
	}
	m.Data = nil
	return nil
}

func (m *NetMessage) Refuse(ec int32) error {
	var ack = AllocNetMessage()
	ack.Seq = m.Seq
	ack.Errno = ec
	return m.Session.SendMsg(ack, SendNonblock)
}

func (m *NetMessage) ReplyAck(ack proto.Message) error {
	var netMsg = NewNetMessage(m.Seq, ack)
	return m.Session.SendMsg(netMsg, SendNonblock)
}

func (m *NetMessage) Reply(msgId uint32, data []byte) error {
	var netMsg = AllocNetMessage()
	netMsg.Seq = m.Seq
	netMsg.MsgID = msgId
	netMsg.Data = data
	return m.Session.SendMsg(netMsg, SendNonblock)
}

// DefaultMsgIDReflector get message ID by reflection
var DefaultMsgIDReflector = func(proto.Message) uint32 {
	panic("not implemented")
	return 0
}
