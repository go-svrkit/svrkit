// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package slice

import (
	"math/rand"
)

// 常用的数组操作，配合下面的包使用
// `golang.org/x/exp/slices`  --> std.slices
// `golang.org/x/exp/maps`	  --> std.maps

// InsertAt 把`v`插入到第`i`个位置
func InsertAt[E any](s []E, i int, v E) []E {
	if i >= 0 && i < len(s) {
		return append(s[:i], append([]E{v}, s[i:]...)...)
	}
	return append(s, v)
}

// RemoveAt 删除第`i`个元素，不保证原来元素的顺序
func RemoveAt[E any](s []E, i int) []E {
	var zero E
	if n := len(s); i >= 0 && i < n {
		s[i] = s[n-1]
		s[n-1] = zero
		return s[:n-1]
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
