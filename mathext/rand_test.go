// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package mathext

import (
	"testing"
)

func TestRandInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		v := RandInt(0, 1000)
		if v < 0 {
			t.Fatalf("%v < 0", v)
		}
		if v > 1000 {
			t.Fatalf("%v > 1000", v)
		}
	}
}

func TestRandFloat(t *testing.T) {
	for i := 0; i < 1000; i++ {
		v := RandFloat(0, 1.0)
		if v < 0 {
			t.Fatalf("%v < 0.0", v)
		}
		if v > 1.0 {
			t.Fatalf("%v > 1.0", v)
		}
	}
}

func TestRangePerm(t *testing.T) {
	var appeared = make(map[int]bool, 1000)
	var list = RangePerm(0, 1000)
	for _, v := range list {
		if v < 0 {
			t.Fatalf("%v < 0", v)
		}
		if v > 1000 {
			t.Fatalf("%v > 1000", v)
		}
		if _, found := appeared[v]; found {
			t.Fatalf("duplicate %v", v)
		}
		appeared[v] = true
	}
}

func TestLazyLCGRand32(t *testing.T) {
	var rng LCGRand32
	for i := 0; i < 1000; i++ {
		rng.Uint32()
	}
}

func TestSetLazyLCGSeed32(t *testing.T) {
	var rng LCGRand32
	rng.Seed(1234567890)
	for i := 0; i < 1000; i++ {
		rng.Uint32()
	}
}

func TestLazyLCGRand64(t *testing.T) {
	var rng LCGRand64
	for i := 0; i < 1000; i++ {
		rng.Uint64()
	}
}

func TestSetLazyLCGSeed64(t *testing.T) {
	var rng LCGRand64
	rng.Seed(1234567890)
	for i := 0; i < 1000; i++ {
		rng.Uint64()
	}
}
