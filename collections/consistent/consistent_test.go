// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package consistent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testAddOneNode(t *testing.T, node int32) *Consistent {
	var c = New()
	c.Add(node)
	assert.Equal(t, 1, c.Len())
	assert.True(t, c.nodes[node])
	assert.Equal(t, len(c.sortedHash), ReplicaCount)
	for i := 0; i < 10; i++ {
		var n = c.GetBy(int64(i))
		assert.Equal(t, n, node) // only one node, all should be equal
	}
	return c
}

func testAddNNode(t *testing.T, N int) *Consistent {
	var c = New()
	for i := 0; i < N; i++ {
		var node = 100 + int32(i)
		c.Add(node)
		assert.Equal(t, i+1, c.Len())
		assert.True(t, c.nodes[node])
		assert.Equal(t, len(c.sortedHash), (i+1)*ReplicaCount)
	}
	for i := 0; i < 10; i++ {
		var n = c.GetBy(int64(i))
		assert.NotEqual(t, n, int32(0)) // only one node, all should be equal
	}
	return c
}

func TestConsistent_Add(t *testing.T) {
	testAddOneNode(t, 1234)
	testAddNNode(t, 10)
}

func TestConsistent_Remove(t *testing.T) {
	var node int32 = 1234
	var c = testAddOneNode(t, node)
	c.Remove(node)
	assert.Equal(t, 0, c.Len())
	assert.False(t, c.nodes[node])
	assert.Equal(t, len(c.sortedHash), 0)
	var n = c.Get("anykey")
	assert.Equal(t, int32(0), n)
}

func TestConsistent_Get(t *testing.T) {
	var c = testAddNNode(t, 10)
	for i := 0; i < 9; i++ {
		var n = 100 + int32(i)
		c.Remove(n)
	}
	assert.Equal(t, 1, c.Len())
	var n = c.Get("anykey")
	assert.NotEqual(t, int32(0), n)
}

func BenchmarkConsistent_Add(b *testing.B) {
	var c = New()
	for i := 1; i <= 500; i++ {
		c.Add(int32(i))
	}
}

func BenchmarkConsistent_AddMany(b *testing.B) {
	var c = New()
	var nodes = make([]int32, 0, 500)
	for i := 1; i <= 500; i++ {
		nodes = append(nodes, int32(i))
	}
	c.Add(nodes...)
}

func BenchmarkConsistent_Remove(b *testing.B) {
	var c = New()
	b.StopTimer()
	for i := 1; i <= 500; i++ {
		c.Add(int32(i))
	}
	b.StartTimer()
	for i := 1; i <= 500; i++ {
		c.Remove(int32(i))
	}
}

func BenchmarkConsistent_RemoveMany(b *testing.B) {
	var c = New()
	b.StopTimer()
	var nodes = make([]int32, 0, 500)
	for i := 1; i <= 500; i++ {
		nodes = append(nodes, int32(i))
	}
	c.Add(nodes...)
	b.StartTimer()
	for i := 1; i <= 500; i++ {
		c.Remove(int32(i))
	}
}

func BenchmarkConsistent_Get(b *testing.B) {
	b.StopTimer()
	var c = New()
	for i := 1; i <= 100; i++ {
		c.Add(int32(i))
	}
	b.StartTimer()
	for i := 0; i < 1000; i++ {
		c.GetBy(12345)
	}
}
