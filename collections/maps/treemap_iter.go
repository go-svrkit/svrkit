// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package maps

type EntryIterator[K, V comparable] struct {
	owner           *TreeMap[K, V]
	next            *TreeEntry[K, V]
	lastReturned    *TreeEntry[K, V]
	expectedVersion int
}

func NewEntryIterator[K, V comparable](m *TreeMap[K, V], first *TreeEntry[K, V]) *EntryIterator[K, V] {
	return &EntryIterator[K, V]{
		owner:           m,
		next:            first,
		expectedVersion: m.version,
	}
}

func (it *EntryIterator[K, V]) HasNext() bool {
	return it.next != nil
}

func (it *EntryIterator[K, V]) nextEntry() *TreeEntry[K, V] {
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

func (it *EntryIterator[K, V]) prevEntry() *TreeEntry[K, V] {
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

func (it *EntryIterator[K, V]) Next() *TreeEntry[K, V] {
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

type DescendingEntryIterator[K, V comparable] struct {
	EntryIterator[K, V]
}

func NewKeyDescendingEntryIterator[K, V comparable](m *TreeMap[K, V], first *TreeEntry[K, V]) *DescendingEntryIterator[K, V] {
	return &DescendingEntryIterator[K, V]{
		EntryIterator: EntryIterator[K, V]{
			owner:           m,
			next:            first,
			expectedVersion: m.version,
		},
	}
}

func (it *DescendingEntryIterator[K, V]) Next() *TreeEntry[K, V] {
	return it.prevEntry()
}

type KeyIterator[K, V comparable] struct {
	EntryIterator[K, V]
}

func NewKeyIterator[K, V comparable](m *TreeMap[K, V], first *TreeEntry[K, V]) *KeyIterator[K, V] {
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

type DescendingKeyIterator[K, V comparable] struct {
	EntryIterator[K, V]
}

func NewDescendingKeyIterator[K, V comparable](m *TreeMap[K, V], first *TreeEntry[K, V]) *DescendingKeyIterator[K, V] {
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

type ValueIterator[K, V comparable] struct {
	EntryIterator[K, V]
}

func NewValueIterator[K, V comparable](m *TreeMap[K, V], first *TreeEntry[K, V]) *ValueIterator[K, V] {
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
