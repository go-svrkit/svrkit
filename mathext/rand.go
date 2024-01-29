// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package mathext

import (
	"math/rand"
)

// LCG 线性同余法的随机数生成器
// see https://en.wikipedia.org/wiki/Linear_congruential_generator
type LCG struct {
	seed uint32
}

func (g *LCG) Seed(seed uint32) {
	g.seed = seed*214013 + 2531011
}

func (g *LCG) Rand() uint32 {
	g.seed = g.seed*214013 + 2531011
	var r = uint32(g.seed>>16) & 0x7fff
	return r
}

// RandInt rand an integer in [min, max]
func RandInt(min, max int) int {
	if min > max {
		panic("RandInt,min greater than max")
	}
	if min == max {
		return min
	}
	return rand.Intn(max-min+1) + min
}

// RandFloat rand a float number in [min, max]
func RandFloat(min, max float64) float64 {
	if min > max {
		panic("RandFloat: min greater than max")
	}
	if min == max {
		return min
	}
	return rand.Float64()*(max-min) + min
}

// RangePerm [min,max]区间内的随机数集合
func RangePerm(min, max int) []int {
	if min > max {
		panic("RangePerm: min greater than max")
	}
	if min == max {
		return []int{min}
	}
	list := rand.Perm(max - min + 1)
	for i := 0; i < len(list); i++ {
		list[i] += min
	}
	return list
}
