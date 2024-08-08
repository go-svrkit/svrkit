// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package lru

import (
	"container/list"
)

type Entry[T comparable] struct {
	Key   T
	Value any
}

// Cache LRU缓存
// https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)
type Cache[T comparable] struct {
	list      *list.List
	items     map[T]*list.Element
	onEvicted func(key T, value any)
	size      int
}

func NewCache[T comparable](size int, onEvicted func(key T, value any)) *Cache[T] {
	if size <= 0 {
		panic("cache capacity out of range")
	}
	cache := &Cache[T]{
		size:      size,
		onEvicted: onEvicted,
		list:      list.New(),
		items:     make(map[T]*list.Element, size),
	}
	return cache
}

func (c *Cache[T]) Len() int {
	return c.list.Len()
}

func (c *Cache[T]) Cap() int {
	return c.size
}

// Contains 查看key是否存在，不移动链表
func (c *Cache[T]) Contains(key T) bool {
	_, found := c.items[key]
	return found
}

// Get 获取key对应的值，并把其移动到链表头部
func (c *Cache[T]) Get(key T) (any, bool) {
	e, found := c.items[key]
	if found {
		c.list.MoveToFront(e)
		kv := e.Value.(*Entry[T])
		if kv == nil {
			return nil, false
		}
		return kv.Value, true
	}
	return nil, false
}

// Peek 获取key对应的值，不移动链表
func (c *Cache[T]) Peek(key T) (any, bool) {
	e, found := c.items[key]
	if found {
		kv := e.Value.(*Entry[T])
		return kv.Value, true
	}
	return nil, false
}

// GetOldest 获取最老的值（链表尾节点）
func (c *Cache[T]) GetOldest() (key T, value any, ok bool) {
	ent := c.list.Back()
	if ent != nil {
		kv := ent.Value.(*Entry[T])
		return kv.Key, kv.Value, true
	}
	var empty T
	return empty, nil, false
}

// Keys 返回所有的key（从旧到新）
func (c *Cache[T]) Keys() []T {
	keys := make([]T, len(c.items))
	i := 0
	for e := c.list.Back(); e != nil; e = e.Prev() {
		keys[i] = e.Value.(*Entry[T]).Key
		i++
	}
	return keys
}

// Put 把key-value加入到cache中，并移动到链表头部
func (c *Cache[T]) Put(key T, value any) bool {
	e, exist := c.items[key]
	if exist {
		c.list.MoveToFront(e)
		e.Value.(*Entry[T]).Value = value
		return false
	}
	entry := &Entry[T]{Key: key, Value: value}
	e = c.list.PushFront(entry) // push entry to list front
	c.items[key] = e
	if c.Len() > c.size {
		c.removeOldest()
	}
	return true
}

// Resize changes the cache size.
func (c *Cache[T]) Resize(size int) int {
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
func (c *Cache[T]) Remove(key T) bool {
	if e, ok := c.items[key]; ok {
		c.removeElement(e)
		return true
	}
	return false
}

// RemoveOldest 删除最老的的key-value，并返回
func (c *Cache[T]) RemoveOldest() (key T, value any, ok bool) {
	e := c.list.Back()
	if e != nil {
		entry := e.Value.(*Entry[T])
		c.removeElement(e)
		return entry.Key, entry.Value, true
	}
	return
}

// Purge 清除所有
func (c *Cache[T]) Purge() {
	for k, v := range c.items {
		if c.onEvicted != nil {
			c.onEvicted(k, v)
		}
		delete(c.items, k)
	}
	c.list.Init()
}

// removeOldest removes the oldest item from the cache.
func (c *Cache[T]) removeOldest() {
	ent := c.list.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

// remove a given list element from the cache
func (c *Cache[T]) removeElement(e *list.Element) {
	entry := e.Value.(*Entry[T])
	c.list.Remove(e)
	delete(c.items, entry.Key)
	if c.onEvicted != nil {
		c.onEvicted(entry.Key, entry.Value)
	}
}
