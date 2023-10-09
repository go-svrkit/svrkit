// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package reflext

import (
	"image"
	"os"
	"reflect"
	"testing"
)

type TB struct {
	A string
	C image.Rectangle
	D []int16
	E map[string]float64
}

func (b *TB) GetC() *image.Rectangle {
	return &b.C
}

type TA struct {
	A int32
	B TB
	C []image.Point
	D map[float64]string
}

func (a *TA) GetB() *TB {
	return &a.B
}

func createTA() *TA {
	return &TA{
		B: TB{
			A: "hello",
			C: image.Rect(12, 34, 56, 78),
			D: []int16{12, 34, 56},
			E: map[string]float64{"hello": 12, "world": 34},
		},
		A: 1234,
		C: []image.Point{image.Pt(12, 34)},
		D: map[float64]string{3.14: "PI"},
	}
}

func TestEvalView(t *testing.T) {
	var obj = createTA()
	tests := []struct {
		expr         string
		shouldHasErr bool
		result       interface{}
	}{
		//{"C", false, obj.C},
		//{"B.A", false, obj.B.A},
		//{"B.C.Max.X", false, obj.B.C.Max.X},
		{"C[0].X", false, obj.C[0].X},
		{"D[3.14]", false, obj.D[3.14]},
		{`GetB().GetC().Min.X`, false, 12},
		{"C[10]", false, nil},
		{"D[10]", false, nil},
		{"B.D[10]", false, nil},
		{"B.E[10]", false, nil},
		{`D["KKK"]`, true, nil},
		{`B.E["KKK"]`, true, nil},
	}
	for _, tc := range tests {
		var evalCtx = NewEvalContext(obj)
		v, err := evalCtx.Eval(tc.expr)
		if tc.shouldHasErr {
			if err == nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
			continue
		}
		t.Logf("%s: %v", tc.expr, err)
		evalCtx.PrintNodes(os.Stdout)
		if tc.result == nil && !IsInterfaceNil(v) {
			t.Fatalf("%s: %v", tc.expr, err)
		} else if !reflect.DeepEqual(v, tc.result) {
			t.Fatalf("%s: %v", tc.expr, err)
		}
	}
}

func TestEvalSet(t *testing.T) {
	var obj = createTA()
	tests := []struct {
		expr         string
		val          interface{}
		shouldHasErr bool
	}{
		{"A", int32(5678), false},
		{"B.A", "hi", false},
		{"B.C.Max.X", 100, false},
		{"C[0].X", 54321, false},
		{"D[3.14]", "pi", false},
		{"D[1.68]", "ratio", false},
		{`B.E["hello"]`, 12345, false},
		{`GetB().E["hello"]`, 54321, false},
	}
	var ctx = NewEvalContext(obj)
	for _, tc := range tests {
		err := ctx.Set(tc.expr, tc.val)
		t.Logf("%s: %v", tc.expr, err)
		if tc.shouldHasErr {
			if err == nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
			continue
		}
		v, err := ctx.Eval(tc.expr)
		if err != nil {
			t.Fatalf("%s: %v", tc.expr, err)
		}
		if !reflect.DeepEqual(v, tc.val) {
			t.Fatalf("%s: %v", tc.expr, err)
		}
	}
}

func TestEvalRemove(t *testing.T) {
	var obj = createTA()
	tests := []struct {
		expr   string
		hasErr bool
	}{
		{"A", true},
		{"B.D[1]", false},
		{`GetB().E["hello"]`, false},
	}
	var ctx = NewEvalContext(obj)
	for _, tc := range tests {
		err := ctx.Delete(tc.expr)
		t.Logf("%s: %v", tc.expr, err)
		if tc.hasErr {
			if err == nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
		} else {
			if err != nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
		}
	}
}
