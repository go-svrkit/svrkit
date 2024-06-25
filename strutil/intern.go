// Copyright (c) 2017-2018 Phil Pearl

package strutil

import (
	"sync"
)

// string interning
var interPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string)
	},
}

// InternStr returns s interned.
func InternStr(s string) string {
	m := interPool.Get().(map[string]string)
	c, ok := m[s]
	if ok {
		interPool.Put(m)
		return c
	}
	m[s] = s
	interPool.Put(m)
	return s
}

// InternBytes returns b converted to a string, interned.
func InternBytes(b []byte) string {
	m := interPool.Get().(map[string]string)
	c, ok := m[string(b)]
	if ok {
		interPool.Put(m)
		return c
	}
	s := string(b)
	m[s] = s
	interPool.Put(m)
	return s
}
