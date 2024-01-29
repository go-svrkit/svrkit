// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package slice

import (
	"math/rand"
)

// 常用的数组操作，配合下面的包使用
// `golang.org/x/exp/slices`  --> std.slices
// `golang.org/x/exp/maps`	  --> std.maps

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
	var r = make([]E, len(s))
	copy(r, s)
	return r
}
