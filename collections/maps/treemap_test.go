// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package maps

import (
	"math"
	"math/rand"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func randText(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b = make([]byte, n)
	for i := range b {
		idx := rand.Int() % len(alphabet)
		b[i] = alphabet[idx]
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func TestTreeMap_Put(t *testing.T) {
	// test insert empty tree
	var m = NewOrderedTreeMap[int, string]()
	m.Put(1, "a")
	v, ok := m.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "a", v)

	// test insert non-empty tree
	m.Put(2, "b")
	v, ok = m.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "b", v)

	// test replace tree element
	m.Put(2, "c")
	v, ok = m.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "c", v)

	// test insert order
	for i := 1; i < 100; i++ {
		m.Put(i, randText(8))
	}
	var maxHeight = math.Ceil(2 * math.Log(float64(m.size+1))) // 高度不应超过2log(n+1)
	var height = m.root.Height()
	assert.Equal(t, int(maxHeight), height)
	var node = m.getFirstEntry()
	assert.NotNil(t, node)
	for i := 1; i < 100 && node != nil; i++ {
		assert.Equal(t, i, node.key)
		node = successor(node)
	}
}
