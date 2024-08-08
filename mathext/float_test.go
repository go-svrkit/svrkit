// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package mathext

import (
	"testing"
)

func TestIsAlmostEqualFloat32(t *testing.T) {
	var a float32 = 0.15 + 0.15
	var b float32 = 0.1 + 0.2
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
	if !IsAlmostEqualFloat32(b, a) {
		t.Fatalf("%v is not almost equal to %v", b, a)
	}
}

func TestIsAlmostEqualFloat64(t *testing.T) {
	var a float64 = 1.0
	var b = float64(RoundHalf(0.5001))
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
	if !IsAlmostEqualFloat64(b, a) {
		t.Fatalf("%v is not almost equal to %v", b, a)
	}
}
