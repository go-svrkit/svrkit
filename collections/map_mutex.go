// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import (
	"maps"
	"sync"
)

// MapInterface is the interface SyncMap implements.
type MapInterface[K comparable, V any] interface {
	Size() int
	Contains(K) bool
	Load(K) (V, bool)
	Store(key K, value V)
	LoadOrStore(key K, value V) (actual V, loaded bool)
	LoadAndDelete(key K) (value V, loaded bool)
	Delete(K)
	Range(func(key K, value V) (shouldContinue bool))
	Keys() []K
	Values() []V
	CloneMap() map[K]V
}

var _ MapInterface[int, int] = (*MutexMap[int, int])(nil)

// MutexMap is an implementation of mapInterface using a sync.RWMutex.
type MutexMap[K comparable, V any] struct {
	mu    sync.RWMutex
	dirty map[K]V
}

// Size returns the number of Key-Value mappings in this map.
func (m *MutexMap[K, V]) Size() int {
	m.mu.RLock()
	size := len(m.dirty)
	m.mu.RUnlock()
	return size
}

func (m *MutexMap[K, V]) Contains(key K) bool {
	m.mu.RLock()
	_, ok := m.dirty[key]
	m.mu.RUnlock()
	return ok
}

func (m *MutexMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

func (m *MutexMap[K, V]) Keys() []K {
	m.mu.RLock()
	var keys = make([]K, 0, len(m.dirty))
	for k := range m.dirty {
		keys = append(keys, k)
	}
	m.mu.RUnlock()
	return keys
}

func (m *MutexMap[K, V]) Values() []V {
	m.mu.RLock()
	var values = make([]V, 0, len(m.dirty))
	for _, val := range m.dirty {
		values = append(values, val)
	}
	m.mu.RUnlock()
	return values
}

func (m *MutexMap[K, V]) CloneMap() map[K]V {
	m.mu.RLock()
	var clone = maps.Clone(m.dirty)
	m.mu.RUnlock()
	return clone
}

func (m *MutexMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.RLock()
	value, ok = m.dirty[key]
	m.mu.RUnlock()
	return
}

func (m *MutexMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[K]V)
	}
	m.dirty[key] = value
	m.mu.Unlock()
}

func (m *MutexMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	actual, loaded = m.dirty[key]
	if !loaded {
		actual = value
		if m.dirty == nil {
			m.dirty = make(map[K]V)
		}
		m.dirty[key] = value
	}
	m.mu.Unlock()
	return actual, loaded
}

func (m *MutexMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.mu.Lock()
	value, loaded = m.dirty[key]
	if !loaded {
		m.mu.Unlock()
		var zero V
		return zero, false
	}
	delete(m.dirty, key)
	m.mu.Unlock()
	return value, loaded
}

func (m *MutexMap[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.dirty, key)
	m.mu.Unlock()
}

func (m *MutexMap[K, V]) Range(f func(key K, value V) (shouldContinue bool)) {
	m.mu.RLock()
	keys := make([]K, 0, len(m.dirty))
	for k := range m.dirty {
		keys = append(keys, k)
	}
	m.mu.RUnlock()

	for _, k := range keys {
		v, ok := m.Load(k)
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}
