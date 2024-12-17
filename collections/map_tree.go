// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package collections

import (
	"cmp"
	"iter"
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

type TreeEntry[K comparable, V any] struct {
	Key                 K
	Value               V
	color               Color
	left, right, parent *TreeEntry[K, V]
}

func NewTreeEntry[K comparable, V any](key K, val V, parent *TreeEntry[K, V]) *TreeEntry[K, V] {
	return &TreeEntry[K, V]{
		Key:    key,
		Value:  val,
		parent: parent,
		color:  BLACK,
	}
}

func (e *TreeEntry[K, V]) GetKey() K {
	return e.Key
}

func (e *TreeEntry[K, V]) GetValue() V {
	return e.Value
}

func (e *TreeEntry[K, V]) SetValue(val V) V {
	var old = e.Value
	e.Value = val
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

// Size returns the number of elements stored in the subtree.
func (e *TreeEntry[K, V]) Size() int {
	if e == nil {
		return 0
	}
	var size = 1
	if e.left != nil {
		size += e.left.Size()
	}
	if e.right != nil {
		size += e.right.Size()
	}
	return size
}

type TreeMap[K comparable, V any] struct {
	root    *TreeEntry[K, V] // root node of the tree
	compare Comparator[K]    // The comparator used to maintain order in this tree map
	size    int              // The number of entries in the tree
	version int              // The number of structural modifications to the tree.
}

var _ MapInterface[int, int] = (*TreeMap[int, int])(nil)

func NewTreeMap[K comparable, V any](comparator Comparator[K]) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		compare: comparator,
	}
}

func NewTreeMapCmp[K cmp.Ordered, V any]() *TreeMap[K, V] {
	return NewTreeMap[K, V](cmp.Compare[K])
}

func TreeMapOf[M ~map[K]V, K cmp.Ordered, V any](m M) *TreeMap[K, V] {
	var treeMap = &TreeMap[K, V]{
		compare: cmp.Compare[K],
	}
	for k, v := range m {
		treeMap.Put(k, v)
	}
	return treeMap
}

// Size returns the number of Key-Value mappings in this map.
func (m *TreeMap[K, V]) Size() int {
	return m.size
}

func (m *TreeMap[K, V]) IsEmpty() bool {
	return m.size == 0
}

// ContainsKey return true if this map contains a mapping for the specified Key
func (m *TreeMap[K, V]) ContainsKey(key K) bool {
	return m.getEntry(key) != nil
}

// Get returns the Value to which the specified Key is mapped,
// or nil if this map contains no mapping for the Key.
func (m *TreeMap[K, V]) Get(key K) (V, bool) {
	var p = m.getEntry(key)
	if p != nil {
		return p.Value, true
	}
	var zero V
	return zero, false
}

func (m *TreeMap[K, V]) Load(key K) (V, bool) {
	return m.Get(key)
}

// GetOrDefault returns the Value to which the specified Key is mapped,
// or `defaultValue` if this map contains no mapping for the Key.
func (m *TreeMap[K, V]) GetOrDefault(key K, defVal V) V {
	var p = m.getEntry(key)
	if p != nil {
		return p.Value
	}
	return defVal
}

func (m *TreeMap[K, V]) GetEntry(key K) *TreeEntry[K, V] {
	return m.getEntry(key)
}

// FirstKey returns the first Key in the TreeMap (according to the Key's order)
func (m *TreeMap[K, V]) FirstKey() (K, bool) {
	return keyOf[K](m.FirstEntry())
}

// LastKey returns the last Key in the TreeMap (according to the Key's order)
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

// FloorEntry gets the entry corresponding to the specified Key;
// if no such entry exists, returns the entry for the greatest Key less than the specified Key;
func (m *TreeMap[K, V]) FloorEntry(key K) *TreeEntry[K, V] {
	return m.getFloorEntry(key)
}

// FloorKey gets the specified Key, returns the greatest Key less than the specified Key if not exist.
func (m *TreeMap[K, V]) FloorKey(key K) (K, bool) {
	var entry = m.getFloorEntry(key)
	if entry != nil {
		return entry.Key, true
	}
	var zero K
	return zero, false
}

// CeilingEntry gets the entry corresponding to the specified Key;
// returns the entry for the least Key greater than the specified Key if not exist.
func (m *TreeMap[K, V]) CeilingEntry(key K) *TreeEntry[K, V] {
	return m.getCeilingEntry(key)
}

// CeilingKey gets the specified Key, return the least Key greater than the specified Key if not exist.
func (m *TreeMap[K, V]) CeilingKey(key K) (K, bool) {
	var entry = m.getCeilingEntry(key)
	if entry != nil {
		return entry.Key, true
	}
	var zero K
	return zero, false
}

// HigherEntry gets the entry for the least Key greater than the specified Key
func (m *TreeMap[K, V]) HigherEntry(key K) *TreeEntry[K, V] {
	return m.getHigherEntry(key)
}

// HigherKey returns the least Key greater than the specified Key
func (m *TreeMap[K, V]) HigherKey(key K) (K, bool) {
	var entry = m.getHigherEntry(key)
	if entry != nil {
		return entry.Key, true
	}
	var zero K
	return zero, false
}

// Foreach performs the given action for each entry in this map until all entries
// have been processed or the action panic
func (m *TreeMap[K, V]) Foreach(visit KeyValVisitor[K, V]) {
	var ver = m.version
	for e := m.FirstEntry(); e != nil; e = successor(e) {
		visit(e.Key, e.Value)
		if ver != m.version {
			panic("concurrent map modification")
		}
	}
}

// Range calls f sequentially for each Key and Value present in the map.
// If f returns false, range stops the iteration.
func (m *TreeMap[K, V]) Range(visit func(key K, value V) (shouldContinue bool)) {
	var ver = m.version
	for e := m.FirstEntry(); e != nil; e = successor(e) {
		if !visit(e.Key, e.Value) {
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
	for e := m.FirstEntry(); e != nil; e = successor(e) {
		keys = append(keys, e.Key)
	}
	return keys
}

// Values return list of all values
func (m *TreeMap[K, V]) Values() []V {
	var values = make([]V, 0, m.size)
	for e := m.FirstEntry(); e != nil; e = successor(e) {
		values = append(values, e.Value)
	}
	return values
}

func (m *TreeMap[K, V]) ToMap() map[K]V {
	var unordered = make(map[K]V, m.size)
	for e := m.FirstEntry(); e != nil; e = successor(e) {
		unordered[e.Key] = e.Value
	}
	return unordered
}

func (m *TreeMap[K, V]) IterSeq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.getFirstEntry(); e != nil; e = successor(e) {
			if !yield(e.Key, e.Value) {
				break
			}
		}
	}
}

func (m *TreeMap[K, V]) IterBackwards() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.getLastEntry(); e != nil; e = predecessor(e) {
			if !yield(e.Key, e.Value) {
				break
			}
		}
	}
}

// Put associates the specified Value with the specified Key in this map.
// If the map previously contained a mapping for the Key, the old Value is replaced.
func (m *TreeMap[K, V]) Put(key K, value V) V {
	return m.put(key, value, true)
}

// Store sets the Value for a Key, equivalent to Put.
func (m *TreeMap[K, V]) Store(key K, value V) {
	m.put(key, value, true)
}

// LoadOrStore returns the existing Value for the Key if present.
// Otherwise, it stores and returns the given Value.
// The loaded result is true if the Value was loaded, false if stored.
func (m *TreeMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	var p = m.getEntry(key)
	if p == nil {
		actual = value
		m.put(key, value, true)
		return value, false
	} else {
		return p.Value, true
	}
}

// Swap swaps the Value for a Key and returns the previous Value if any.
// The loaded result reports whether the Key was present.
func (m *TreeMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	var p = m.getEntry(key)
	if p != nil {
		loaded = true
		previous = p.Value
	}
	m.put(key, value, true)
	return
}

// PutIfAbsent put a Key-Value pair if the Key is not already associated with a Value.
func (m *TreeMap[K, V]) PutIfAbsent(key K, value V) V {
	return m.put(key, value, false)
}

// LoadAndDelete deletes the Value for a Key, returning the previous Value if any.
// The loaded result reports whether the Key was present.
func (m *TreeMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var p = m.getEntry(key)
	if p != nil {
		loaded = true
		value = p.Value
		m.deleteEntry(p)
	}
	return
}

// Remove removes the mapping for this Key from this TreeMap if present.
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

// Clear removes all the mappings from this map.
func (m *TreeMap[K, V]) Clear() {
	m.version++
	m.size = 0
	m.root = nil
}

// Returns the first TreeEntry in the TreeMap (according to the Key's order)
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

// Returns the last TreeEntry in the TreeMap (according to the Key's order)
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

// Returns this map's entry for the given Key,
// or nil if the map does not contain an entry for the Key.
func (m *TreeMap[K, V]) getEntry(key K) *TreeEntry[K, V] {
	var p = m.root
	for p != nil {
		var cp = m.compare(key, p.Key)
		if cp < 0 {
			p = p.left
		} else if cp > 0 {
			p = p.right
		} else {
			return p
		}
	}
	return nil
}

// Gets the entry corresponding to the specified Key;
// if no such entry exists, returns the entry for the least Key greater than the specified Key;
// if no such entry exists returns nil.
func (m *TreeMap[K, V]) getCeilingEntry(key K) *TreeEntry[K, V] {
	var p = m.root
	for p != nil {
		var cp = m.compare(key, p.Key)
		if cp < 0 {
			if p.left != nil {
				p = p.left
			} else {
				return p
			}
		} else if cp > 0 {
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

// Gets the entry corresponding to the specified Key;
// if no such entry exists, returns the entry for the greatest Key less than the specified Key;
// if no such entry exists, returns nil.
func (m *TreeMap[K, V]) getFloorEntry(key K) *TreeEntry[K, V] {
	var p = m.root
	for p != nil {
		var cp = m.compare(key, p.Key)
		if cp > 0 {
			if p.right != nil {
				p = p.right
			} else {
				return p
			}
		} else if cp < 0 {
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

// Gets the entry for the least Key greater than the specified Key;
// if no such entry exists, returns the entry for the least Key greater than the specified Key;
// if no such entry exists returns nil.
func (m *TreeMap[K, V]) getHigherEntry(key K) *TreeEntry[K, V] {
	var p = m.root
	for p != nil {
		var cp = m.compare(key, p.Key)
		if cp < 0 {
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

// Returns the entry for the greatest Key less than the specified Key;
// if no such entry exists (i.e., the least Key in the Tree is greater than the specified Key), returns nil
func (m *TreeMap[K, V]) getLowerEntry(key K) *TreeEntry[K, V] {
	var p = m.root
	for p != nil {
		var cp = m.compare(key, p.Key)
		if cp > 0 {
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
	var t = m.root
	if t == nil {
		m.addEntryToEmptyMap(key, value)
		var zero V
		return zero
	}
	var cp = 0
	var parent *TreeEntry[K, V]
	var loop = true
	for loop {
		parent = t
		cp = m.compare(key, t.Key)
		if cp < 0 {
			t = t.left
		} else if cp > 0 {
			t = t.right
		} else {
			var oldValue = t.Value
			if replaceOld {
				t.Value = value
			}
			return oldValue
		}
		loop = t != nil
	}
	m.addEntry(key, value, parent, cp < 0)
	var zero V
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
		p.Key = s.Key
		p.Value = s.Value
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
// Implementations of re-balancings during insertion and deletion are
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

func keyOf[K comparable, V any](e *TreeEntry[K, V]) (K, bool) {
	if e != nil {
		return e.Key, true
	}
	var zero K
	return zero, false
}

func colorOf[K comparable, V any](p *TreeEntry[K, V]) Color {
	if p != nil {
		return p.color
	}
	return BLACK
}

func parentOf[K comparable, V any](p *TreeEntry[K, V]) *TreeEntry[K, V] {
	if p != nil {
		return p.parent
	}
	return nil
}

func setColor[K comparable, V any](p *TreeEntry[K, V], color Color) {
	if p != nil {
		p.color = color
	}
}

func leftOf[K comparable, V any](p *TreeEntry[K, V]) *TreeEntry[K, V] {
	if p != nil {
		return p.left
	}
	return nil
}

func rightOf[K comparable, V any](p *TreeEntry[K, V]) *TreeEntry[K, V] {
	if p != nil {
		return p.right
	}
	return nil
}

// Returns the successor of the specified TreeEntry, or null if no such.
func successor[K comparable, V any](t *TreeEntry[K, V]) *TreeEntry[K, V] {
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
func predecessor[K comparable, V any](t *TreeEntry[K, V]) *TreeEntry[K, V] {
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

type TreeIterator[K comparable, V any] struct {
	tree     *TreeMap[K, V]
	node     *TreeEntry[K, V]
	position int8
}
