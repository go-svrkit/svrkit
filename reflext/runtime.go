// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflext

import (
	"unsafe"
)

// Everything below is taken from the Go 1.21 runtime package, and must stay in sync with it.

//go:linkname firstmoduledata runtime.firstmoduledata
var firstmoduledata Moduledata

func GetFirstModuleData() *Moduledata {
	return &firstmoduledata
}

//go:linkname memmove reflect.memmove
func memmove(dst, src unsafe.Pointer, size uintptr)

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(t unsafe.Pointer, dst, src unsafe.Pointer)

//go:linkname typedmemclr reflect.typedmemclr
func typedmemclr(t unsafe.Pointer, ptr unsafe.Pointer)

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(unsafe.Pointer) unsafe.Pointer

//go:linkname unsafe_NewArray reflect.unsafe_NewArray
func unsafe_NewArray(unsafe.Pointer, int) unsafe.Pointer

//go:noescape
//go:linkname GetItab runtime.getitab
func GetItab(inter unsafe.Pointer, typ *GoType, canfail bool) *GoItab

//go:linkname typelinks reflect.typelinks
func typelinks() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer
