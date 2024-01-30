// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package util

type Visitor[V any] func(value V) bool

type KeyValVisitor[K, V any] func(key K, value V) bool

type Iterator[T any] interface {
	// HasNext returns true if the iteration has more elements.
	HasNext() bool

	// Next returns the next element in the iteration.
	Next() T

	// Remove removes from the underlying collection the last element returned by this iterator.
	Remove()
}

// Container a base linear container interface
type Container[T any] interface {
	Len() int
	IsEmpty() bool
	Front() T
	Back() T

	PushBack(value T)
	PushFront(value T)
	PopFront() T
	PopBack() T
	Clear()
}
