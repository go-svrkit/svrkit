// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package mathext

import (
	"math/bits"
	"math/rand/v2"
)

// RandInt rand an integer in [min, max]
func RandInt(min, max int) int {
	if min > max {
		max, min = min, max
	}
	if min == max {
		return min
	}
	return rand.IntN(max-min+1) + min
}

// RandFloat rand a float number in [min, max]
func RandFloat(min, max float64) float64 {
	if min > max {
		max, min = min, max
	}
	if min == max {
		return min
	}
	return rand.Float64()*(max-min) + min
}

// RangePerm [min,max]区间内的随机数集合
func RangePerm(min, max int) []int {
	if min > max {
		max, min = min, max
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

// LCGRand 线性同余法的随机数生成器
// see https://en.wikipedia.org/wiki/Linear_congruential_generator

const (
	rngMax          = 1 << 63
	lcg32Multiplier = 214013 // Visual C++
	lcg32Increment  = 2531011
	lcg64Multiplier = 6364136223846793005 // musl
	lcg64Increment  = 1
)

type LCGRand32 uint32

func (g *LCGRand32) Seed(seed uint32) {
	*(*uint32)(g) = seed*lcg32Multiplier + lcg32Increment
}

func (g *LCGRand32) Uint32() uint32 {
	*(*uint32)(g) = uint32(*g)*lcg32Multiplier + lcg32Increment
	return uint32(uint32(*g)>>16) & 0x7fff
}

func (g *LCGRand32) Uint64() uint64 {
	return (uint64(g.Uint32()) << 32) | uint64(g.Uint32())
}

func (g *LCGRand32) Int63() int64 {
	return int64(g.Uint64() & (rngMax - 1))
}

type LCGRand64 uint64

func (g *LCGRand64) Seed(seed uint64) {
	*(*uint64)(g) = seed - 1
}

func (g *LCGRand64) Uint64() uint64 {
	*(*uint64)(g) = uint64(*g)*lcg32Multiplier + lcg32Increment
	return uint64(*g) >> 33
}

func (g *LCGRand64) Int63() int64 {
	return int64(g.Uint64() & (rngMax - 1))
}

// WyRand see https://github.com/wangyi-fudan/wyhash
type WyRand uint64

func _wymix(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}

func (r *WyRand) Uint64() uint64 {
	*r += WyRand(0xa0761d6478bd642f)
	return _wymix(uint64(*r), uint64(*r^WyRand(0xe7037ed1a0b428db)))
}

func (r *WyRand) Uint64n(n uint64) uint64 {
	return r.Uint64() % n
}

func (r *WyRand) Uint32() uint32 {
	return uint32(r.Uint64())
}

func (r *WyRand) Uint32n(n int) uint32 {
	// This is similar to Uint32() % n, but faster.
	// See https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
	return uint32(uint64(r.Uint32()) * uint64(n) >> 32)
}
