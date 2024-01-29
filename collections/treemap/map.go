// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package treemap

import (
	"gopkg.in/svrkit.v1/collections/util"
)

// Map is a Red-Black tree based map implementation.
// The map is sorted according to the Comparable natural ordering of its keys
// This implementation provides guaranteed log(n) time cost for the
// ContainsKey(), Get(), Put() and Remove() operations.
//
// more details see links below
// https://github.com/openjdk/jdk/blob/jdk-17+35/src/java.base/share/classes/java/util/TreeMap.java
type Map[K, V any] struct {
	root       *Entry[K, V]
	comparator util.Comparator[K]
	size       int // The number of entries in the tree
	version    int // The number of structural modifications to the tree.
}

func New[K, V any](comparator util.Comparator[K]) *Map[K, V] {
	return &Map[K, V]{
		comparator: comparator,
	}
}

// Size returns the number of key-value mappings in this map.
func (m *Map[K, V]) Size() int {
	return m.size
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.size == 0
}

// ContainsKey return true if this map contains a mapping for the specified key
func (m *Map[K, V]) ContainsKey(key K) bool {
	return m.getEntry(key) != nil
}

// Get returns the value to which the specified key is mapped,
// or nil if this map contains no mapping for the key.
func (m *Map[K, V]) Get(key K) (V, bool) {
	var p = m.getEntry(key)
	if p != nil {
		return p.value, true
	}
	return zeroOf[V](), false
}

// GetOrDefault returns the value to which the specified key is mapped,
// or `defaultValue` if this map contains no mapping for the key.
func (m *Map[K, V]) GetOrDefault(key K, defVal V) V {
	var p = m.getEntry(key)
	if p != nil {
		return p.value
	}
	return defVal
}

// FirstKey returns the first key in the TreeMap (according to the key's order)
func (m *Map[K, V]) FirstKey() (K, bool) {
	return keyOf[K](m.getFirstEntry())
}

// LastKey returns the last key in the TreeMap (according to the key's order)
func (m *Map[K, V]) LastKey() (K, bool) {
	return keyOf(m.getLastEntry())
}

func (m *Map[K, V]) RootEntry() *Entry[K, V] {
	return m.root
}

func (m *Map[K, V]) FirstEntry() *Entry[K, V] {
	return m.getFirstEntry()
}

func (m *Map[K, V]) LastEntry() *Entry[K, V] {
	return m.getLastEntry()
}

// FloorEntry gets the entry corresponding to the specified key;
// if no such entry exists, returns the entry for the greatest key less than the specified key;
func (m *Map[K, V]) FloorEntry(key K) *Entry[K, V] {
	return m.getFloorEntry(key)
}

// FloorKey gets the specified key, returns the greatest key less than the specified key if not exist.
func (m *Map[K, V]) FloorKey(key K) (K, bool) {
	var entry = m.getFloorEntry(key)
	if entry != nil {
		return entry.key, true
	}
	return zeroOf[K](), false
}

// CeilingEntry gets the entry corresponding to the specified key;
// returns the entry for the least key greater than the specified key if not exist.
func (m *Map[K, V]) CeilingEntry(key K) *Entry[K, V] {
	return m.getCeilingEntry(key)
}

// CeilingKey gets the specified key, return the least key greater than the specified key if not exist.
func (m *Map[K, V]) CeilingKey(key K) (K, bool) {
	var entry = m.getCeilingEntry(key)
	if entry != nil {
		return entry.key, true
	}
	return zeroOf[K](), false
}

// HigherEntry gets the entry for the least key greater than the specified key
func (m *Map[K, V]) HigherEntry(key K) *Entry[K, V] {
	return m.getHigherEntry(key)
}

// HigherKey returns the least key greater than the specified key
func (m *Map[K, V]) HigherKey(key K) (K, bool) {
	var entry = m.getHigherEntry(key)
	if entry != nil {
		return entry.key, true
	}
	return zeroOf[K](), false
}

// Foreach performs the given action for each entry in this map until all entries
// have been processed or the action panic
func (m *Map[K, V]) Foreach(visit util.KeyValVisitor[K, V]) {
	var ver = m.version
	for e := m.getFirstEntry(); e != nil; e = successor(e) {
		visit(e.key, e.value)
		if ver != m.version {
			panic("concurrent map modification")
		}
	}
}

// Keys return list of all keys
func (m *Map[K, V]) Keys() []K {
	var keys = make([]K, 0, m.size)
	for e := m.getFirstEntry(); e != nil; e = successor(e) {
		keys = append(keys, e.key)
	}
	return keys
}

// Values return list of all values
func (m *Map[K, V]) Values() []V {
	var values = make([]V, 0, m.size)
	for e := m.getFirstEntry(); e != nil; e = successor(e) {
		values = append(values, e.value)
	}
	return values
}

func (m *Map[K, V]) Iterator() util.Iterator[*Entry[K, V]] {
	return NewEntryIterator(m, m.getFirstEntry())
}

func (m *Map[K, V]) DescendingIterator() util.Iterator[*Entry[K, V]] {
	return NewKeyDescendingEntryIterator(m, m.getLastEntry())
}

func (m *Map[K, V]) KeyIterator() util.Iterator[K] {
	return NewKeyIterator(m, m.getFirstEntry())
}

func (m *Map[K, V]) DescendingKeyIterator() util.Iterator[K] {
	return NewDescendingKeyIterator(m, m.getLastEntry())
}

func (m *Map[K, V]) ValueIterator() util.Iterator[V] {
	return NewValueIterator(m, m.getFirstEntry())
}

// Put associates the specified value with the specified key in this map.
// If the map previously contained a mapping for the key, the old value is replaced.
func (m *Map[K, V]) Put(key K, value V) V {
	return m.put(key, value, true)
}

func (m *Map[K, V]) PutIfAbsent(key K, value V) V {
	return m.put(key, value, false)
}

// Remove removes the mapping for this key from this TreeMap if present.
func (m *Map[K, V]) Remove(key K) bool {
	var p = m.getEntry(key)
	if p != nil {
		m.deleteEntry(p)
		return true
	}
	return false
}

// Clear removes all the mappings from this map.
func (m *Map[K, V]) Clear() {
	m.version++
	m.size = 0
	m.root = nil
}

// Returns the first Entry in the TreeMap (according to the key's order)
// Returns nil if the TreeMap is empty.
func (m *Map[K, V]) getFirstEntry() *Entry[K, V] {
	var p = m.root
	if p != nil {
		for p.left != nil {
			p = p.left
		}
	}
	return p
}

// Returns the last Entry in the TreeMap (according to the key's order)
// Returns nil if the TreeMap is empty.
func (m *Map[K, V]) getLastEntry() *Entry[K, V] {
	var p = m.root
	if p != nil {
		for p.right != nil {
			p = p.right
		}
	}
	return p
}

// Returns this map's entry for the given key,
// or nil if the map does not contain an entry for the key.
func (m *Map[K, V]) getEntry(key K) *Entry[K, V] {
	var p = m.root
	for p != nil {
		var cmp = m.comparator(key, p.key)
		if cmp < 0 {
			p = p.left
		} else if cmp > 0 {
			p = p.right
		} else {
			return p
		}
	}
	return nil
}

// Gets the entry corresponding to the specified key;
// if no such entry exists, returns the entry for the least key greater than the specified key;
// if no such entry exists returns nil.
func (m *Map[K, V]) getCeilingEntry(key K) *Entry[K, V] {
	var p = m.root
	for p != nil {
		var cmp = m.comparator(key, p.key)
		if cmp < 0 {
			if p.left != nil {
				p = p.left
			} else {
				return p
			}
		} else if cmp > 0 {
			if p.right != nil {
				p = p.right
			} else {
				var parent = p.parent
				var ch = p
				for parent != nil && ch == parent.right {
					ch = parent
					parent = parent.parent
				}
				return parent
			}
		} else {
			return p
		}
	}
	return nil
}

// Gets the entry corresponding to the specified key;
// if no such entry exists, returns the entry for the greatest key less than the specified key;
// if no such entry exists, returns nil.
func (m *Map[K, V]) getFloorEntry(key K) *Entry[K, V] {
	var p = m.root
	for p != nil {
		var cmp = m.comparator(key, p.key)
		if cmp > 0 {
			if p.right != nil {
				p = p.right
			} else {
				return p
			}
		} else if cmp < 0 {
			if p.left != nil {
				p = p.left
			} else {
				var parent = p.parent
				var ch = p
				for parent != nil && ch == parent.left {
					ch = parent
					parent = parent.parent
				}
				return parent
			}
		} else {
			return p
		}

	}
	return nil
}

// Gets the entry for the least key greater than the specified key;
// if no such entry exists, returns the entry for the least key greater than the specified key;
// if no such entry exists returns nil.
func (m *Map[K, V]) getHigherEntry(key K) *Entry[K, V] {
	var p = m.root
	for p != nil {
		var cmp = m.comparator(key, p.key)
		if cmp < 0 {
			if p.left != nil {
				p = p.left
			} else {
				return p
			}
		} else {
			if p.right != nil {
				p = p.right
			} else {
				var parent = p.parent
				var ch = p
				for parent != nil && ch == parent.right {
					ch = parent
					parent = parent.parent
				}
				return parent
			}
		}
	}
	return nil
}

// Returns the entry for the greatest key less than the specified key;
// if no such entry exists (i.e., the least key in the Tree is greater than the specified key), returns nil
func (m *Map[K, V]) getLowerEntry(key K) *Entry[K, V] {
	var p = m.root
	for p != nil {
		var cmp = m.comparator(key, p.key)
		if cmp > 0 {
			if p.right != nil {
				p = p.right
			} else {
				return p
			}
		} else {
			if p.left != nil {
				p = p.left
			} else {
				var parent = p.parent
				var ch = p
				for parent != nil && ch == parent.left {
					ch = parent
					parent = parent.parent
				}
				return parent
			}
		}
	}
	return nil
}

func (m *Map[K, V]) put(key K, value V, replaceOld bool) V {
	var zero = zeroOf[V]()
	var t = m.root
	if t == nil {
		m.addEntryToEmptyMap(key, value)
		return zero
	}
	var cmp int
	var parent *Entry[K, V]
	for {
		parent = t
		cmp = m.comparator(key, t.key)
		if cmp < 0 {
			t = t.left
		} else if cmp > 0 {
			t = t.right
		} else {
			var oldValue = t.value
			if replaceOld {
				t.value = value
			}
			return oldValue
		}
		if t == nil {
			break
		}
	}
	m.addEntry(key, value, parent, cmp < 0)
	return zero
}

func (m *Map[K, V]) addEntry(key K, value V, parent *Entry[K, V], addToLeft bool) {
	var entry = NewEntry(key, value, parent)
	if addToLeft {
		parent.left = entry
	} else {
		parent.right = entry
	}
	m.fixAfterInsertion(entry)
	m.size++
	m.version++
}

func (m *Map[K, V]) addEntryToEmptyMap(key K, value V) {
	m.root = NewEntry(key, value, nil)
	m.size = 1
	m.version++
}

func (m *Map[K, V]) deleteEntry(p *Entry[K, V]) {
	m.version++
	m.size--

	// If strictly internal, copy successor's element to p and then make p
	// point to successor.
	if p.left != nil && p.right != nil {
		var s = successor(p)
		p.key = s.key
		p.value = s.value
		p = s
	} // p has 2 children

	// Start fixup at replacement node, if it exists.
	var replacement = p.left
	if p.left == nil {
		replacement = p.right
	}

	if replacement != nil {
		// Link replacement to parent
		replacement.parent = p.parent
		if p.parent == nil {
			m.root = replacement
		} else if p == p.parent.left {
			p.parent.left = replacement
		} else {
			p.parent.right = replacement
		}

		// Null out links, so they are OK to use by fixAfterDeletion.
		p.left = nil
		p.right = nil
		p.parent = nil

		// Fix replacement
		if p.color == BLACK {
			m.fixAfterDeletion(replacement)
		}
	} else if p.parent == nil { // return if we are the only node.
		m.root = nil
	} else { //  No children. Use self as phantom replacement and unlink.
		if p.color == BLACK {
			m.fixAfterDeletion(p)
		}
		if p.parent != nil {
			if p == p.parent.left {
				p.parent.left = nil
			} else if p == p.parent.right {
				p.parent.right = nil
			}
			p.parent = nil
		}
	}
}

// Balancing operations.
//
// Implementations of rebalancings during insertion and deletion are
// slightly different from the Cormen, Leiserson, and Rivest's <Introduction to Algorithms> version.
// Rather than using dummy nil nodes, we use a set of accessors that deal properly with nil.
// They are used to avoid messiness surrounding nullness checks in the main algorithms.
//
// see original version at http://staff.ustc.edu.cn/~csli/graduate/algorithms/book6/chap14.htm

func (m *Map[K, V]) rotateLeft(p *Entry[K, V]) {
	if p != nil {
		var r = p.right
		p.right = r.left
		if r.left != nil {
			r.left.parent = p
		}
		r.parent = p.parent
		if p.parent == nil {
			m.root = r
		} else if p.parent.left == p {
			p.parent.left = r
		} else {
			p.parent.right = r
		}
		r.left = p
		p.parent = r
	}
}

func (m *Map[K, V]) rotateRight(p *Entry[K, V]) {
	if p != nil {
		var l = p.left
		p.left = l.right
		if l.right != nil {
			l.right.parent = p
		}
		l.parent = p.parent
		if p.parent == nil {
			m.root = l
		} else if p.parent.right == p {
			p.parent.right = l
		} else {
			p.parent.left = l
		}
		l.right = p
		p.parent = l
	}
}

func (m *Map[K, V]) fixAfterInsertion(x *Entry[K, V]) {
	x.color = RED
	for x != nil && x != m.root && x.parent.color == RED {
		if parentOf(x) == leftOf(parentOf(parentOf(x))) {
			var y = rightOf(parentOf(parentOf(x)))
			if colorOf(y) == RED {
				setColor(parentOf(x), BLACK)
				setColor(y, BLACK)
				setColor(parentOf(parentOf(x)), RED)
				x = parentOf(parentOf(x))
			} else {
				if x == rightOf(parentOf(x)) {
					x = parentOf(x)
					m.rotateLeft(x)
				}
				setColor(parentOf(x), BLACK)
				setColor(parentOf(parentOf(x)), RED)
				m.rotateRight(parentOf(parentOf(x)))
			}
		} else {
			var y = leftOf(parentOf(parentOf(x)))
			if colorOf(y) == RED {
				setColor(parentOf(x), BLACK)
				setColor(y, BLACK)
				setColor(parentOf(parentOf(x)), RED)
				x = parentOf(parentOf(x))
			} else {
				if x == leftOf(parentOf(x)) {
					x = parentOf(x)
					m.rotateRight(x)
				}
				setColor(parentOf(x), BLACK)
				setColor(parentOf(parentOf(x)), RED)
				m.rotateLeft(parentOf(parentOf(x)))
			}
		}
	}
	m.root.color = BLACK
}

func (m *Map[K, V]) fixAfterDeletion(x *Entry[K, V]) {
	for x != m.root && colorOf(x) == BLACK {
		if x == leftOf(parentOf(x)) {
			var sib = rightOf(parentOf(x))

			if colorOf(sib) == RED {
				setColor(sib, BLACK)
				setColor(parentOf(x), RED)
				m.rotateLeft(parentOf(x))
				sib = rightOf(parentOf(x))
			}

			if colorOf(leftOf(sib)) == BLACK &&
				colorOf(rightOf(sib)) == BLACK {
				setColor(sib, RED)
				x = parentOf(x)
			} else {
				if colorOf(rightOf(sib)) == BLACK {
					setColor(leftOf(sib), BLACK)
					setColor(sib, RED)
					m.rotateRight(sib)
					sib = rightOf(parentOf(x))
				}
				setColor(sib, colorOf(parentOf(x)))
				setColor(parentOf(x), BLACK)
				setColor(rightOf(sib), BLACK)
				m.rotateLeft(parentOf(x))
				x = m.root
			}
		} else { // symmetric
			var sib = leftOf(parentOf(x))

			if colorOf(sib) == RED {
				setColor(sib, BLACK)
				setColor(parentOf(x), RED)
				m.rotateRight(parentOf(x))
				sib = leftOf(parentOf(x))
			}

			if colorOf(rightOf(sib)) == BLACK &&
				colorOf(leftOf(sib)) == BLACK {
				setColor(sib, RED)
				x = parentOf(x)
			} else {
				if colorOf(leftOf(sib)) == BLACK {
					setColor(rightOf(sib), BLACK)
					setColor(sib, RED)
					m.rotateLeft(sib)
					sib = leftOf(parentOf(x))
				}
				setColor(sib, colorOf(parentOf(x)))
				setColor(parentOf(x), BLACK)
				setColor(leftOf(sib), BLACK)
				m.rotateRight(parentOf(x))
				x = m.root
			}
		}
	}
	setColor(x, BLACK)
}
