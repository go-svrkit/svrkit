// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package trie

const WildCardChar = '*' // 通配符

// WalkAction defines some action to take on the given key and value during
// a Trie Walk. Returning a non-nil error will terminate the Walk.
type WalkAction func(key string, value any) error

// Node trie map node
type Node struct {
	children map[rune]*Node
	value    any
}

func (n *Node) IsLeaf() bool {
	return len(n.children) == 0
}

func (n *Node) Walk(path string, action WalkAction) error {
	if n.value != nil {
		if err := action(path, n.value); err != nil {
			return err
		}
	}
	for r, child := range n.children {
		if err := child.Walk(path+string(r), action); err != nil {
			return err
		}
	}
	return nil
}

// Map rune trie map
type Map struct {
	root    Node //
	size    int  // count of nodes
	version int  //
}

func NewMap() *Map {
	return new(Map)
}

func (m *Map) Size() int {
	return m.size
}

// Contains returns true if map has the given key stored
func (m *Map) Contains(key string) bool {
	var node = m.getNode(key)
	return node != nil
}

// Get returns the value stored at the given key
func (m *Map) Get(key string) any {
	var node = m.getNode(key)
	if node != nil {
		return node.value
	}
	return nil
}

// ShortestPrefixOf find shortest prefix of `query`
// ['the','them'] ShortestPrefixOf("themxyz") -> "the"
func (m *Map) ShortestPrefixOf(query string) string {
	var node = &m.root
	for i, r := range query {
		node = node.children[r]
		if node == nil {
			return ""
		}
		if node.value != nil {
			return query[:i+1]
		}
	}
	if node != nil && node.value != nil {
		return query // its query itself
	}
	return ""
}

// LongestPrefixOf find longtest prefix of `query`
// ['the','them'] LongestPrefixOf("themxyz") -> "them"
func (m *Map) LongestPrefixOf(query string) string {
	var node = &m.root
	var maxLen = 0
	for i, r := range query {
		node = node.children[r]
		if node == nil {
			if maxLen > 0 {
				return query[:maxLen+1]
			}
			return ""
		}
		if node.value != nil {
			maxLen = i
		}
	}
	if node != nil && node.value != nil {
		return query // its query itself
	}
	if maxLen > 0 {
		return query[:maxLen+1]
	}
	return ""
}

// KeysWithPrefix find all words with prefix `prefix`
// keysWithPrefix("th") -> ["that", "the", "them"]
func (m *Map) KeysWithPrefix(prefix string) []string {
	var node = m.getNode(prefix)
	if node == nil {
		return nil
	}
	var result = make([]string, 0, len(node.children))
	node.Walk("", func(key string, value any) error {
		result = append(result, prefix+key)
		return nil
	})
	return result
}

func (m *Map) HasKeyWithPrefix(prefix string) bool {
	var node = m.getNode(prefix)
	return node != nil
}

// KeysWithPattern KeysWithPattern("t*a*") -> ["team", "that"]
func (m *Map) KeysWithPattern(pattern string) []string {
	var result = make([]string, 0)
	m.Walk(func(key string, value any) error {
		if len(key) < len(pattern) {
			return nil
		}
		for i := 0; i < len(pattern); i++ {
			if pattern[i] != WildCardChar && pattern[i] != key[i] {
				return nil
			}
		}
		result = append(result, key)
		return nil
	})
	return result
}

// Walk iterates over each key/value stored in the trie and calls the given
// walker function with the key and value. If the walker function returns
// an error, the walk is aborted.
func (m *Map) Walk(action WalkAction) error {
	return m.root.Walk("", action)
}

// Put inserts the value into the trie at the given key, replacing any existing items.
func (m *Map) Put(key string, v any) bool {
	var node = m.getNode(key)
	if node != nil {
		node.value = v
		return false
	}
	node = m.insertNode(key)
	node.value = v
	m.version++
	return true
}

// PutIfAbsent inserts the value into the trie at the given key only when key not exists.
func (m *Map) PutIfAbsent(key string, v any) bool {
	var node = m.getNode(key)
	if node != nil {
		return false
	}
	node = m.insertNode(key)
	node.value = v
	m.version++
	return true
}

// Delete removes the value associated with the given key.
func (m *Map) Delete(key string) bool {
	var deleted = m.removeNodes(key)
	m.size--
	m.version++
	return deleted
}

func (m *Map) getNode(key string) *Node {
	var node = &m.root
	for _, r := range key {
		node = node.children[r]
		if node == nil {
			return nil
		}
	}
	return node
}

func (m *Map) insertNode(key string) *Node {
	var node = &m.root
	for _, r := range key {
		child := node.children[r]
		if child == nil {
			if node.children == nil {
				node.children = make(map[rune]*Node)
			}
			child = new(Node)
			node.children[r] = child
			m.size++
		}
		node = child
	}
	return node
}

func (m *Map) removeNodes(key string) bool {
	var node = &m.root
	var path = make([]nodeRune, 0, len(key))
	for _, r := range key {
		path = append(path, nodeRune{node: node, r: r})
		node = node.children[r]
		if node == nil {
			return false // node not exist
		}
	}
	node.value = nil // mark deleted

	// remove leaf node from its parent. repeat for ancestor path
	if node.IsLeaf() {
		for i := len(path) - 1; i >= 0; i-- {
			var parent = path[i].node
			var r = path[i].r
			delete(parent.children, r)
			m.size--
			if !parent.IsLeaf() {
				break // parent has other children
			}
			parent.children = nil
			if parent.value != nil {
				break // parent has value
			}
		}
	}
	return true
}

// used for node deletion
type nodeRune struct {
	node *Node
	r    rune
}
