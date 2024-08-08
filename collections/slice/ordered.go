// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package slice

import (
	"cmp"
	"slices"
	"sort"
)

// 有序数组实现的集合，用于小量数据的场合

// OrderedIndexOf `elem`在有序数组`set`中的索引，如果不存在返回-1
func OrderedIndexOf[T cmp.Ordered](set []T, elem T) int {
	var i = sort.Search(len(set), func(i int) bool {
		return set[i] >= elem
	})
	if i < len(set) && set[i] == elem {
		return i
	}
	return -1
}

// OrderedContains `elem`是否在有序数组`set`中
func OrderedContains[T cmp.Ordered](set []T, elem T) bool {
	return OrderedIndexOf(set, elem) >= 0
}

// PutIfAbsent 插入`n`到有序数组
func PutIfAbsent[T cmp.Ordered](set []T, elem T) []T {
	var i = sort.Search(len(set), func(i int) bool {
		return set[i] >= elem
	})
	if i < len(set) && set[i] == elem {
		return set
	}
	return InsertAt(set, i, elem)
}

// OrderedDelete 把`n`从有序数组中删除
func OrderedDelete[T cmp.Ordered](set []T, elem T) []T {
	var i = OrderedIndexOf(set, elem)
	if i >= 0 {
		return slices.Delete(set, i, i+1)
	}
	return set
}

// OrderedUnion 有序数组并集, A ∪ B
func OrderedUnion[T cmp.Ordered](a, b []T) []T {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	var c = make([]T, 0)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			c = append(c, a[i])
			i++
		} else if a[i] > b[j] {
			c = append(c, b[j])
			j++
		} else {
			c = append(c, a[i])
			i++
			j++
		}
	}
	for i < len(a) {
		c = append(c, a[i])
		i++
	}
	for j < len(b) {
		c = append(c, b[j])
		j++
	}
	return c
}

// OrderedIntersect 有序数组交集, A ∩ B
func OrderedIntersect[T cmp.Ordered](a, b []T) []T {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	var c = make([]T, 0)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			c = append(c, a[i])
			i++
			j++
		}
	}
	return c
}
