// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package cluster

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	var node = NewNode("GATE", 1)
	t.Logf("%v", node)
	node.ID = 2
	t.Logf("%v", node)

	node.Set("hello", "world")
	var s = node.Get("hello")
	if s != "world" {
		t.Fatalf("get hello: %s", s)
	}

	node.SetInt("hello2", 123)
	var n = node.GetInt("hello2")
	if n != 123 {
		t.Fatalf("get hello2: %d", n)
	}

	node.SetBool("hello3", true)
	var b = node.GetBool("hello3")
	if !b {
		t.Fatalf("get hello3: %v", b)
	}
}

func TestNodeMap(t *testing.T) {
	var nm = NewNodeMap()
	nm.InsertNode(NewNode("GATE", 1))
	nm.InsertNode(NewNode("GATE", 2))
	nm.InsertNode(NewNode("GAME", 1))
	nm.InsertNode(NewNode("GAME", 2))
	nm.InsertNode(NewNode("LOGIN", 1))
	nm.InsertNode(NewNode("LOGIN", 2))
	t.Logf("initial nodes: %v", nm.String())

	var r = nm.GetNodes("GATE")
	t.Logf("get nodes: %v", r)

	nm.DeleteNode("LOGIN", 1)
	t.Logf("after del 1: %v", nm.String())

}
