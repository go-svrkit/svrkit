// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"sync/atomic"
)

// StreamConnBase base stream connection
type StreamConnBase struct {
	node      NodeID             //
	running   atomic.Int32       //
	recvQueue chan<- *NetMessage // 收消息队列
	sendQueue chan *NetMessage   // 发消息队列
	errChan   chan *Error        // error signal
	encrypt   Encryptor          // 加密
	decrypt   Encryptor          // 解密
	addr      string             // 缓存的远端地址
	userdata  any                // user data
}

func (c *StreamConnBase) Init(node NodeID, recvQueue chan<- *NetMessage, errChan chan *Error) {
	c.node = node
	c.recvQueue = recvQueue
	c.errChan = errChan
}

func (c *StreamConnBase) IsRunning() bool {
	return c.running.Load() != 0
}

func (c *StreamConnBase) Node() NodeID {
	return c.node
}

func (c *StreamConnBase) SetNode(node NodeID) {
	c.node = node
}

func (c *StreamConnBase) SetSendQueue(sendQueue chan *NetMessage) {
	c.sendQueue = sendQueue
}

func (c *StreamConnBase) SetRemoteAddr(addr string) {
	c.addr = addr
}

func (c *StreamConnBase) RemoteAddr() string {
	return c.addr
}

func (c *StreamConnBase) SetUserData(ud any) {
	c.userdata = ud
}

func (c *StreamConnBase) UserData() any {
	return c.userdata
}

func (c *StreamConnBase) SetEncryption(encrypt, decrypt Encryptor) {
	c.encrypt = encrypt
	c.decrypt = decrypt
}

func (c *StreamConnBase) SendNonBlock(msg *NetMessage) bool {
	select {
	case c.sendQueue <- msg:
		return true
	default:
		return false
	}
}

func (c *StreamConnBase) SendMsg(msg *NetMessage, mode int) error {
	if !c.IsRunning() {
		return ErrConnNotRunning
	}
	switch mode {
	case SendNonblock:
		if !c.SendNonBlock(msg) {
			return ErrConnOutboundOverflow
		}
		return nil
	default:
		c.sendQueue <- msg // blocking send
		return nil
	}
}
