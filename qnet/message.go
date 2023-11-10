// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"encoding/binary"
	"hash/crc32"
	"reflect"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"gopkg.in/svrkit.v1/pool"
)

const (
	MsgArenaPoolSize = 1024
)

var msgAlloc = pool.NewArenaAllocator[NetMessage](MsgArenaPoolSize)

func AllocNetMessage() *NetMessage {
	return msgAlloc.Alloc()
}

type SessionMessage struct {
	Session Endpoint
	MsgId   uint32
	MsgBody proto.Message
}

type NetMessage struct {
	Command   uint32        `json:"cmd"`
	Route     uint32        `json:"route,omitempty"`
	Seq       uint32        `json:"seq,omitempty"`
	Errno     uint32        `json:"errno,omitempty"`
	Data      []byte        `json:"data,omitempty"`
	CreatedAt int64         `json:"created_at,omitempty"`
	Body      proto.Message `json:"body,omitempty"`
	Session   Endpoint      `json:"-"`
}

func NewNetMessage(cmd, seq uint32, data []byte) *NetMessage {
	var msg = AllocNetMessage()
	msg.Command = cmd
	msg.Seq = seq
	msg.Data = data
	return msg
}

func CreateNetMessage(cmd, seq uint32, body proto.Message) *NetMessage {
	var msg = AllocNetMessage()
	msg.Command = cmd
	msg.Seq = seq
	msg.Body = body
	return msg
}

func CreateNetMessageWith(body proto.Message) *NetMessage {
	var msg = AllocNetMessage()
	msg.Command = DefaultMsgIDReflector(body)
	msg.Body = body
	return nil
}

func (m *NetMessage) Reset() {
	*m = NetMessage{}
}

func (m *NetMessage) Clone() *NetMessage {
	return &NetMessage{
		CreatedAt: time.Now().UnixNano() / 1e6,
		Seq:       m.Seq,
		Errno:     m.Errno,
		Command:   m.Command,
		Body:      m.Body,
	}
}

func (m *NetMessage) SetMsgID(msgId uint32) {
	m.Command = msgId
}

func (m *NetMessage) ErrCode() uint32 {
	return m.Errno
}

// Encode encode `Body` to `Data`
func (m *NetMessage) Encode() error {
	if m.Data != nil {
		return nil
	}
	if m.Errno != 0 {
		var buf [binary.MaxVarintLen32]byte
		var i = binary.PutUvarint(buf[:], uint64(m.Errno))
		m.Data = buf[:i]
	} else if m.Body != nil {
		data, err := proto.Marshal(m.Body)
		if err != nil {
			return err
		}
		m.Data = data
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

func (m *NetMessage) Refuse(ec uint32) error {
	var ack = AllocNetMessage()
	ack.Seq = m.Seq
	ack.Errno = ec
	return m.Session.SendMsg(ack, SendNonblock)
}

func (m *NetMessage) ReplyAck(ack proto.Message) error {
	var netMsg = CreateNetMessageWith(ack)
	netMsg.Seq = m.Seq
	return m.Session.SendMsg(netMsg, SendNonblock)
}

func (m *NetMessage) Reply(cmd uint32, data []byte) error {
	var netMsg = AllocNetMessage()
	netMsg.Seq = m.Seq
	netMsg.Command = cmd
	netMsg.Data = data
	return m.Session.SendMsg(netMsg, SendNonblock)
}

// DefaultMsgIDReflector get message ID by reflection
var DefaultMsgIDReflector = func(msg proto.Message) uint32 {
	var name = reflect.TypeOf(msg).String()
	var idx = strings.LastIndex(name, ".") // *protos.LoginReq --> LoginReq
	if idx > 0 {
		name = name[idx+1:]
	}
	var crc = crc32.NewIEEE()
	crc.Write([]byte(name))
	return crc.Sum32()
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
