// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package helper

import (
	"unsafe"
)

// BytesAsStr 对[]byte的修改会影响到返回的string
func BytesAsStr(b []byte) string {
	//return *(*string)(unsafe.Pointer(&b))
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StrAsBytes returns the bytes backing a string
// see https://pkg.go.dev/unsafe#Pointer rule(6)
func StrAsBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))

	// We need to declare a real byte slice so internally the compiler
	// knows to use an unsafe.Pointer to keep track of the underlying memory so that
	// once the slice's array pointer is updated with the pointer to the string's
	// underlying bytes, the compiler won't prematurely GC the memory when the string
	// goes out of scope.
	//var b []byte
	//var hdr = (*reflect.SliceHeader)(unsafe.Pointer(&b))
	//
	//// This satisfies unsafe.Pointer rule 5, makes sure that even if GC relocates the string's
	//// underlying memory after this assignment, the corresponding unsafe.Pointer in the internal
	//// slice struct will be updated accordingly to reflect the memory relocation.
	//hdr.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	//
	//// It is important that we access s after we assign the Data
	//// pointer of the string header to the Data pointer of the slice header to
	//// make sure the string (and the underlying bytes backing the string) don't get
	//// GC'ed before the assignment happens.
	//hdr.Len = len(s)
	//hdr.Cap = len(s)
	//
	//return b
}
