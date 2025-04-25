// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package tree

import (
	"cmp"
	"fmt"
	"iter"
	"strings"
)

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

type Entry[K comparable, V any] struct {
	Key                 K
	Value               V
	color               Color
	left, right, parent *Entry[K, V]
}

func NewTreeEntry[K comparable, V any](key K, val V, parent *Entry[K, V]) *Entry[K, V] {
	return &Entry[K, V]{
		Key:    key,
		Value:  val,
		parent: parent,
		color:  BLACK,
	}
}

func (e *Entry[K, V]) GetKey() K {
	return e.Key
}

func (e *Entry[K, V]) GetValue() V {
	return e.Value
}

func (e *Entry[K, V]) SetValue(val V) V {
	var old = e.Value
	e.Value = val
	return old
}

func (e *Entry[K, V]) Height() int {
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
func (e *Entry[K, V]) Size() int {
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

type Comparator[T any] func(a, b T) int

// Map is a Red-Black tree based map implementation.
// The map is sorted by a comparator passed to its constructor.
// This implementation provides guaranteed log(n) time cost for the
// Contains(), Get(), Put() and Remove() operations.
// Algorithms are adaptations of those in Cormen, Leiserson, and Rivest's <Introduction to Algorithms>.
type Map[K comparable, V any] struct {
	root    *Entry[K, V]  // root node of the tree
	compare Comparator[K] // The comparator used to maintain order in this tree map
	size    int           // The number of entries in the tree
	version int           // The number of structural modifications to the tree.
}

func NewMap[K comparable, V any](comparator Comparator[K]) *Map[K, V] {
	return &Map[K, V]{
		compare: comparator,
	}
}

func NewMapCmp[K cmp.Ordered, V any]() *Map[K, V] {
	return NewMap[K, V](cmp.Compare[K])
}

func MapOf[M ~map[K]V, K cmp.Ordered, V any](m M) *Map[K, V] {
	var treeMap = &Map[K, V]{
		compare: cmp.Compare[K],
	}
	for k, v := range m {
		treeMap.Put(k, v)
	}
	return treeMap
}

// Size returns the number of Key-Value mappings in this map.
func (m *Map[K, V]) Size() int {
	return m.size
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.size == 0
}

// Contains return true if this map contains a mapping for the specified Key
func (m *Map[K, V]) Contains(key K) bool {
	return m.getEntry(key) != nil
}

// Get returns the Value to which the specified Key is mapped,
// or nil if this map contains no mapping for the Key.
func (m *Map[K, V]) Get(key K) (V, bool) {
	var p = m.getEntry(key)
	if p != nil {
		return p.Value, true
	}
	var zero V
	return zero, false
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	return m.Get(key)
}

// GetOrDefault returns the Value to which the specified Key is mapped,
// or `defaultValue` if this map contains no mapping for the Key.
func (m *Map[K, V]) GetOrDefault(key K, defVal V) V {
	var p = m.getEntry(key)
	if p != nil {
		return p.Value
	}
	return defVal
}

func (m *Map[K, V]) GetEntry(key K) *Entry[K, V] {
	return m.getEntry(key)
}

// FirstKey returns the first Key in the Map (according to the Key's order)
func (m *Map[K, V]) FirstKey() (K, bool) {
	return keyOf[K](m.FirstEntry())
}

// LastKey returns the last Key in the Map (according to the Key's order)
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

// FloorEntry gets the entry corresponding to the specified Key;
// if no such entry exists, returns the entry for the greatest Key less than the specified Key;
func (m *Map[K, V]) FloorEntry(key K) *Entry[K, V] {
	return m.getFloorEntry(key)
}

// FloorKey gets the specified Key, returns the greatest Key less than the specified Key if not exist.
func (m *Map[K, V]) FloorKey(key K) (K, bool) {
	var entry = m.getFloorEntry(key)
	if entry != nil {
		return entry.Key, true
	}
	var zero K
	return zero, false
}

// CeilingEntry gets the entry corresponding to the specified Key;
// returns the entry for the least Key greater than the specified Key if not exist.
func (m *Map[K, V]) CeilingEntry(key K) *Entry[K, V] {
	return m.getCeilingEntry(key)
}

// CeilingKey gets the specified Key, return the least Key greater than the specified Key if not exist.
func (m *Map[K, V]) CeilingKey(key K) (K, bool) {
	var entry = m.getCeilingEntry(key)
	if entry != nil {
		return entry.Key, true
	}
	var zero K
	return zero, false
}

// HigherEntry gets the entry for the least Key greater than the specified Key
func (m *Map[K, V]) HigherEntry(key K) *Entry[K, V] {
	return m.getHigherEntry(key)
}

// HigherKey returns the least Key greater than the specified Key
func (m *Map[K, V]) HigherKey(key K) (K, bool) {
	var entry = m.getHigherEntry(key)
	if entry != nil {
		return entry.Key, true
	}
	var zero K
	return zero, false
}

// Range calls f sequentially for each Key and Value present in the map.
// If f returns false, range stops the iteration.
func (m *Map[K, V]) Range(visit func(key K, value V) (shouldContinue bool)) {
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

func (m *Map[K, V]) IterSeq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.getFirstEntry(); e != nil; e = successor(e) {
			if !yield(e.Key, e.Value) {
				break
			}
		}
	}
}

func (m *Map[K, V]) IterBackwards() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.getLastEntry(); e != nil; e = predecessor(e) {
			if !yield(e.Key, e.Value) {
				break
			}
		}
	}
}

// Keys return list of all keys
func (m *Map[K, V]) Keys() []K {
	var keys = make([]K, 0, m.size)
	for e := m.FirstEntry(); e != nil; e = successor(e) {
		keys = append(keys, e.Key)
	}
	return keys
}

// Values return list of all values
func (m *Map[K, V]) Values() []V {
	var values = make([]V, 0, m.size)
	for e := m.FirstEntry(); e != nil; e = successor(e) {
		values = append(values, e.Value)
	}
	return values
}

func (m *Map[K, V]) CloneMap() map[K]V {
	var unordered = make(map[K]V, m.size)
	for e := m.FirstEntry(); e != nil; e = successor(e) {
		unordered[e.Key] = e.Value
	}
	return unordered
}

func (m *Map[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	var cnt = 0
	for k, v := range m.IterSeq() {
		fmt.Fprintf(&sb, "%v:%v", k, v)
		cnt++
		if cnt < m.size {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

// Put associates the specified Value with the specified Key in this map.
// If the map previously contained a mapping for the Key, the old Value is replaced.
func (m *Map[K, V]) Put(key K, value V) V {
	return m.put(key, value, true)
}

// Store sets the Value for a Key, equivalent to Put.
func (m *Map[K, V]) Store(key K, value V) {
	m.put(key, value, true)
}

// LoadOrStore returns the existing Value for the Key if present.
// Otherwise, it stores and returns the given Value.
// The loaded result is true if the Value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
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
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	var p = m.getEntry(key)
	if p != nil {
		loaded = true
		previous = p.Value
	}
	m.put(key, value, true)
	return
}

// PutIfAbsent put a Key-Value pair if the Key is not already associated with a Value.
func (m *Map[K, V]) PutIfAbsent(key K, value V) V {
	return m.put(key, value, false)
}

// LoadAndDelete deletes the Value for a Key, returning the previous Value if any.
// The loaded result reports whether the Key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var p = m.getEntry(key)
	if p != nil {
		loaded = true
		value = p.Value
		m.deleteEntry(p)
	}
	return
}

// Remove removes the mapping for this Key from this Map if present.
func (m *Map[K, V]) Remove(key K) bool {
	var p = m.getEntry(key)
	if p != nil {
		m.deleteEntry(p)
		return true
	}
	return false
}

func (m *Map[K, V]) Delete(key K) {
	m.Remove(key)
}

// Clear removes all the mappings from this map.
func (m *Map[K, V]) Clear() {
	m.version++
	m.size = 0
	m.root = nil
}

// Returns the first Entry in the Map (according to the Key's order)
// Returns nil if the Map is empty.
func (m *Map[K, V]) getFirstEntry() *Entry[K, V] {
	var p = m.root
	if p != nil {
		for p.left != nil {
			p = p.left
		}
	}
	return p
}

// Returns the last Entry in the Map (according to the Key's order)
// Returns nil if the Map is empty.
func (m *Map[K, V]) getLastEntry() *Entry[K, V] {
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
func (m *Map[K, V]) getEntry(key K) *Entry[K, V] {
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
func (m *Map[K, V]) getCeilingEntry(key K) *Entry[K, V] {
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
func (m *Map[K, V]) getFloorEntry(key K) *Entry[K, V] {
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
func (m *Map[K, V]) getHigherEntry(key K) *Entry[K, V] {
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
func (m *Map[K, V]) getLowerEntry(key K) *Entry[K, V] {
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

func (m *Map[K, V]) put(key K, value V, replaceOld bool) V {
	var t = m.root
	if t == nil {
		m.addEntryToEmptyMap(key, value)
		var zero V
		return zero
	}
	var cp = 0
	var parent *Entry[K, V]
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

func (m *Map[K, V]) addEntry(key K, value V, parent *Entry[K, V], addToLeft bool) {
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

func (m *Map[K, V]) addEntryToEmptyMap(key K, value V) {
	m.root = NewTreeEntry(key, value, nil)
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

func keyOf[K comparable, V any](e *Entry[K, V]) (K, bool) {
	if e != nil {
		return e.Key, true
	}
	var zero K
	return zero, false
}

func colorOf[K comparable, V any](p *Entry[K, V]) Color {
	if p != nil {
		return p.color
	}
	return BLACK
}

func parentOf[K comparable, V any](p *Entry[K, V]) *Entry[K, V] {
	if p != nil {
		return p.parent
	}
	return nil
}

func setColor[K comparable, V any](p *Entry[K, V], color Color) {
	if p != nil {
		p.color = color
	}
}

func leftOf[K comparable, V any](p *Entry[K, V]) *Entry[K, V] {
	if p != nil {
		return p.left
	}
	return nil
}

func rightOf[K comparable, V any](p *Entry[K, V]) *Entry[K, V] {
	if p != nil {
		return p.right
	}
	return nil
}

// Returns the successor of the specified Entry, or null if no such.
func successor[K comparable, V any](t *Entry[K, V]) *Entry[K, V] {
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

// Returns the predecessor of the specified Entry, or null if no such.
func predecessor[K comparable, V any](t *Entry[K, V]) *Entry[K, V] {
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
