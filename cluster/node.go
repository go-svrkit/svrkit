// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package cluster

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
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
	Args      map[string]string `json:"args,omitempty"`
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

func (n *Node) GetStr(key string) string {
	return n.Args[key]
}

func (n *Node) SetStr(key, val string) {
	if n.Args == nil {
		n.Args = make(map[string]string)
	}
	n.Args[key] = val
}

func (n *Node) Set(key string, val any) {
	if n.Args == nil {
		n.Args = make(map[string]string)
	}
	data, _ := sonic.MarshalString(val)
	n.Args[key] = data
}

func (n *Node) GetInt(key string) int {
	var s = n.Args[key]
	val, _ := strconv.Atoi(s)
	return val
}

func (n *Node) SetInt(key string, val int) {
	var s = strconv.Itoa(val)
	n.SetStr(key, s)
}

func (n *Node) GetBool(key string) bool {
	var s = n.Args[key]
	val, _ := strconv.ParseBool(s)
	return val
}

func (n *Node) SetBool(key string, val bool) {
	var s = strconv.FormatBool(val)
	n.SetStr(key, s)
}

func (n *Node) GetFloat(key string) float64 {
	var s = n.Args[key]
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func (n *Node) SetFloat(key string, val float64) {
	var s = strconv.FormatFloat(val, 'f', 5, 64)
	n.SetStr(key, s)
}

func (n *Node) Clone() Node {
	var clone = *n
	clone.Args = maps.Clone(n.Args)
	return clone
}

func (n *Node) String() string {
	data, _ := sonic.MarshalString(n)
	return data
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

// NodeMap 按服务类型区分的节点信息
type NodeMap struct {
	sync.RWMutex
	nodes map[string][]Node
}

func NewNodeMap() *NodeMap {
	return &NodeMap{
		nodes: make(map[string][]Node),
	}
}

// Count 所有节点数量
func (m *NodeMap) Count() int {
	m.RLock()
	var count = 0
	for _, nodes := range m.nodes {
		count += len(nodes)
	}
	m.RUnlock()
	return count
}

func (m *NodeMap) CountOf(nodeType string) int {
	m.RLock()
	v := m.nodes[nodeType]
	m.RUnlock()
	return len(v)
}

func (m *NodeMap) GetKeys() []string {
	m.RLock()
	var names = make([]string, 0, len(m.nodes))
	for name := range m.nodes {
		names = append(names, name)
	}
	m.RUnlock()
	return names
}

// GetNodes 所有本类型的节点，不要修改返回值
func (m *NodeMap) GetNodes(nodeType string) []Node {
	m.RLock()
	v := m.nodes[nodeType]
	m.RUnlock()
	return v
}

func (m *NodeMap) FindNodeOf(nodeType string, id uint32) int {
	var nodes = m.nodes[nodeType]
	for i := 0; i < len(nodes); i++ {
		if nodes[i].ID == id {
			return i
		}
	}
	return -1
}

// AddNode 添加一个节点
func (m *NodeMap) AddNode(node Node) {
	m.Lock()
	defer m.Unlock()

	var nodes = m.nodes[node.Type]
	for i, v := range nodes {
		if v.ID == node.ID {
			nodes[i] = node
			return
		}
	}
	m.nodes[node.Type] = append(nodes, node)
}

func (m *NodeMap) Clear() {
	m.Lock()
	m.nodes = make(map[string][]Node)
	m.Unlock()
}

// DeleteNodesOf 删除某一类型的所有节点
func (m *NodeMap) DeleteNodesOf(nodeType string) {
	m.Lock()
	m.nodes[nodeType] = nil
	m.Unlock()
}

// DeleteNode 删除一个节点
func (m *NodeMap) DeleteNode(nodeType string, id uint32) {
	m.Lock()
	defer m.Unlock()

	var idx = m.FindNodeOf(nodeType, id)
	if idx >= 0 {
		var nodes = m.nodes[nodeType]
		nodes = slices.Delete(nodes, idx, idx+1)
		if len(m.nodes[nodeType]) == 0 {
			delete(m.nodes, nodeType)
		} else {
			m.nodes[nodeType] = nodes
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
		var dec = decoder.NewStreamDecoder(bytes.NewReader(data))
		dec.UseInt64()
		return dec.Decode(node)
	}
	return nil
}
