// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package trie

import (
	"testing"
)

func TestMap_PutDelete(t *testing.T) {
	const firstPutValue = "first put"
	cases := []struct {
		key   string
		value interface{}
	}{
		{"fish", 0},
		{"/cat", 1},
		{"/dog", 2},
		{"/cats", 3},
		{"/caterpillar", 4},
		{"/cat/gideon", 5},
		{"/cat/giddy", 6},
	}

	var trie = NewMap()
	// get missing keys
	for _, c := range cases {
		if value := trie.Get(c.key); value != nil {
			t.Errorf("expected key %s to be missing, found value %v", c.key, value)
		}
	}

	// initial put
	for _, c := range cases {
		if isNew := trie.Put(c.key, firstPutValue); !isNew {
			t.Errorf("expected key %s to be missing", c.key)
		}
	}

	// subsequent put
	for _, c := range cases {
		if isNew := trie.Put(c.key, c.value); isNew {
			t.Errorf("expected key %s to have a value already", c.key)
		}
	}

	// get
	for _, c := range cases {
		if value := trie.Get(c.key); value != c.value {
			t.Errorf("expected key %s to have value %v, got %v", c.key, c.value, value)
		}
	}

	// delete, expect Delete to return true indicating a node was nil'd
	for _, c := range cases {
		if deleted := trie.Delete(c.key); !deleted {
			t.Errorf("expected key %s to be deleted", c.key)
		}
	}

	// delete cleaned all the way to the first character
	// expect Delete to return false bc no node existed to nil
	for _, c := range cases {
		if deleted := trie.Delete(string(c.key[0])); deleted {
			t.Errorf("expected key %s to be cleaned by delete", string(c.key[0]))
		}
	}

	// get deleted keys
	for _, c := range cases {
		if value := trie.Get(c.key); value != nil {
			t.Errorf("expected key %s to be deleted, got value %v", c.key, value)
		}
	}
}

func TestMap_ShortestPrefixOf(t *testing.T) {
	var trie = NewMap()
	trie.Put("the", true)
	trie.Put("them", true)
	var r = trie.ShortestPrefixOf("themxyz")
	if r != "the" {
		t.Fatalf("unexpected result %v", r)
	}
}

func TestMap_LongestPrefixOf(t *testing.T) {
	var trie = NewMap()
	trie.Put("the", true)
	trie.Put("them", true)
	var r = trie.LongestPrefixOf("themxyz")
	if r != "them" {
		t.Fatalf("unexpected result %v", r)
	}
}

func TestMap_KeyWithPrefix(t *testing.T) {
	var trie = NewMap()
	trie.Put("the", true)
	trie.Put("them", true)
	trie.Put("this", true)
	trie.Put("that", true)
	var keys = trie.KeysWithPrefix("th")
	t.Logf("keys with prefix `th`: %v", keys)
	if len(keys) != 4 {
		t.Fatalf("unexpected result %v", len(keys))
	}
}

func TestMap_KeysWithPattern(t *testing.T) {
	var trie = NewMap()
	trie.Put("the", true)
	trie.Put("them", true)
	trie.Put("this", true)
	trie.Put("that", true)
	var keys = trie.KeysWithPattern("t*e")
	t.Logf("keys with pattern `t*e*`: %v", keys)
	if len(keys) != 2 {
		t.Fatalf("unexpected result %v", len(keys))
	}
}
