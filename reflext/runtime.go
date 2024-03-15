// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflext

import (
	"unsafe"
)

type flag uintptr

const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5
	flagEmbedRO     flag = 1 << 6
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
	flagMethod      flag = 1 << 9
	flagMethodShift      = 10
	flagRO          flag = flagStickyRO | flagEmbedRO
)

func (f flag) kind() Kind {
	return Kind(f & flagKindMask)
}

func (f flag) ro() flag {
	if f&flagRO != 0 {
		return flagStickyRO
	}
	return 0
}

// Value is the reflection interface to a Go value.
type Value struct {
	// typ_ holds the type of the value represented by a Value.
	// Access using the typ method to avoid escape of v.
	Typ_ *Type

	// Pointer-valued data or, if flagIndir is set, pointer to data.
	// Valid when either flagIndir is set or typ.pointers() is true.
	Ptr unsafe.Pointer

	// flag holds metadata about the value.
	Flag flag
}

// Before Go 1.21, ValueOf always escapes and a Value's content
// is always heap allocated.

// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i. ValueOf(nil) returns the zero Value.
func ValueOf(i any) Value {
	if i == nil {
		return Value{}
	}
	return UnpackEface(i)
}

// UnpackEface converts the empty interface i to a Value.
func UnpackEface(i any) Value {
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	// NOTE: don't read e.word until we know whether it is really a pointer or not.
	t := e.Typ
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())
	if t.IfaceIndir() {
		f |= flagIndir
	}
	return Value{t, e.Word, f}
}

// PackEface converts v to the empty interface.
func PackEface(v Value) any {
	t := v.Typ_
	var i any
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	// First, fill in the data portion of the interface.
	switch {
	case t.IfaceIndir():
		if v.Flag&flagIndir == 0 {
			panic("bad indir")
		}
		// Value is indirect, and so is the interface we're making.
		ptr := v.Ptr
		if v.Flag&flagAddr != 0 {
			// TODO: pass safe boolean from valueInterface so
			// we don't need to copy if safe==true?
			c := unsafe_New(t)
			typedmemmove(t, c, ptr)
			ptr = c
		}
		e.Word = ptr
	case v.Flag&flagIndir != 0:
		// Value is indirect, but interface is direct. We need
		// to load the data at v.ptr into the interface data word.
		e.Word = *(*unsafe.Pointer)(v.Ptr)
	default:
		// Value is direct, and so is the interface.
		e.Word = v.Ptr
	}
	// Now, fill in the type portion. We're very careful here not
	// to have any operation between the e.word and e.typ assignments
	// that would let the garbage collector observe the partially-built
	// interface value.
	e.Typ = t
	return i
}

// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i any) *Type {
	eface := *(*EmptyInterface)(unsafe.Pointer(&i))
	// Noescape so this doesn't make i to escape. See the comment
	// at Value.typ for why this is safe.
	return (*Type)(noescape(unsafe.Pointer(eface.Typ)))
}

// EmptyInterface is the header for an interface{} value.
type EmptyInterface struct {
	Typ  *Type
	Word unsafe.Pointer
}

// Dummy annotation marking that the value x escapes,
// for use in cases where the reflect code is so clever that
// the compiler cannot follow.
func escapes(x any) {
	if dummy.b {
		dummy.x = x
	}
}

var dummy struct {
	b bool
	x any
}

//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

// typelinks is implemented in package runtime.
// It returns a slice of the sections in each module,
// and a slice of *rtype offsets in each module.
//
//go:linkname typelinks reflect.typelinks
func typelinks() (sections []unsafe.Pointer, offset [][]int32)

// resolveTypeOff resolves an *rtype offset from a base type.
// The (*rtype).typeOff method is a convenience wrapper for this function.
// Implemented in the runtime package.
//
//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

// memmove copies size bytes to dst from src. No write barriers are used.
//
//go:linkname memmove reflect.memmove
func memmove(dst, src unsafe.Pointer, size uintptr)

// typedmemmove copies a value of type t to dst from src.
//
//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(t *Type, dst, src unsafe.Pointer)

// typedmemclr zeros the value at ptr of type t.
//
//go:linkname typedmemclr reflect.typedmemclr
func typedmemclr(t *Type, ptr unsafe.Pointer)

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(*Type) unsafe.Pointer

//go:linkname unsafe_NewArray reflect.unsafe_NewArray
func unsafe_NewArray(*Type, int) unsafe.Pointer
