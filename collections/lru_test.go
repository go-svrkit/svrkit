// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package collections

import (
	"strconv"
	"testing"
)

func TestLRU(t *testing.T) {
	c := NewLRUCache[string, int](128, nil)
	for i := 0; i < 256; i++ {
		key := strconv.Itoa(i)
		c.Put(key, i)
	}
	if c.Len() != 128 {
		t.Fatalf("expect len %d, got: %d", 128, c.Len())
	}
	keys := c.Keys()
	for i, k := range keys {
		if v, ok := c.Get(k); !ok {
			t.Fatalf("bad Key: %v", k)
		} else {
			key := strconv.Itoa(v)
			if key != k || v != i+128 {
				t.Fatalf("bad Key: %v", k)
			}
		}
	}
	for i := 0; i < 128; i++ {
		key := strconv.Itoa(i)
		_, ok := c.Get(key)
		if ok {
			t.Fatalf("should be evicted: %v", key)
		}
	}
	for i := 128; i < 256; i++ {
		key := strconv.Itoa(i)
		_, found := c.Get(key)
		if !found {
			t.Fatalf("should be evicted: %v", key)
		}
	}
	for i := 128; i < 192; i++ {
		key := strconv.Itoa(i)
		ok := c.Remove(key)
		if !ok {
			t.Fatalf("should be contained: %v", key)
		}
		ok = c.Remove(key)
		if ok {
			t.Fatalf("should not be contained: %v", key)
		}
		_, ok = c.Get(key)
		if ok {
			t.Fatalf("should be deleted: %v", key)
		}
	}

	c.Get("192") // expect 192 to be last Key in l.Keys()
	keys = c.Keys()
	for i, k := range keys {
		vk := strconv.Itoa(i + 193)
		if (i < 63 && k != vk) || (i == 63 && k != "192") {
			t.Fatalf("out of order Key: %v", k)
		}
	}

	c.Purge()
	if c.Len() != 0 {
		t.Fatalf("bad len: %v", c.Len())
	}
	if _, ok := c.Get("200"); ok {
		t.Fatalf("should contain nothing")
	}
}

func TestLRUOldest(t *testing.T) {
	c := NewLRUCache[string, int](128, nil)
	for i := 0; i < 256; i++ {
		k := strconv.Itoa(i)
		c.Put(k, i)
	}
	k, _, ok := c.GetOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != "128" {
		t.Fatalf("bad: %v", k)
	}

	k, _, ok = c.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != "128" {
		t.Fatalf("bad: %v", k)
	}

	k, _, ok = c.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != "129" {
		t.Fatalf("bad: %v", k)
	}
}

func TestLRU_Put(t *testing.T) {
	c := NewLRUCache[string, int](1, nil)

	c.Put("1", 1)

	if c.Len() != 1 {
		t.Errorf("bad len: %v", c.Len())
	}
	if !c.Contains("1") {
		t.Errorf("should exist Key 1")
	}
	if v, found := c.Get("1"); !found || v != 1 {
		t.Errorf("bad Value: %v", v)
	}

	c.Put("2", 2)

	if c.Len() != 1 {
		t.Errorf("bad len: %v", c.Len())
	}
	if c.Contains("1") {
		t.Errorf("should not exist Key 1")
	}
	if !c.Contains("2") {
		t.Errorf("should exist Key 2")
	}
}

func TestLRU_Exist(t *testing.T) {
	c := NewLRUCache[string, int](2, nil)

	c.Put("1", 1)
	c.Put("2", 2)
	if !c.Contains("1") {
		t.Errorf("1 should be contained")
	}

	c.Put("3", 3)
	if c.Contains("1") {
		t.Errorf("Contains should not have updated recent-ness of 1")
	}
}

// Test that Peek doesn't update recent-ness
func TestLRU_Peek(t *testing.T) {
	c := NewLRUCache[string, int](2, nil)

	c.Put("1", 1)
	c.Put("2", 2)
	if v, ok := c.Peek("1"); !ok || v != 1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	c.Put("3", 3)
	if c.Contains("1") {
		t.Errorf("should not have updated recent-ness of 1")
	}
}
