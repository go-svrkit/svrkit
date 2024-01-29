// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package treemap

type (
	KeyType     Comparable
	EntryAction func(key KeyType, val any)
)

// Comparable 丐版java.lang.Comparable
// 内部实现要符合结合律:  (a.CompareTo(b) > 0 && b.CompareTo(c) > 0) implies a.CompareTo(c) > 0
type Comparable interface {
	// CompareTo returns an integer comparing two Comparables.
	// a.CompareTo(b) < 0 implies a < b
	// a.CompareTo(b) > 0 implies a > b
	// a.CompareTo(b) == 0 implies a == b
	CompareTo(Comparable) int
}

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

type Entry struct {
	key                 KeyType
	value               any
	left, right, parent *Entry
	color               Color
}

func NewEntry(key KeyType, val any, parent *Entry) *Entry {
	return &Entry{
		key:    key,
		value:  val,
		parent: parent,
		color:  BLACK,
	}
}

func (e *Entry) GetKey() KeyType {
	return e.key
}

func (e *Entry) GetValue() any {
	return e.value
}

func (e *Entry) SetValue(val any) any {
	var old = e.value
	e.value = val
	return old
}

func (e *Entry) Equals(other *Entry) bool {
	if e == other {
		return true
	}
	return e.key == other.key && e.value == other.value
}

func key(e *Entry) KeyType {
	if e != nil {
		return e.key
	}
	return nil
}

func colorOf(p *Entry) Color {
	if p != nil {
		return p.color
	}
	return BLACK
}

func parentOf(p *Entry) *Entry {
	if p != nil {
		return p.parent
	}
	return nil
}

func setColor(p *Entry, color Color) {
	if p != nil {
		p.color = color
	}
}

func leftOf(p *Entry) *Entry {
	if p != nil {
		return p.left
	}
	return nil
}

func rightOf(p *Entry) *Entry {
	if p != nil {
		return p.right
	}
	return nil
}

// Returns the successor of the specified Entry, or null if no such.
func successor(t *Entry) *Entry {
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
func predecessor(t *Entry) *Entry {
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
