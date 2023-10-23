// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"fmt"
)

const (
	NodeInstanceShift          = 32
	NodeTypeShift              = 48
	MaxNodeInstance            = (1 << NodeInstanceShift) - 1
	NodeBackendTypeMask NodeID = 1 << NodeTypeShift
)

// NodeID 节点ID
// 一个64位整数表示的节点号，用以标识一个service(最高位为1)，低32位为服务实例编号，32-48位为服务类型；
// 或者一个客户端session(最高位为0)，低32位为GATE内部的session编号，32-48位为GATE编号；
type NodeID uint64

// NewBackendNode 根据服务号和实例号创建一个节点ID
func NewBackendNode(service uint16, instance uint32) NodeID {
	return NodeBackendTypeMask | NodeID((uint64(service)<<NodeInstanceShift)|uint64(instance))
}

// IsBackend 是否backend节点
func (n NodeID) IsBackend() bool {
	return n > NodeBackendTypeMask
}

// IsSession 是否client会话
func (n NodeID) IsSession() bool {
	return n < NodeBackendTypeMask
}

// Service 服务型
func (n NodeID) Service() uint16 {
	return uint16((n >> NodeInstanceShift) & 0xFF)
}

// GateID client会话的网关ID
func (n NodeID) GateID() uint16 {
	return uint16((n >> NodeInstanceShift) & 0xFF)
}

// Instance service节点的实例编号
func (n NodeID) Instance() uint32 {
	return uint32(n)
}

func (n NodeID) String() string {
	if n.IsSession() {
		return fmt.Sprintf("G%d#%d", n.GateID(), n.Instance())
	}
	return fmt.Sprintf("%02x#%d", n.Service(), n.Instance())
}

// NodeIDSet 没有重复ID的有序集合
type NodeIDSet = []NodeID
