// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// Node 表示用于服务发现的节点信息
type Node struct {
	Type      string            `json:"type"`
	ID        uint32            `json:"id"`
	Status    uint32            `json:"status,omitempty"`
	PID       int               `json:"pid,omitempty"`
	Host      string            `json:"host,omitempty"`
	Interface string            `json:"interface,omitempty"`
	URI       string            `json:"uri,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
}

func NewNode(nodeType string, id uint32) Node {
	node := Node{
		Type: nodeType,
		ID:   id,
		PID:  os.Getpid(),
	}
	if hostname, err := os.Hostname(); err == nil {
		node.Host = hostname
	}
	return node
}

func (n *Node) Get(key string) string {
	return n.Data[key]
}

func (n *Node) Set(key, val string) {
	if n.Data == nil {
		n.Data = make(map[string]string)
	}
	n.Data[key] = val
}

func (n *Node) GetInt(key string) int {
	var s = n.Data[key]
	val, _ := strconv.Atoi(s)
	return val
}

func (n *Node) SetInt(key string, val int) {
	var s = strconv.Itoa(val)
	n.Set(key, s)
}

func (n *Node) GetBool(key string) bool {
	var s = n.Data[key]
	val, _ := strconv.ParseBool(s)
	return val
}

func (n *Node) SetBool(key string, val bool) {
	var s = strconv.FormatBool(val)
	n.Set(key, s)
}

func (n *Node) GetFloat(key string) float64 {
	var s = n.Data[key]
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func (n *Node) SetFloat(key string, val float64) {
	var s = strconv.FormatFloat(val, 'f', 3, 64)
	n.Set(key, s)
}

func (n *Node) Clone() Node {
	var clone = *n
	clone.Data = make(map[string]string, len(n.Data))
	for k, v := range n.Data {
		clone.Data[k] = v
	}
	return clone
}

func (n *Node) String() string {
	data, _ := json.Marshal(n)
	return string(data)
}

type NodeEventType int

const (
	EventUnknown NodeEventType = 0
	EventCreate  NodeEventType = 1
	EventUpdate  NodeEventType = 2
	EventDelete  NodeEventType = 3
)

func (e NodeEventType) String() string {
	switch e {
	case EventCreate:
		return "create"
	case EventUpdate:
		return "update"
	case EventDelete:
		return "delete"
	}
	return "???"
}

// NodeEvent 节点变化事件
type NodeEvent struct {
	Type NodeEventType
	Key  string
	Node Node
}

func (e NodeEvent) String() string {
	return fmt.Sprintf("%v %s: %v", e.Type, e.Key, e.Node)
}

// NodeSet 节点信息列表
type NodeSet []Node

// NodeMap 按服务类型区分的节点信息
type NodeMap struct {
	guard sync.RWMutex
	nodes map[string]NodeSet
}

func NewNodeMap() *NodeMap {
	return &NodeMap{
		nodes: make(map[string]NodeSet),
	}
}

// Count 所有节点数量
func (m *NodeMap) Count() int {
	m.guard.RLock()
	var count = 0
	for _, nodes := range m.nodes {
		count += len(nodes)
	}
	m.guard.RUnlock()
	return count
}

func (m *NodeMap) GetKeys() []string {
	m.guard.RLock()
	var names = make([]string, 0, len(m.nodes))
	for name := range m.nodes {
		names = append(names, name)
	}
	m.guard.RUnlock()
	return names
}

// GetNodes 所有本类型的节点，不要修改返回值
func (m *NodeMap) GetNodes(nodeType string) NodeSet {
	m.guard.RLock()
	v := m.nodes[nodeType]
	m.guard.RUnlock()
	return v
}

// InsertNode 添加一个节点
func (m *NodeMap) InsertNode(node Node) {
	m.guard.Lock()
	defer m.guard.Unlock()

	slice := m.nodes[node.Type]
	for i, v := range slice {
		if v.ID == node.ID {
			slice[i] = node
			return
		}
	}
	m.nodes[node.Type] = append(slice, node)
}

func (m *NodeMap) Clear() {
	m.guard.Lock()
	m.nodes = make(map[string]NodeSet)
	m.guard.Unlock()
}

// DeleteNodes 删除某一类型的所有节点
func (m *NodeMap) DeleteNodes(nodeType string) {
	m.guard.Lock()
	m.nodes[nodeType] = nil
	m.guard.Unlock()
}

// DeleteNode 删除一个节点
func (m *NodeMap) DeleteNode(nodeType string, id uint32) {
	m.guard.Lock()
	defer m.guard.Unlock()

	var a = m.nodes[nodeType]
	var idx = -1
	for i, v := range a {
		if v.ID == id {
			idx = i
			break
		}
	}
	if idx >= 0 {
		m.nodes[nodeType] = slices.Delete(a, idx, idx+1)
		if len(m.nodes[nodeType]) == 0 {
			delete(m.nodes, nodeType)
		}
	}
}

func (m *NodeMap) String() string {
	var sb strings.Builder
	for name, set := range m.nodes {
		fmt.Fprintf(&sb, "%s: %v,\n", name, set)
	}
	return sb.String()
}

// 使用bigint序列化大整数
func unmarshalNode(data []byte, node *Node) error {
	if len(data) > 0 {
		var dec = json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		return dec.Decode(node)
	}
	return nil
}
