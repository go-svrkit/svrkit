// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package list

import "testing"

func checkListLen[T comparable](t *testing.T, l *List[T], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers[T comparable](t *testing.T, l *List[T], es []*Elem[T]) {
	root := &l.root

	if !checkListLen(t, l, len(es)) {
		return
	}

	// zero length lists must be the zero Value or properly initialized (sentinel circle)
	if len(es) == 0 {
		if l.root.next != nil && l.root.next != root || l.root.prev != nil && l.root.prev != root {
			t.Errorf("l.root.next = %p, l.root.prev = %p; both should both be nil or %p", l.root.next, l.root.prev, root)
		}
		return
	}
	// len(es) > 0

	// check internal and external prev/next connections
	for i, e := range es {
		prev := root
		Prev := (*Elem[T])(nil)
		if i > 0 {
			prev = es[i-1]
			Prev = prev
		}
		if p := e.prev; p != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p", i, e, p, prev)
		}
		if p := e.Prev(); p != Prev {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, e, p, Prev)
		}

		next := root
		Next := (*Elem[T])(nil)
		if i < len(es)-1 {
			next = es[i+1]
			Next = next
		}
		if n := e.next; n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, e, n, next)
		}
		if n := e.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, e, n, Next)
		}
	}
}

func TestList(t *testing.T) {
	l := New[string]()
	checkListPointers(t, l, []*Elem[string]{})

	// Single element list
	e := l.PushFront("a")
	checkListPointers(t, l, []*Elem[string]{e})
	l.MoveToFront(e)
	checkListPointers(t, l, []*Elem[string]{e})
	l.MoveToBack(e)
	checkListPointers(t, l, []*Elem[string]{e})
	l.Remove(e)
	checkListPointers(t, l, []*Elem[string]{})

	// Bigger list
	e2 := l.PushFront("2")
	e1 := l.PushFront("1")
	e3 := l.PushBack("3")
	e4 := l.PushBack("banana")
	checkListPointers(t, l, []*Elem[string]{e1, e2, e3, e4})

	l.Remove(e2)
	checkListPointers(t, l, []*Elem[string]{e1, e3, e4})

	l.MoveToFront(e3) // move from middle
	checkListPointers(t, l, []*Elem[string]{e3, e1, e4})

	l.MoveToFront(e1)
	l.MoveToBack(e3) // move from middle
	checkListPointers(t, l, []*Elem[string]{e1, e4, e3})

	l.MoveToFront(e3) // move from back
	checkListPointers(t, l, []*Elem[string]{e3, e1, e4})
	l.MoveToFront(e3) // should be no-op
	checkListPointers(t, l, []*Elem[string]{e3, e1, e4})

	l.MoveToBack(e3) // move from front
	checkListPointers(t, l, []*Elem[string]{e1, e4, e3})
	l.MoveToBack(e3) // should be no-op
	checkListPointers(t, l, []*Elem[string]{e1, e4, e3})

	e2 = l.InsertBefore("2", e1) // insert before front
	checkListPointers(t, l, []*Elem[string]{e2, e1, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore("2", e4) // insert before middle
	checkListPointers(t, l, []*Elem[string]{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore("2", e3) // insert before back
	checkListPointers(t, l, []*Elem[string]{e1, e4, e2, e3})
	l.Remove(e2)

	e2 = l.InsertAfter("2", e1) // insert after front
	checkListPointers(t, l, []*Elem[string]{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertAfter("2", e4) // insert after middle
	checkListPointers(t, l, []*Elem[string]{e1, e4, e2, e3})
	l.Remove(e2)
	e2 = l.InsertAfter("2", e3) // insert after back
	checkListPointers(t, l, []*Elem[string]{e1, e4, e3, e2})
	l.Remove(e2)

	// Clear all elements by iterating
	var next *Elem[string]
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		l.Remove(e)
	}
	checkListPointers(t, l, []*Elem[string]{})
}

func checkList[T comparable](t *testing.T, l *List[T], es []T) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		le := e.Value
		if le != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, es[i])
		}
		i++
	}
}

func TestExtending(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := New[int]()
	l3.PushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushBackList(l2)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	l3 = New[int]()
	l3.PushFrontList(l2)
	checkList(t, l3, []int{4, 5})
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	checkList(t, l1, []int{1, 2, 3})
	checkList(t, l2, []int{4, 5})

	l3 = New[int]()
	l3.PushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushBackList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = New[int]()
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushFrontList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = New[int]()
	l1.PushBackList(l3)
	checkList(t, l1, []int{1, 2, 3})
	l1.PushFrontList(l3)
	checkList(t, l1, []int{1, 2, 3})
}

func TestRemove(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	checkListPointers(t, l, []*Elem[int]{e1, e2})
	e := l.Front()
	l.Remove(e)
	checkListPointers(t, l, []*Elem[int]{e2})
	l.Remove(e)
	checkListPointers(t, l, []*Elem[int]{e2})
}

func TestIssue4103(t *testing.T) {
	l1 := New[int]()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := New[int]()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	l2.Remove(e) // l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(8, e)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()
	l.Remove(e)
	if e.Value != 1 {
		t.Errorf("e.Value = %d, want 1", e.Value)
	}
	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.Prev() != nil {
		t.Errorf("e.Prev() != nil")
	}
}

func TestMove(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	e4 := l.PushBack(4)

	l.MoveAfter(e3, e3)
	checkListPointers(t, l, []*Elem[int]{e1, e2, e3, e4})
	l.MoveBefore(e2, e2)
	checkListPointers(t, l, []*Elem[int]{e1, e2, e3, e4})

	l.MoveAfter(e3, e2)
	checkListPointers(t, l, []*Elem[int]{e1, e2, e3, e4})
	l.MoveBefore(e2, e3)
	checkListPointers(t, l, []*Elem[int]{e1, e2, e3, e4})

	l.MoveBefore(e2, e4)
	checkListPointers(t, l, []*Elem[int]{e1, e3, e2, e4})
	e2, e3 = e3, e2

	l.MoveBefore(e4, e1)
	checkListPointers(t, l, []*Elem[int]{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.MoveAfter(e4, e1)
	checkListPointers(t, l, []*Elem[int]{e1, e4, e2, e3})
	e2, e3, e4 = e4, e2, e3

	l.MoveAfter(e2, e3)
	checkListPointers(t, l, []*Elem[int]{e1, e3, e2, e4})
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized List
func TestZeroList(t *testing.T) {
	var l1 = new(List[int])
	l1.PushFront(1)
	checkList(t, l1, []int{1})

	var l2 = new(List[int])
	l2.PushBack(1)
	checkList(t, l2, []int{1})

	var l3 = new(List[int])
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1})

	var l4 = new(List[int])
	l4.PushBackList(l2)
	checkList(t, l4, []int{1})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	var l List[int]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertBefore(1, new(Elem[int]))
	checkList(t, &l, []int{1, 2, 3})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	var l List[int]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertAfter(1, new(Elem[int]))
	checkList(t, &l, []int{1, 2, 3})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestMoveUnknownMark(t *testing.T) {
	var l1 List[int]
	e1 := l1.PushBack(1)

	var l2 List[int]
	e2 := l2.PushBack(2)

	l1.MoveAfter(e1, e2)
	checkList(t, &l1, []int{1})
	checkList(t, &l2, []int{2})

	l1.MoveBefore(e1, e2)
	checkList(t, &l1, []int{1})
	checkList(t, &l2, []int{2})
}
