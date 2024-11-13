// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflext

import (
	"reflect"
	"unsafe"
)

// Everything below is taken from the Go 1.21 runtime package, and must stay in sync with it.

func FuncPCABI0(f interface{}) uintptr {
	words := (*[2]unsafe.Pointer)(unsafe.Pointer(&f))
	return *(*uintptr)(unsafe.Pointer(words[1]))
}

type GoSlice struct {
	Ptr unsafe.Pointer
	Len int
	Cap int
}

type GoString struct {
	Ptr unsafe.Pointer
	Len int
}

type GoMap struct {
	Count      int            // size of map
	Flags      uint8          //
	B          uint8          // log_2 of # of buckets
	NOverflow  uint16         // approximate number of overflow buckets
	Hash0      uint32         // hash seed
	Buckets    unsafe.Pointer // array of 2^B Buckets
	OldBuckets unsafe.Pointer
	Evacuate   uintptr
	Extra      unsafe.Pointer
}

type bmap struct {
	tophash [8]uint8
}

// GoEface empty interface
type GoEface struct {
	Typ  *GoType
	Data unsafe.Pointer
}

func PackEface(e *GoEface) (v interface{}) {
	*(*GoEface)(unsafe.Pointer(&v)) = *e
	return
}

func UnpackEface(v interface{}) GoEface {
	return *(*GoEface)(unsafe.Pointer(&v))
}

// GoIface none empty interface
type GoIface struct {
	Itab  *GoItab
	Value unsafe.Pointer
}

func getReflectTypeItab() *GoItab {
	v := reflect.TypeOf(struct{}{})
	return (*GoIface)(unsafe.Pointer(&v)).Itab
}

type RValue struct {
	Type *GoType
	Ptr  unsafe.Pointer
	Flag uintptr
}

// GoItab itable
type GoItab struct {
	inter *GoIfaceType
	_type *GoType
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

// GoType is the runtime representation of a Go type.
type GoType struct {
	Size_       uintptr
	PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash        uint32  // hash of type; avoids computation in hash tables
	TFlag       TFlag   // extra type information flags
	Align_      uint8   // alignment of variable with this type
	FieldAlign_ uint8   // alignment of struct field with this type
	Kind_       uint8   // enumeration for C
	Equal       func(unsafe.Pointer, unsafe.Pointer) bool
	GCData      *byte
	Str         int32
	PtrToThis   int32
}

var rtypeItab = getReflectTypeItab()

func PackReflectType(typ *GoType) (rtyp reflect.Type) {
	(*GoIface)(unsafe.Pointer(&rtyp)).Itab = rtypeItab
	(*GoIface)(unsafe.Pointer(&rtyp)).Value = unsafe.Pointer(typ)
	return
}

func UnpackReflectType(t reflect.Type) *GoType {
	return (*GoType)((*GoIface)(unsafe.Pointer(&t)).Value)
}

const (
	KindDirectIface = 1 << 5
	KindGCProg      = 1 << 6 // Type.gc points to GC program
	KindMask        = (1 << 5) - 1
)

// TFlag is used by a Type to signal what extra type information is
// available in the memory directly following the Type value.
type TFlag uint8

const (
	TFlagUncommon      TFlag = 1 << 0
	TFlagExtraStar     TFlag = 1 << 1
	TFlagNamed         TFlag = 1 << 2
	TFlagRegularMemory TFlag = 1 << 3
)

func (t *GoType) Kind() reflect.Kind { return reflect.Kind(t.Kind_ & KindMask) }

func (t *GoType) HasName() bool {
	return t.TFlag&TFlagNamed != 0
}

func (t *GoType) Pointers() bool { return t.PtrBytes != 0 }

// IfaceIndir reports whether t is stored indirectly in an interface value.
func (t *GoType) IfaceIndir() bool {
	return t.Kind_&KindDirectIface == 0
}

func (t *GoType) IsDirectIface() bool {
	return t.Kind_&KindDirectIface != 0
}

func (t *GoType) GcSlice(begin, end uintptr) []byte {
	return unsafe.Slice(t.GCData, int(end))[begin:]
}

// Method on non-interface type
type Method struct {
	Name int32 // name of method
	Mtyp int32 // method type (without receiver)
	Ifn  int32 // fn used in interface call (one-word receiver)
	Tfn  int32 // fn used for normal method call
}

type UncommonType struct {
	PkgPath int32  // import path; empty for built-in types like int, string
	Mcount  uint16 // number of methods
	Xcount  uint16 // number of exported methods
	Moff    uint32 // offset from this uncommontype to [mcount]Method
	_       uint32 // unused
}

// Imethod represents a method on an interface type
type Imethod struct {
	Name int32 // name of method
	Typ  int32 // .(*FuncType) underneath
}

// GoArrayType represents a fixed array type.
type GoArrayType struct {
	GoType
	Elem  *GoType // array element type
	Slice *GoType // slice type
	Len   uintptr
}

type ChanDir int

const (
	RecvDir    ChanDir = 1 << iota         // <-chan
	SendDir                                // chan<-
	BothDir            = RecvDir | SendDir // chan
	InvalidDir ChanDir = 0
)

// GoChanType represents a channel type
type GoChanType struct {
	GoType
	Elem *GoType
	Dir  ChanDir
}

type GoIfaceType struct {
	GoType
	PkgPath Name      // import path
	Methods []Imethod // sorted by hash
}

type GoMapType struct {
	GoType
	Key    *GoType
	Elem   *GoType
	Bucket *GoType // internal type representing a hash bucket
	// function for hashing keys (ptr to key, seed) -> hash
	Hasher     func(unsafe.Pointer, uintptr) uintptr
	KeySize    uint8  // size of key slot
	ValueSize  uint8  // size of elem slot
	BucketSize uint16 // size of bucket
	Flags      uint32
}

type GoSliceType struct {
	GoType
	Elem *GoType // slice element type
}

type GoFuncType struct {
	GoType
	InCount  uint16
	OutCount uint16 // top bit is set if last input parameter is ...
}

type GoPtrType struct {
	GoType
	Elem *GoType // pointer element (pointed at) type
}

type GoStructField struct {
	Name   Name    // name is always non-empty
	Typ    *GoType // type of field
	Offset uintptr // byte offset of field
}

type GoStructType struct {
	GoType
	PkgPath Name
	Fields  []GoStructField
}

type Name struct {
	Bytes *byte
}

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

// PcHeader holds data used by the pclntab lookups.
type PcHeader struct {
	Magic          uint32  // 0xFFFFFFF1
	pad1, pad2     uint8   // 0,0
	MinLC          uint8   // min instruction size
	PtrSize        uint8   // size of a ptr in bytes
	Nfunc          int     // number of functions in the module
	Nfiles         uint    // number of entries in the file tab
	TextStart      uintptr // base for function entry PC offsets in this module, equal to moduledata.text
	FuncnameOffset uintptr // offset to the funcnametab variable from pcHeader
	CuOffset       uintptr // offset to the cutab variable from pcHeader
	FiletabOffset  uintptr // offset to the filetab variable from pcHeader
	PctabOffset    uintptr // offset to the pctab variable from pcHeader
	PclnOffset     uintptr // offset to the pclntab variable from pcHeader
}

type Functab struct {
	Entryoff uint32 // relative to runtime.text
	Funcoff  uint32
}

type Textsect struct {
	Vaddr    uintptr // prelinked section vaddr
	End      uintptr // vaddr + section length
	Baseaddr uintptr // relocated section address
}

type PtabEntry struct {
	Name int32
	Typ  int32
}

type Modulehash struct {
	Modulename   string
	Linktimehash string
	Runtimehash  *string
}

type InitTask struct {
	State uint32 // 0 = uninitialized, 1 = in progress, 2 = done
	Nfns  uint32
	// followed by nfns pcs, uintptr sized, one per init function to run
}

type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}

type Moduledata struct {
	PcHeader     *PcHeader
	Funcnametab  []byte
	Cutab        []uint32
	Filetab      []byte
	Pctab        []byte
	Pclntable    []byte
	Ftab         []Functab
	Findfunctab  uintptr
	MinPc, MaxPc uintptr

	Text, Etext           uintptr
	Noptrdata, Enoptrdata uintptr
	Data, Edata           uintptr
	Bss, Ebss             uintptr
	Noptrbss, Enoptrbss   uintptr
	Covctrs, Ecovctrs     uintptr
	End, Gcdata, Gcbss    uintptr
	Types, Etypes         uintptr
	Rodata                uintptr
	Gofunc                uintptr // go.func.*

	Textsectmap []Textsect
	Typelinks   []int32 // offsets from types
	Itablinks   []*GoItab

	Ptab []PtabEntry

	Pluginpath string
	Pkghashes  []Modulehash

	// This slice records the initializing tasks that need to be
	// done to start up the program. It is built by the linker.
	Inittasks []*InitTask

	Modulename   string
	Modulehashes []Modulehash

	Hasmain uint8 // 1 if module contains the main function, 0 otherwise

	Gcdatamask, Gcbssmask bitvector

	Typemap map[int32]*GoType // offset to *_rtype in previous module

	Bad bool // module failed to load and should be ignored

	Next *Moduledata
}

func escapes(x any) {
	if dummy.b {
		dummy.x = x
	}
}

var dummy struct {
	b bool
	x any
}
