// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package treemap

import (
	"gopkg.in/svrkit.v1/collections/util"
)

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
	left, right, parent *Entry[K, V]
	color               Color
	key                 K
	value               V
}

func NewEntry[K comparable, V any](key K, val V, parent *Entry[K, V]) *Entry[K, V] {
	return &Entry[K, V]{
		key:    key,
		value:  val,
		parent: parent,
		color:  BLACK,
	}
}

func (e *Entry[K, V]) GetKey() K {
	return e.key
}

func (e *Entry[K, V]) GetValue() V {
	return e.value
}

func (e *Entry[K, V]) SetValue(val V) V {
	var old = e.value
	e.value = val
	return old
}

//func (e *Entry[K, V]) Equals(other *Entry[K, V]) bool {
//	if e == other {
//		return true
//	}
//	return e.key == other.key && e.value == other.value
//}

func keyOf[K comparable, V any](e *Entry[K, V]) (K, bool) {
	if e != nil {
		return e.key, true
	}
	return util.ZeroOf[K](), false
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
