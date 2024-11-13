// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

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
// 一个64位整数表示的节点号，用以标识一个service，低32位为服务实例编号，32-48位为服务类型；
// 或者一个客户端session，低32位为GATE内部的session编号，32-48位为GATE编号；
type NodeID uint64

// MakeBackendNode 根据服务号和实例号创建一个节点ID
func MakeBackendNode(service uint16, instance uint32) NodeID {
	return NodeBackendTypeMask | NodeID((uint64(service)<<NodeInstanceShift)|uint64(instance))
}

// MakeGateSession `instance`指GATE的实例编号，限定为16位
func MakeGateSession(instance uint16, session uint32) NodeID {
	return NodeID((uint64(instance) << NodeInstanceShift) | uint64(session))
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
func (n NodeID) Service() int16 {
	return int16(n >> NodeInstanceShift)
}

// Instance 节点的实例编号
func (n NodeID) Instance() uint32 {
	return uint32(n)
}

// GateID client会话的网关ID
func (n NodeID) GateID() uint16 {
	return uint16(n >> NodeInstanceShift)
}

func (n NodeID) Session() uint32 {
	return uint32(n)
}

func (n NodeID) String() string {
	if n.IsSession() {
		return fmt.Sprintf("G%d#%d", n.GateID(), n.Session())
	}
	return fmt.Sprintf("%02X#%d", n.Service(), n.Instance())
}

// NodeIDSet 没有重复ID的有序集合
type NodeIDSet = []NodeID
