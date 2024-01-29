// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package treemap

type EntryIterator[K comparable, V any] struct {
	owner           *Map[K, V]
	next            *Entry[K, V]
	lastReturned    *Entry[K, V]
	expectedVersion int
}

func NewEntryIterator[K comparable, V any](m *Map[K, V], first *Entry[K, V]) *EntryIterator[K, V] {
	return &EntryIterator[K, V]{
		owner:           m,
		next:            first,
		expectedVersion: m.version,
	}
}

func (it *EntryIterator[K, V]) HasNext() bool {
	return it.next != nil
}

func (it *EntryIterator[K, V]) nextEntry() *Entry[K, V] {
	var e = it.next
	if e == nil {
		panic("EntryIterator: no such element")
	}
	if it.expectedVersion != it.owner.version {
		panic("EntryIterator: concurrent modification")
	}
	it.next = successor(e)
	it.lastReturned = e
	return e
}

func (it *EntryIterator[K, V]) prevEntry() *Entry[K, V] {
	var e = it.next
	if e == nil {
		panic("EntryIterator: no such element")
	}
	if it.expectedVersion != it.owner.version {
		panic("EntryIterator: concurrent modification")
	}
	it.next = predecessor(e)
	it.lastReturned = e
	return e
}

func (it *EntryIterator[K, V]) Next() *Entry[K, V] {
	return it.nextEntry()
}

// Remove removes from the underlying collection the last element returned
func (it *EntryIterator[K, V]) Remove() {
	if it.lastReturned == nil {
		panic("EntryIterator: illegal state")
	}
	if it.expectedVersion != it.owner.version {
		panic("EntryIterator: concurrent modification")
	}
	if it.lastReturned.left != nil && it.lastReturned.right != nil {
		it.next = it.lastReturned
	}
	it.owner.deleteEntry(it.lastReturned)
	it.expectedVersion = it.owner.version
	it.lastReturned = nil
}

type DescendingEntryIterator[K comparable, V any] struct {
	EntryIterator[K, V]
}

func NewKeyDescendingEntryIterator[K comparable, V any](m *Map[K, V], first *Entry[K, V]) *DescendingEntryIterator[K, V] {
	return &DescendingEntryIterator[K, V]{
		EntryIterator: EntryIterator[K, V]{
			owner:           m,
			next:            first,
			expectedVersion: m.version,
		},
	}
}

func (it *DescendingEntryIterator[K, V]) Next() *Entry[K, V] {
	return it.prevEntry()
}

type KeyIterator[K comparable, V any] struct {
	EntryIterator[K, V]
}

func NewKeyIterator[K comparable, V any](m *Map[K, V], first *Entry[K, V]) *KeyIterator[K, V] {
	return &KeyIterator[K, V]{
		EntryIterator: EntryIterator[K, V]{
			owner:           m,
			next:            first,
			expectedVersion: m.version,
		},
	}
}

func (it *KeyIterator[K, V]) Next() K {
	return it.nextEntry().key
}

type DescendingKeyIterator[K comparable, V any] struct {
	EntryIterator[K, V]
}

func NewDescendingKeyIterator[K comparable, V any](m *Map[K, V], first *Entry[K, V]) *DescendingKeyIterator[K, V] {
	return &DescendingKeyIterator[K, V]{
		EntryIterator: EntryIterator[K, V]{
			owner:           m,
			next:            first,
			expectedVersion: m.version,
		},
	}
}

func (it *DescendingKeyIterator[K, V]) Next() K {
	return it.prevEntry().key
}

func (it *DescendingKeyIterator[K, V]) Remove() {
	if it.lastReturned == nil {
		panic("DescendingKeyIterator: illegal state")
	}
	if it.expectedVersion != it.owner.version {
		panic("DescendingKeyIterator: concurrent modification")
	}
	it.owner.deleteEntry(it.lastReturned)
	it.lastReturned = nil
	it.expectedVersion = it.owner.version
}

type ValueIterator[K comparable, V any] struct {
	EntryIterator[K, V]
}

func NewValueIterator[K comparable, V any](m *Map[K, V], first *Entry[K, V]) *ValueIterator[K, V] {
	return &ValueIterator[K, V]{
		EntryIterator: EntryIterator[K, V]{
			owner:           m,
			next:            first,
			expectedVersion: m.version,
		},
	}
}

func (it *ValueIterator[K, V]) Next() V {
	return it.nextEntry().value
}
