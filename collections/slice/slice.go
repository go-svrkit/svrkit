// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package slice

import (
	"cmp"
	"math/rand"
	"slices"
	"sort"
)

// 提供一些常用的slice操作函数
// See also:
// 	https://pkg.go.dev/slices
// 	https://ueokande.github.io/go-slice-tricks/

// InsertAt 把`v`插入到第`i`个位置
func InsertAt[E any](s []E, i int, v E) []E {
	if i >= 0 && i < len(s) {
		return append(s[:i], append([]E{v}, s[i:]...)...)
	}
	return append(s, v)
}

// RemoveAt 删除第`i`个元素，不保证原来元素的顺序
func RemoveAt[E any](s []E, i int) []E {
	if n := len(s); i >= 0 && i < n {
		var zero E
		s[i] = s[n-1]
		s[n-1] = zero //  GC friendly
		return s[:n-1]
	}
	return s
}

// RemoveFirst 删除第一个查询到的元素
func RemoveFirst[E comparable](s []E, elem E) []E {
	for i, v := range s {
		if v == elem {
			return slices.Delete(s, i, i+1)
		}
	}
	return s
}

func Shuffle[E any](s []E) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

func Shrink[E any](s []E) []E {
	if len(s) == 0 {
		return nil
	}
	if len(s) == cap(s) {
		return s
	}
	var a = make([]E, len(s))
	copy(a, s)
	return a
}

// SortAndRemoveDup 去重并排序
func SortAndRemoveDup[E cmp.Ordered](s []E) []E {
	if len(s) <= 1 {
		return s
	}
	slices.Sort(s)
	s = slices.Compact(s)
	return s
}

// IsAllZeroElem 是否数组的所有元素都为0
func IsAllZeroElem[E cmp.Ordered | ~bool](s []E) bool {
	var zero E
	for i := 0; i < len(s); i++ {
		if s[i] != zero {
			return false
		}
	}
	return true
}

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

// OrderedPutIfAbsent 插入`n`到有序数组
func OrderedPutIfAbsent[T cmp.Ordered](set []T, elem T) []T {
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
