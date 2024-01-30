// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package consistent

import (
	"slices"
)

const (
	ReplicaCount = 20 // 固定大小的虚拟节点数量
)

// Consistent 一致性hash
// code inspired by github.com/stathat/consistent
type Consistent struct {
	circle     map[uint32]int32 // <hash, 节点ID>
	nodes      map[int32]bool   // 所有节点
	sortedHash []uint32         //
}

func New() *Consistent {
	return &Consistent{
		circle: make(map[uint32]int32),
		nodes:  make(map[int32]bool),
	}
}

func (c *Consistent) Len() int {
	return len(c.nodes)
}

// eltKey generates a string key for an element with an index.
func (c *Consistent) eltKey(elt int32, idx int) int64 {
	return int64(elt)<<32 + int64(idx)
}

// Add inserts an int32 node in the consistent hash.
func (c *Consistent) Add(nodes ...int32) {
	if len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		for i := 0; i < ReplicaCount; i++ {
			var key = c.eltKey(node, i)
			c.circle[hashIntKey(key)] = node
		}
		c.nodes[node] = true
	}
	c.updateSortedHash()
}

// Remove removes an int32 node from the hash.
func (c *Consistent) Remove(nodes ...int32) {
	for _, node := range nodes {
		for i := 0; i < ReplicaCount; i++ {
			var key = c.eltKey(node, i)
			delete(c.circle, hashIntKey(key))
		}
		delete(c.nodes, node)
	}
	c.updateSortedHash()
}

func (c *Consistent) Clear() {
	clear(c.circle)
	clear(c.nodes)
	clear(c.sortedHash)
}

func (c *Consistent) Members() []int32 {
	var nodes = make([]int32, len(c.nodes))
	for elem, _ := range c.nodes {
		nodes = append(nodes, elem)
	}
	return nodes
}

// Get returns an element close to where name hashes to in the circle.
func (c *Consistent) Get(name string) int32 {
	if len(c.circle) == 0 {
		return 0
	}
	var i = c.search(hashKey(name))
	var hash = c.sortedHash[i]
	return c.circle[hash]
}

// GetBy returns an element close to where name hashes to in the circle.
func (c *Consistent) GetBy(key int64) int32 {
	if len(c.circle) == 0 {
		return 0
	}
	var i = c.search(hashIntKey(key))
	return c.circle[c.sortedHash[i]]
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
func hashKey(key string) uint32 {
	// see src/hash/fnv.go sum32a.Write
	var hash = uint32(2166136261)
	for i := 0; i < len(key); i++ {
		var ch = key[i]
		hash ^= uint32(ch)
		hash *= 16777619
	}
	return hash
}

func (c *Consistent) search(hash uint32) int {
	i, _ := slices.BinarySearch(c.sortedHash, hash)
	if i >= len(c.sortedHash) {
		i = 0 // fallback to first node
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
	slices.Sort(hashes)
	c.sortedHash = hashes
}
