// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package consistent

import (
	"slices"
	"strconv"
)

const (
	ReplicaCount = 20 // 固定大小的虚拟节点数量
)

// Consistent 一致性hash
// code inspired by github.com/stathat/consistent
type Consistent struct {
	circle     map[uint32]string // <hash, 节点ID>
	nodes      map[string]bool   // 所有节点
	sortedHash []uint32          //
}

func New() *Consistent {
	return &Consistent{
		circle: make(map[uint32]string),
		nodes:  make(map[string]bool),
	}
}

func (c *Consistent) Len() int {
	return len(c.nodes)
}

func (c *Consistent) HasMember(node string) bool {
	return c.nodes[node]
}

func (c *Consistent) HasMemberNode(node uint64) bool {
	var key = strconv.FormatUint(node, 10)
	return c.nodes[key]
}

func consistentElementKey(elt string, idx int) string {
	return strconv.Itoa(idx) + elt
}

// fnv hash
func hashConsistentNodeKey(key string) uint32 {
	// see src/hash/fnv.go sum32a.Write
	var hash = uint32(2166136261)
	for i := 0; i < len(key); i++ {
		var ch = key[i]
		hash ^= uint32(ch)
		hash *= 16777619
	}
	return hash
}

// AddElem inserts an element node in the consistent hash.
func (c *Consistent) AddElem(node string) {
	if c.nodes[node] {
		return
	}
	for i := 0; i < ReplicaCount; i++ {
		var key = consistentElementKey(node, i)
		var hashCode = hashConsistentNodeKey(key)
		c.circle[hashCode] = node
	}
	c.nodes[node] = true
	c.updateSortedHash()
}

func (c *Consistent) Add(node uint64) {
	c.AddElem(strconv.FormatUint(node, 10))
}

// RemoveElem removes an element node from the hash.
func (c *Consistent) RemoveElem(node string) bool {
	if !c.nodes[node] {
		return false
	}
	for i := 0; i < ReplicaCount; i++ {
		var key = consistentElementKey(node, i)
		var hashCode = hashConsistentNodeKey(key)
		delete(c.circle, hashCode)
	}
	delete(c.nodes, node)
	c.updateSortedHash()
	return true
}

func (c *Consistent) Remove(node uint64) bool {
	return c.RemoveElem(strconv.FormatUint(node, 10))
}

func (c *Consistent) Clear() {
	clear(c.circle)
	clear(c.nodes)
	c.sortedHash = nil
}

func (c *Consistent) Members() []string {
	var nodes = make([]string, 0, len(c.nodes))
	for elem, _ := range c.nodes {
		nodes = append(nodes, elem)
	}
	return nodes
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
	// reallocate if we're holding on to too much (1/4th)
	if cap(c.sortedHash)/(ReplicaCount*4) > len(c.circle) {
		hashes = nil
	}
	for k, _ := range c.circle {
		hashes = append(hashes, k)
	}
	slices.Sort(hashes)
	c.sortedHash = hashes
}

// Get returns an element close to where name hashes to in the circle.
func (c *Consistent) Get(name string) string {
	if len(c.circle) == 0 {
		return ""
	}
	var key = hashConsistentNodeKey(name)
	var i = c.search(key)
	var hash = c.sortedHash[i]
	return c.circle[hash]
}

func (c *Consistent) GetNode(node uint64) uint64 {
	if len(c.circle) == 0 {
		return 0
	}
	var key = strconv.FormatUint(node, 10)
	var elem = c.Get(key)
	if elem != "" {
		n, _ := strconv.ParseUint(elem, 10, 64)
		return n
	}
	return 0
}

// GetTwo returns the two closest distinct elements to the name input in the circle.
func (c *Consistent) GetTwo(name string) []string {
	if len(c.circle) == 0 {
		return nil
	}
	var key = hashConsistentNodeKey(name)
	var i = c.search(key)
	var hash = c.sortedHash[i]
	var a = c.circle[hash]
	if len(c.circle) == 1 {
		return []string{a}
	}
	var b string
	var start = i
	for i = start + 1; i != start; i++ {
		if i >= len(c.sortedHash) {
			i = 0
		}
		b = c.circle[c.sortedHash[i]]
		if b != a {
			break
		}
	}
	return []string{a, b}
}

func parseIntNodes(ss ...string) []uint64 {
	var list = make([]uint64, 0, len(ss))
	for _, s := range ss {
		if n, err := strconv.ParseUint(s, 10, 64); err == nil {
			list = append(list, n)
		}
	}
	return list
}

// GetTwoNodes returns the two closest distinct elements to the name input in the circle.
func (c *Consistent) GetTwoNodes(node uint64) []uint64 {
	var key = strconv.FormatUint(node, 10)
	var elem = c.GetTwo(key)
	return parseIntNodes(elem...)
}
