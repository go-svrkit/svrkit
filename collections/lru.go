// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package collections

import (
	"log"
)

type LRUEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// LRUCache LRU缓存
// https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)
type LRUCache[K comparable, V any] struct {
	list      *List[LRUEntry[K, V]]
	items     map[K]*ListElem[LRUEntry[K, V]]
	onEvicted func(key K, value V)
	size      int
}

func NewLRUCache[K comparable, V any](size int, onEvicted func(key K, value V)) *LRUCache[K, V] {
	if size <= 0 {
		log.Panicln("cache capacity out of range")
	}
	cache := &LRUCache[K, V]{
		size:      size,
		onEvicted: onEvicted,
		list:      NewList[LRUEntry[K, V]](),
		items:     make(map[K]*ListElem[LRUEntry[K, V]], size),
	}
	return cache
}

func (c *LRUCache[K, V]) Len() int {
	return c.list.Len()
}

func (c *LRUCache[K, V]) Cap() int {
	return c.size
}

// Contains 查看key是否存在，不移动链表
func (c *LRUCache[K, V]) Contains(key K) bool {
	_, found := c.items[key]
	return found
}

// Get 获取key对应的值，并把其移动到链表头部
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	entry, found := c.items[key]
	if found {
		c.list.MoveToFront(entry)
		kv := entry.Value
		return kv.Value, true
	}
	var zero V
	return zero, false
}

// Peek 获取key对应的值，不移动链表
func (c *LRUCache[K, V]) Peek(key K) (V, bool) {
	elem, found := c.items[key]
	if found && elem != nil {
		entry := elem.Value
		return entry.Value, true
	}
	var zero V
	return zero, false
}

// GetOldest 获取最老的值（链表尾节点）
func (c *LRUCache[K, V]) GetOldest() (key K, value V, ok bool) {
	var elem = c.list.Back()
	if elem != nil {
		entry := elem.Value
		return entry.Key, entry.Value, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Keys 返回所有的key（从旧到新）
func (c *LRUCache[K, V]) Keys() []K {
	var keys = make([]K, 0, len(c.items))
	for elem := c.list.Back(); elem != nil; elem = elem.Prev() {
		var key = elem.Value.Key
		keys = append(keys, key)
	}
	return keys
}

// Put 把key-value加入到cache中，并移动到链表头部
func (c *LRUCache[K, V]) Put(key K, value V) bool {
	elem, exist := c.items[key]
	if exist {
		c.list.MoveToFront(elem)
		elem.Value.Value = value
		return false
	}
	entry := LRUEntry[K, V]{Key: key, Value: value}
	elem = c.list.PushFront(entry) // push entry to list front
	c.items[key] = elem
	if c.Len() > c.size {
		c.removeOldest()
	}
	return true
}

// Resize changes the cache size.
func (c *LRUCache[K, V]) Resize(size int) int {
	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.removeOldest()
	}
	c.size = size
	return diff
}

// Remove 把key从cache中删除
func (c *LRUCache[K, V]) Remove(key K) bool {
	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
		return true
	}
	return false
}

// RemoveOldest 删除最老的的key-Value，并返回
func (c *LRUCache[K, V]) RemoveOldest() (key K, value V, ok bool) {
	elem := c.list.Back()
	if elem != nil {
		c.removeElement(elem)
		var entry = elem.Value
		return entry.Key, entry.Value, true
	}
	return
}

// Purge 清除所有
func (c *LRUCache[K, V]) Purge() {
	for k, elem := range c.items {
		if c.onEvicted != nil {
			c.onEvicted(elem.Value.Key, elem.Value.Value)
		}
		delete(c.items, k)
	}
	c.list.Init()
}

// removeOldest removes the oldest item from the cache.
func (c *LRUCache[K, V]) removeOldest() {
	entry := c.list.Back()
	if entry != nil {
		c.removeElement(entry)
	}
}

// remove a given list element from the cache
func (c *LRUCache[K, V]) removeElement(elem *ListElem[LRUEntry[K, V]]) {
	entry := elem.Value
	c.list.Remove(elem)
	delete(c.items, entry.Key)
	if c.onEvicted != nil {
		c.onEvicted(entry.Key, entry.Value)
	}
}
