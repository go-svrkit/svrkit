// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"context"
	"net"
	"sync"
)

const (
	SendBlock    = 0
	SendNonblock = 1
)

const (
	DefaultRecvQueueSize        = 1 << 16
	DefaultBackendSendQueueSize = 1 << 14
	DefaultSessionSendQueueSize = 1024
	DefaultBacklogSize          = 128
	DefaultErrorChanSize        = 64
)

type Encryptor interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

// Endpoint 网络端点
type Endpoint interface {
	GetNode() NodeID
	SetNode(NodeID)

	GetRemoteAddr() string
	UnderlyingConn() net.Conn

	GetUserData() any
	SetEncryption(Encryptor, Encryptor)
	SetSendQueue(chan *NetMessage)

	Go(ctx context.Context, reader, writer bool) // 开启read/write线程

	SendMsg(*NetMessage, int) error
	Close() error
	ForceClose(error)
}

// EndpointMap 线程安全的Endpoint Map
type EndpointMap struct {
	guard     sync.RWMutex
	endpoints map[NodeID]Endpoint
}

func NewEndpointMap() *EndpointMap {
	return new(EndpointMap).Init()
}

func (m *EndpointMap) Init() *EndpointMap {
	m.endpoints = make(map[NodeID]Endpoint, 64)
	return m
}

func (m *EndpointMap) Size() int {
	m.guard.RLock()
	var n = len(m.endpoints)
	m.guard.RUnlock()
	return n
}

func (m *EndpointMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *EndpointMap) Has(node NodeID) bool {
	m.guard.RLock()
	_, ok := m.endpoints[node]
	m.guard.RUnlock()
	return ok
}

func (m *EndpointMap) Get(node NodeID) Endpoint {
	m.guard.RLock()
	var v = m.endpoints[node]
	m.guard.RUnlock()
	return v
}

func (m *EndpointMap) Keys() []NodeID {
	m.guard.RLock()
	var keys = make([]NodeID, 0, len(m.endpoints))
	for k, _ := range m.endpoints {
		keys = append(keys, k)
	}
	m.guard.RUnlock()
	return keys
}

// Foreach action应该对map是read-only
func (m *EndpointMap) Foreach(action func(Endpoint) bool) {
	m.guard.RLock()
	defer m.guard.RUnlock()
	for _, endpoint := range m.endpoints {
		if !action(endpoint) {
			break
		}
	}
}

func (m *EndpointMap) Put(node NodeID, endpoint Endpoint) {
	m.guard.Lock()
	m.endpoints[node] = endpoint
	m.guard.Unlock()
}

func (m *EndpointMap) PutIfAbsent(node NodeID, endpoint Endpoint) {
	m.guard.Lock()
	if _, ok := m.endpoints[node]; !ok {
		m.endpoints[node] = endpoint
	}
	m.guard.Unlock()
}

func (m *EndpointMap) Remove(node NodeID) {
	m.guard.Lock()
	delete(m.endpoints, node)
	m.guard.Unlock()
}

func (m *EndpointMap) Clear() {
	m.guard.Lock()
	clear(m.endpoints)
	m.guard.Unlock()
}
