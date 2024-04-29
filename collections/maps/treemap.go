// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package maps

import (
	"cmp"

	"gopkg.in/svrkit.v1/collections/cutil"
)

// TreeMap is a Red-Black tree based map implementation.
// The map is sorted by a comparator passed to its constructor.
// This implementation provides guaranteed log(n) time cost for the
// ContainsKey(), Get(), Put() and Remove() operations.
//
// detail code taken from JDK, see links below
// https://github.com/openjdk/jdk/blob/jdk-11+28/src/java.base/share/classes/java/util/TreeMap.java

// Color Red-black mechanics
type Color uint8

const (
	RED   Color = 0
	BLACK Color = 1
)

func (c Color) String() string {
	switch c {
	case RED:
		return "red"
	case BLACK:
		return "black"
	default:
		return "??"
	}
}

type TreeEntry[K, V comparable] struct {
	left, right, parent *TreeEntry[K, V]
	color               Color
	key                 K
	value               V
}

func NewTreeEntry[K, V comparable](key K, val V, parent *TreeEntry[K, V]) *TreeEntry[K, V] {
	return &TreeEntry[K, V]{
		key:    key,
		value:  val,
		parent: parent,
		color:  BLACK,
	}
}

func (e *TreeEntry[K, V]) GetKey() K {
	return e.key
}

func (e *TreeEntry[K, V]) GetValue() V {
	return e.value
}

func (e *TreeEntry[K, V]) SetValue(val V) V {
	var old = e.value
	e.value = val
	return old
}

func (e *TreeEntry[K, V]) Height() int {
	if e == nil {
		return 0
	}
	var leftHeight = e.left.Height()
	var rightHeight = e.right.Height()
	if leftHeight > rightHeight {
		return leftHeight + 1
	} else {
		return rightHeight + 1
	}
}

type TreeMap[K, V comparable] struct {
	root       *TreeEntry[K, V]
	comparator cutil.Comparator[K]
	size       int // The number of entries in the tree
	version    int // The number of structural modifications to the tree.
}

var _ MapInterface[int, int] = (*TreeMap[int, int])(nil)

func NewTreeMap[K, V comparable](comparator cutil.Comparator[K]) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		comparator: comparator,
	}
}

func NewOrderedTreeMap[K cmp.Ordered, V comparable]() *TreeMap[K, V] {
	return &TreeMap[K, V]{
		comparator: cmp.Compare[K],
	}
}

func NewTreeMapOf[M ~map[K]V, K cmp.Ordered, V comparable](unordered M) *TreeMap[K, V] {
	var m = &TreeMap[K, V]{
		comparator: cmp.Compare[K],
	}
	for k, v := range unordered {
		m.Put(k, v)
	}
	return m
}

// Size returns the number of key-value mappings in this map.
func (m *TreeMap[K, V]) Size() int {
	return m.size
}

func (m *TreeMap[K, V]) IsEmpty() bool {
	return m.size == 0
}

// ContainsKey return true if this map contains a mapping for the specified key
func (m *TreeMap[K, V]) ContainsKey(key K) bool {
	return m.getEntry(key) != nil
}

// Get returns the value to which the specified key is mapped,
// or nil if this map contains no mapping for the key.
func (m *TreeMap[K, V]) Get(key K) (V, bool) {
	var p = m.getEntry(key)
	if p != nil {
		return p.value, true
	}
	var zero V
	return zero, false
}

func (m *TreeMap[K, V]) Load(key K) (V, bool) {
	return m.Get(key)
}

// GetOrDefault returns the value to which the specified key is mapped,
// or `defaultValue` if this map contains no mapping for the key.
func (m *TreeMap[K, V]) GetOrDefault(key K, defVal V) V {
	var p = m.getEntry(key)
	if p != nil {
		return p.value
	}
	return defVal
}

// FirstKey returns the first key in the TreeMap (according to the key's order)
func (m *TreeMap[K, V]) FirstKey() (K, bool) {
	return keyOf[K](m.getFirstEntry())
}

// LastKey returns the last key in the TreeMap (according to the key's order)
func (m *TreeMap[K, V]) LastKey() (K, bool) {
	return keyOf(m.getLastEntry())
}

func (m *TreeMap[K, V]) RootEntry() *TreeEntry[K, V] {
	return m.root
}

func (m *TreeMap[K, V]) FirstEntry() *TreeEntry[K, V] {
	return m.getFirstEntry()
}

func (m *TreeMap[K, V]) LastEntry() *TreeEntry[K, V] {
	return m.getLastEntry()
}

// FloorEntry gets the entry corresponding to the specified key;
// if no such entry exists, returns the entry for the greatest key less than the specified key;
func (m *TreeMap[K, V]) FloorEntry(key K) *TreeEntry[K, V] {
	return m.getFloorEntry(key)
}

// FloorKey gets the specified key, returns the greatest key less than the specified key if not exist.
func (m *TreeMap[K, V]) FloorKey(key K) (K, bool) {
	var entry = m.getFloorEntry(key)
	if entry != nil {
		return entry.key, true
	}
	var zero K
	return zero, false
}

// CeilingEntry gets the entry corresponding to the specified key;
// returns the entry for the least key greater than the specified key if not exist.
func (m *TreeMap[K, V]) CeilingEntry(key K) *TreeEntry[K, V] {
	return m.getCeilingEntry(key)
}

// CeilingKey gets the specified key, return the least key greater than the specified key if not exist.
func (m *TreeMap[K, V]) CeilingKey(key K) (K, bool) {
	var entry = m.getCeilingEntry(key)
	if entry != nil {
		return entry.key, true
	}
	var zero K
	return zero, false
}

// HigherEntry gets the entry for the least key greater than the specified key
func (m *TreeMap[K, V]) HigherEntry(key K) *TreeEntry[K, V] {
	return m.getHigherEntry(key)
}

// HigherKey returns the least key greater than the specified key
func (m *TreeMap[K, V]) HigherKey(key K) (K, bool) {
	var entry = m.getHigherEntry(key)
	if entry != nil {
		return entry.key, true
	}
	var zero K
	return zero, false
}

// Foreach performs the given action for each entry in this map until all entries
// have been processed or the action panic
func (m *TreeMap[K, V]) Foreach(visit cutil.KeyValVisitor[K, V]) {
	var ver = m.version
	for e := m.getFirstEntry(); e != nil; e = successor(e) {
		visit(e.key, e.value)
		if ver != m.version {
			panic("concurrent map modification")
		}
	}
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (m *TreeMap[K, V]) Range(visit func(key K, value V) (shouldContinue bool)) {
	var ver = m.version
	for e := m.getFirstEntry(); e != nil; e = successor(e) {
		if !visit(e.key, e.value) {
			break
		}
		if ver != m.version {
			panic("concurrent map modification")
		}
	}
}

// Keys return list of all keys
func (m *TreeMap[K, V]) Keys() []K {
	var keys = make([]K, 0, m.size)
	for e := m.getFirstEntry(); e != nil; e = successor(e) {
		keys = append(keys, e.key)
	}
	return keys
}

// Values return list of all values
func (m *TreeMap[K, V]) Values() []V {
	var values = make([]V, 0, m.size)
	for e := m.getFirstEntry(); e != nil; e = successor(e) {
		values = append(values, e.value)
	}
	return values
}

func (m *TreeMap[K, V]) ToHashMap() map[K]V {
	var unordered = make(map[K]V, m.size)
	for e := m.getFirstEntry(); e != nil; e = successor(e) {
		unordered[e.key] = e.value
	}
	return unordered
}

//func (m *TreeMap[K, V]) Iterator() cutil.Iterator[*TreeEntry[K, V]] {
//	return NewEntryIterator(m, m.getFirstEntry())
//}
//
//func (m *TreeMap[K, V]) DescendingIterator() cutil.Iterator[*TreeEntry[K, V]] {
//	return NewKeyDescendingEntryIterator(m, m.getLastEntry())
//}
//
//func (m *TreeMap[K, V]) KeyIterator() cutil.Iterator[K] {
//	return NewKeyIterator(m, m.getFirstEntry())
//}
//
//func (m *TreeMap[K, V]) DescendingKeyIterator() cutil.Iterator[K] {
//	return NewDescendingKeyIterator(m, m.getLastEntry())
//}
//
//func (m *TreeMap[K, V]) ValueIterator() cutil.Iterator[V] {
//	return NewValueIterator(m, m.getFirstEntry())
//}

// Put associates the specified value with the specified key in this map.
// If the map previously contained a mapping for the key, the old value is replaced.
func (m *TreeMap[K, V]) Put(key K, value V) V {
	return m.put(key, value, true)
}

// Store sets the value for a key, equivalent to Put.
func (m *TreeMap[K, V]) Store(key K, value V) {
	m.put(key, value, true)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *TreeMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	var p = m.getEntry(key)
	if p == nil {
		actual = value
		m.put(key, value, true)
		return value, false
	} else {
		return p.value, true
	}
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *TreeMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	var p = m.getEntry(key)
	if p != nil {
		loaded = true
		previous = p.value
	}
	m.put(key, value, true)
	return
}

// PutIfAbsent put a key-value pair if the key is not already associated with a value.
func (m *TreeMap[K, V]) PutIfAbsent(key K, value V) V {
	return m.put(key, value, false)
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *TreeMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var p = m.getEntry(key)
	if p != nil {
		loaded = true
		value = p.value
		m.deleteEntry(p)
	}
	return
}

// CompareAndSwap swaps the old and new values for key if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (m *TreeMap[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	var p = m.getEntry(key)
	if p != nil && p.value == old {
		m.put(key, new, true)
		return true
	}
	return false
}

// Remove removes the mapping for this key from this TreeMap if present.
func (m *TreeMap[K, V]) Remove(key K) bool {
	var p = m.getEntry(key)
	if p != nil {
		m.deleteEntry(p)
		return true
	}
	return false
}

func (m *TreeMap[K, V]) Delete(key K) {
	m.Remove(key)
}

func (m *TreeMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	var p = m.getEntry(key)
	if p != nil && p.value == old {
		m.deleteEntry(p)
		return true
	}
	return false
}

// Clear removes all the mappings from this map.
func (m *TreeMap[K, V]) Clear() {
	m.version++
	m.size = 0
	m.root = nil
}

// Returns the first TreeEntry in the TreeMap (according to the key's order)
// Returns nil if the TreeMap is empty.
func (m *TreeMap[K, V]) getFirstEntry() *TreeEntry[K, V] {
	var p = m.root
	if p != nil {
		for p.left != nil {
			p = p.left
		}
	}
	return p
}

// Returns the last TreeEntry in the TreeMap (according to the key's order)
// Returns nil if the TreeMap is empty.
func (m *TreeMap[K, V]) getLastEntry() *TreeEntry[K, V] {
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
func (m *TreeMap[K, V]) getEntry(key K) *TreeEntry[K, V] {
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
func (m *TreeMap[K, V]) getCeilingEntry(key K) *TreeEntry[K, V] {
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
func (m *TreeMap[K, V]) getFloorEntry(key K) *TreeEntry[K, V] {
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
func (m *TreeMap[K, V]) getHigherEntry(key K) *TreeEntry[K, V] {
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
func (m *TreeMap[K, V]) getLowerEntry(key K) *TreeEntry[K, V] {
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

func (m *TreeMap[K, V]) put(key K, value V, replaceOld bool) V {
	var zero V
	var t = m.root
	if t == nil {
		m.addEntryToEmptyMap(key, value)
		return zero
	}
	var cmp = 0
	var parent *TreeEntry[K, V]
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

func (m *TreeMap[K, V]) addEntry(key K, value V, parent *TreeEntry[K, V], addToLeft bool) {
	var entry = NewTreeEntry(key, value, parent)
	if addToLeft {
		parent.left = entry
	} else {
		parent.right = entry
	}
	m.fixAfterInsertion(entry)
	m.size++
	m.version++
}

func (m *TreeMap[K, V]) addEntryToEmptyMap(key K, value V) {
	m.root = NewTreeEntry(key, value, nil)
	m.size = 1
	m.version++
}

func (m *TreeMap[K, V]) deleteEntry(p *TreeEntry[K, V]) {
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

func (m *TreeMap[K, V]) rotateLeft(p *TreeEntry[K, V]) {
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

func (m *TreeMap[K, V]) rotateRight(p *TreeEntry[K, V]) {
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

func (m *TreeMap[K, V]) fixAfterInsertion(x *TreeEntry[K, V]) {
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

func (m *TreeMap[K, V]) fixAfterDeletion(x *TreeEntry[K, V]) {
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

func keyOf[K, V comparable](e *TreeEntry[K, V]) (K, bool) {
	if e != nil {
		return e.key, true
	}
	var zero K
	return zero, false
}

func colorOf[K, V comparable](p *TreeEntry[K, V]) Color {
	if p != nil {
		return p.color
	}
	return BLACK
}

func parentOf[K, V comparable](p *TreeEntry[K, V]) *TreeEntry[K, V] {
	if p != nil {
		return p.parent
	}
	return nil
}

func setColor[K, V comparable](p *TreeEntry[K, V], color Color) {
	if p != nil {
		p.color = color
	}
}

func leftOf[K, V comparable](p *TreeEntry[K, V]) *TreeEntry[K, V] {
	if p != nil {
		return p.left
	}
	return nil
}

func rightOf[K, V comparable](p *TreeEntry[K, V]) *TreeEntry[K, V] {
	if p != nil {
		return p.right
	}
	return nil
}

// Returns the successor of the specified TreeEntry, or null if no such.
func successor[K, V comparable](t *TreeEntry[K, V]) *TreeEntry[K, V] {
	if t == nil {
		return nil
	} else if t.right != nil {
		var p = t.right
		for p.left != nil {
			p = p.left
		}
		return p
	} else {
		var p = t.parent
		var ch = t
		for p != nil && ch == p.right {
			ch = p
			p = p.parent
		}
		return p
	}
}

// Returns the predecessor of the specified TreeEntry, or null if no such.
func predecessor[K, V comparable](t *TreeEntry[K, V]) *TreeEntry[K, V] {
	if t == nil {
		return nil
	} else if t.left != nil {
		var p = t.left
		for p.right != nil {
			p = p.right
		}
		return p
	} else {
		var p = t.parent
		var ch = t
		for p != nil && ch == p.left {
			ch = p
			p = p.parent
		}
		return p
	}
}
