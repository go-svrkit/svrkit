// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"sync/atomic"
)

// StreamConnBase base stream connection
type StreamConnBase struct {
	Node       NodeID             //
	Running    atomic.Bool        //
	RecvQueue  chan<- *NetMessage // 收消息队列
	SendQueue  chan *NetMessage   // 发消息队列
	ErrChan    chan *Error        // error signal
	Encrypt    Encryptor          // 加密
	Decrypt    Encryptor          // 解密
	RemoteAddr string             // 缓存的远端地址
	Userdata   any                // user data
}

func (c *StreamConnBase) IsRunning() bool {
	return c.Running.Load()
}

func (c *StreamConnBase) GetNode() NodeID {
	return c.Node
}

func (c *StreamConnBase) SetNode(node NodeID) {
	c.Node = node
}

func (c *StreamConnBase) GetRemoteAddr() string {
	return c.RemoteAddr
}

func (c *StreamConnBase) SetUserData(val any) {
	c.Userdata = val
}

func (c *StreamConnBase) GetUserData() any {
	return c.Userdata
}

func (c *StreamConnBase) SetSendQueue(sendQueue chan *NetMessage) {
	c.SendQueue = sendQueue
}

func (c *StreamConnBase) SetEncryption(encrypt, decrypt Encryptor) {
	c.Encrypt = encrypt
	c.Decrypt = decrypt
}

func (c *StreamConnBase) SendNonBlock(msg *NetMessage) bool {
	select {
	case c.SendQueue <- msg:
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
		c.SendQueue <- msg // blocking send
		return nil
	}
}
