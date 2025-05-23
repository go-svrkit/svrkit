package queue

import (
	"log"
)

// Deque generalizes a queue and a stack, to efficiently add and remove items at
// either end with O(1) performance.  Queue (FIFO) operations are supported using
// PushBack() and PopFront().  Stack (LIFO) operations are supported using
// PushBack() and PopBack().
//
// Ring-buffer Performance
//
// The ring-buffer automatically resizes by
// powers of two, growing when additional capacity is needed and shrinking when
// only a quarter of the capacity is used, and uses bitwise arithmetic for all
// calculations.
//
// The ring-buffer implementation significantly improves memory and time
// performance with fewer GC pauses, compared to implementations based on slices
// and linked lists.
//
// For maximum speed, this deque implementation leaves concurrency safety up to
// the application to provide, however the application chooses, if needed at all.
//
// Reading Empty Deque
//
// Since it is OK for the deque to contain a nil Value, it is necessary to either
// panic or return a second boolean Value to indicate the deque is empty, when
// reading or removing an element.  This deque panics when reading from an empty
// deque.  This is a run-time check to help catch programming errors, which may be
// missed if a second return Value is ignored.  Simply check Deque.Len() before
// reading from the deque.

// minCapacity is the smallest capacity that deque may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

// Deque provides a fast ring-buffer deque (double-ended queue)
// implementation.
//

type Deque[T any] struct {
	buf    []T
	head   int
	tail   int
	count  int
	minCap int
}

// NewDeque creates a new Deque, optionally setting the current and minimum capacity
// when non-zero values are given for these.
//
// To create a Deque with capacity to store 2048 items without resizing, and
// that will not resize below space for 32 items when removing itmes:
//
//	d := deque.NewList(2048, 32)
//
// To create a Deque that has not yet allocated memory, but after it does will
// never resize to have space for less than 64 items:
//
//	d := deque.NewList(0, 64)
//
// Note that interface{} values supplied here are rounded up to the nearest power of 2.
func NewDeque[T any](size ...int) *Deque[T] {
	var capacity, minimum int
	if len(size) >= 1 {
		capacity = size[0]
		if len(size) >= 2 {
			minimum = size[1]
		}
	}

	minCap := minCapacity
	for minCap < minimum {
		minCap <<= 1
	}

	var buf []T
	if capacity != 0 {
		bufSize := minCap
		for bufSize < capacity {
			bufSize <<= 1
		}
		buf = make([]T, bufSize)
	}

	return &Deque[T]{
		buf:    buf,
		minCap: minCap,
	}
}

// Cap returns the current capacity of the Deque.
func (q *Deque[T]) Cap() int {
	return len(q.buf)
}

// Len returns the number of elements currently stored in the queue.
func (q *Deque[T]) Len() int {
	return q.count
}

func (q *Deque[T]) IsEmpty() bool {
	return q.count == 0
}

// PushBack appends an element to the back of the queue.  Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (q *Deque[T]) PushBack(elem T) {
	q.growIfFull()

	q.buf[q.tail] = elem
	// Calculate new tail position.
	q.tail = q.next(q.tail)
	q.count++
}

// PushFront prepends an element to the front of the queue.
func (q *Deque[T]) PushFront(elem T) {
	q.growIfFull()

	// Calculate new head position.
	q.head = q.prev(q.head)
	q.buf[q.head] = elem
	q.count++
}

// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *Deque[T]) PopFront() T {
	if q.count <= 0 {
		log.Panicln("deque: PopFront() called on empty queue")
	}
	var zero T
	ret := q.buf[q.head]
	q.buf[q.head] = zero
	// Calculate new head position.
	q.head = q.next(q.head)
	q.count--

	q.shrinkIfExcess()
	return ret
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *Deque[T]) PopBack() T {
	if q.count <= 0 {
		log.Panicln("deque: PopBack() called on empty queue")
	}

	// Calculate new tail position
	q.tail = q.prev(q.tail)

	// Remove Value at tail.
	var zero T
	ret := q.buf[q.tail]
	q.buf[q.tail] = zero
	q.count--

	q.shrinkIfExcess()
	return ret
}

// Front returns the element at the front of the queue.  This is the element
// that would be returned by PopFront().  This call panics if the queue is
// empty.
func (q *Deque[T]) Front() T {
	if q.count <= 0 {
		log.Panicln("deque: Front() called when empty")
	}
	return q.buf[q.head]
}

// Back returns the element at the back of the queue.  This is the element
// that would be returned by PopBack().  This call panics if the queue is
// empty.
func (q *Deque[T]) Back() T {
	if q.count <= 0 {
		log.Panicln("deque: Back() called when empty")
	}
	return q.buf[q.prev(q.tail)]
}

// At returns the element at index i in the queue without removing the element
// from the queue.  This method accepts only non-negative index values.  At(0)
// refers to the first element and is the same as Front().  At(Len()-1) refers
// to the last element and is the same as Back().  If the index is invalid, the
// call panics.
//
// The purpose of At is to allow Deque to serve as a more general purpose
// circular buffer, where items are only added to and removed from the ends of
// the deque, but may be read from T place within the deque.  Consider the
// case of a fixed-size circular log buffer: A new entry is pushed onto one end
// and when full the oldest is popped from the other end.  All the log entries
// in the buffer must be readable without altering the buffer contents.
func (q *Deque[T]) At(i int) T {
	if i < 0 || i >= q.count {
		log.Panicln("deque: At() called with index out of range")
	}
	// bitwise modulus
	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Set puts the element at index i in the queue. Set shares the same purpose
// than At() but perform the opposite operation. The index i is the same
// index defined by At(). If the index is invalid, the call panics.
func (q *Deque[T]) Set(i int, elem T) {
	if i < 0 || i >= q.count {
		log.Panicln("deque: Set() called with index out of range")
	}
	// bitwise modulus
	q.buf[(q.head+i)&(len(q.buf)-1)] = elem
}

// Clear removes all elements from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse.  The queue will not be resized smaller as long as items are
// only added.  Only when items are removed is the queue subject to getting
// resized smaller.
func (q *Deque[T]) Clear() {
	// bitwise modulus
	var zero T
	modBits := len(q.buf) - 1
	for h := q.head; h != q.tail; h = (h + 1) & modBits {
		q.buf[h] = zero
	}
	q.head = 0
	q.tail = 0
	q.count = 0
}

// Rotate rotates the deque n steps front-to-back.  If n is negative, rotates
// back-to-front.  Having Deque provide Rotate() avoids resizing that could
// happen if implementing rotation using only Pop and Push methods.
func (q *Deque[T]) Rotate(n int) {
	if q.count <= 1 {
		return
	}
	// Rotating a multiple of q.count is same as no rotation.
	n %= q.count
	if n == 0 {
		return
	}

	modBits := len(q.buf) - 1
	// If no empty space in buffer, only move head and tail indexes.
	if q.head == q.tail {
		// Calculate new head and tail using bitwise modulus.
		q.head = (q.head + n) & modBits
		q.tail = (q.tail + n) & modBits
		return
	}

	var zero T
	if n < 0 {
		// Rotate back to front.
		for ; n < 0; n++ {
			// Calculate new head and tail using bitwise modulus.
			q.head = (q.head - 1) & modBits
			q.tail = (q.tail - 1) & modBits
			// Put tail Value at head and remove Value at tail.
			q.buf[q.head] = q.buf[q.tail]
			q.buf[q.tail] = zero
		}
		return
	}

	// Rotate front to back.
	for ; n > 0; n-- {
		// Put head Value at tail and remove Value at head.
		q.buf[q.tail] = q.buf[q.head]
		q.buf[q.head] = zero
		// Calculate new head and tail using bitwise modulus.
		q.head = (q.head + 1) & modBits
		q.tail = (q.tail + 1) & modBits
	}
}

// SetMinCapacity sets a minimum capacity of 2^minCapacityExp.  If the Value of
// the minimum capacity is less than or equal to the minimum allowed, then
// capacity is set to the minimum allowed.  This may be called at anytime to
// set a new minimum capacity.
//
// Setting a larger minimum capacity may be used to prevent resizing when the
// number of stored items changes frequently across a wide range.
func (q *Deque[T]) SetMinCapacity(minCapacityExp uint) {
	if 1<<minCapacityExp > minCapacity {
		q.minCap = 1 << minCapacityExp
	} else {
		q.minCap = minCapacity
	}
}

// prev returns the previous buffer position wrapping around buffer.
func (q *Deque[T]) prev(i int) int {
	return (i - 1) & (len(q.buf) - 1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (q *Deque[T]) next(i int) int {
	return (i + 1) & (len(q.buf) - 1) // bitwise modulus
}

// growIfFull resizes up if the buffer is full.
func (q *Deque[T]) growIfFull() {
	if q.count != len(q.buf) {
		return
	}
	if len(q.buf) == 0 {
		if q.minCap == 0 {
			q.minCap = minCapacity
		}
		q.buf = make([]T, q.minCap)
		return
	}
	q.resize()
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (q *Deque[T]) shrinkIfExcess() {
	if len(q.buf) > q.minCap && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

// resize resizes the deque to fit exactly twice its current contents.  This is
// used to grow the queue when it is full, and also to shrink it when it is
// only a quarter full.
func (q *Deque[T]) resize() {
	newBuf := make([]T, q.count<<1)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}
