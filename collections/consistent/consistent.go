// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package consistent

import (
	"sort"

	"gopkg.in/svrkit.v1/collections/slice"
)

const (
	ReplicaCount = 20 // 虚拟节点数量
)

// Consistent 一致性hash
// code inspired by github.com/stathat/consistent
type Consistent struct {
	circle     map[uint32]int32 // hash环
	nodes      map[int32]bool   // 所有节点
	sortedHash []uint32         // 环hash排序
}

func New() *Consistent {
	return &Consistent{
		circle: make(map[uint32]int32),
		nodes:  make(map[int32]bool),
	}
}

// AddNode 添加一个节点
func (c *Consistent) AddNode(node int32) {
	var hi = int64(node) << 32
	for i := 0; i < ReplicaCount; i++ {
		var hash = hashIntKey(hi | int64(i+1))
		c.circle[hash] = node
	}
	c.nodes[node] = true
	c.updateSortedHash()
}

func (c *Consistent) RemoveNode(node int32) {
	var hi = int64(node) << 32
	for i := 0; i < ReplicaCount; i++ {
		var hash = hashIntKey(hi | int64(i+1))
		delete(c.circle, hash)
	}
	delete(c.nodes, node)
	c.updateSortedHash()
}

func (c *Consistent) Clear() {
	clear(c.circle)
	clear(c.nodes)
	clear(c.sortedHash)
}

func (c *Consistent) Members() []int32 {
	var nodes = make([]int32, len(c.nodes))
	for node, _ := range c.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetNode 获取一个节点
func (c *Consistent) GetNode(key int64) int32 {
	if len(c.circle) == 0 {
		return 0
	}
	var i = c.search(hashIntKey(key))
	var hash = c.sortedHash[i]
	return c.circle[hash]
}

func (c *Consistent) GetNodeBy(key string) int32 {
	if len(c.circle) == 0 {
		return 0
	}
	var i = c.search(hashStrKey(key))
	var hash = c.sortedHash[i]
	return c.circle[hash]
}

// fnv hash
func hashIntKey(key int64) uint32 {
	// see src/hash/fnv.go sum32a.Write
	var hash = uint32(2166136261)
	for i := 0; i < 8; i++ {
		var ch = byte(key & 0xFF)
		hash ^= uint32(ch)
		hash *= 16777619
		key >>= 8
	}
	return hash
}

// fnv hash
func hashStrKey(key string) uint32 {
	// see src/hash/fnv.go sum32a.Write
	var hash = uint32(2166136261)
	for i := 0; i < len(key); i++ {
		var ch = key[i]
		hash ^= uint32(ch)
		hash *= 16777619
	}
	return hash
}

// 找到第一个大于等于`hash`的节点
func (c *Consistent) search(hash uint32) int {
	var i = sort.Search(len(c.sortedHash), func(x int) bool {
		return c.sortedHash[x] > hash
	})
	if i >= len(c.sortedHash) {
		i = 0
	}
	return i
}

func (c *Consistent) updateSortedHash() {
	hashes := c.sortedHash[:0]
	if cap(c.sortedHash)/(ReplicaCount*4) > len(c.circle) {
		hashes = nil // 使用率低于1/4重新分配内存
	}
	for k, _ := range c.circle {
		hashes = append(hashes, k)
	}
	sort.Sort(slice.Uint32Slice(hashes))
	c.sortedHash = hashes
}
