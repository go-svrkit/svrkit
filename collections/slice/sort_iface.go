// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package slice

import (
	"sort"
)

// hand-made generics
// 常用基本类型的sort.Interface wrapper

type (
	Int8Slice    []int8
	Uint8Slice   []uint8
	Int16Slice   []int16
	Uint16Slice  []uint16
	Int32Slice   []int32
	Uint32Slice  []uint32
	UintSlice    []uint
	Int64Slice   []int64
	Uint64Slice  []uint64
	Float32Slice []float32

	IntSlice     = sort.IntSlice
	Float64Slice = sort.Float64Slice
	StringSlice  = sort.StringSlice
)

func (x Int8Slice) Len() int           { return len(x) }
func (x Int8Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int8Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x Uint8Slice) Len() int           { return len(x) }
func (x Uint8Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint8Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x Int16Slice) Len() int           { return len(x) }
func (x Int16Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int16Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x Uint16Slice) Len() int           { return len(x) }
func (x Uint16Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint16Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x Int32Slice) Len() int           { return len(x) }
func (x Int32Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x Uint32Slice) Len() int           { return len(x) }
func (x Uint32Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x UintSlice) Len() int           { return len(x) }
func (x UintSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x UintSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x Int64Slice) Len() int           { return len(x) }
func (x Int64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x Uint64Slice) Len() int           { return len(x) }
func (x Uint64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (x Float32Slice) Len() int           { return len(x) }
func (x Float32Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Float32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
