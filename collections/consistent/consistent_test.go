// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package consistent

import (
	"testing"
)

func TestConsistentExample(t *testing.T) {
	var c = New()
	t.Logf("add node 1 2 3 4")
	c.AddNode(0x10001)
	c.AddNode(0x10002)
	c.AddNode(0x10003)
	c.AddNode(0x10004)
	var key int64 = 1234
	node := c.GetNode(key)
	t.Logf("get node %v by %d", node, key)
	c.RemoveNode(0x10004)
	c.RemoveNode(0x10003)
	t.Logf("remove node 1 2")
	node = c.GetNode(key)
	t.Logf("get node %v by %d", node, key)
}

func BenchmarkConsistentAdd(b *testing.B) {
	var c = New()
	for i := 1; i <= 200; i++ {
		c.AddNode(int32(i))
	}
}

func BenchmarkConsistentGet(b *testing.B) {
	b.StopTimer()
	var c = New()
	for i := 1; i <= 200; i++ {
		c.AddNode(int32(i))
	}
	b.StartTimer()
	for i := 0; i < 1000; i++ {
		c.GetNode(12345)
	}
}
