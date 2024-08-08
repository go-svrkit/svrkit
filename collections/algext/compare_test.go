// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package algext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedCmp(t *testing.T) {
	assert.Equal(t, 0, OrderedCmp(123, 123))
	assert.Equal(t, 1, OrderedCmp(132, 123))
	assert.Equal(t, -1, OrderedCmp(123, 132))
	assert.Equal(t, 0, OrderedCmp(3.14, 3.14))
	assert.Equal(t, 1, OrderedCmp(3.41, 3.14))
	assert.Equal(t, -1, OrderedCmp(3.14, 3.41))
	assert.Equal(t, 0, OrderedCmp("a", "a"))
	assert.Equal(t, -1, OrderedCmp("a", "b"))
	assert.Equal(t, 1, OrderedCmp("b", "a"))
}

func TestReverseOrderedCmp(t *testing.T) {
	assert.Equal(t, 0, ReverseOrderedCmp(123, 123))
	assert.Equal(t, -1, ReverseOrderedCmp(132, 123))
	assert.Equal(t, 1, ReverseOrderedCmp(123, 132))
	assert.Equal(t, 0, ReverseOrderedCmp(3.14, 3.14))
	assert.Equal(t, -1, ReverseOrderedCmp(3.41, 3.14))
	assert.Equal(t, 1, ReverseOrderedCmp(3.14, 3.41))
	assert.Equal(t, 0, ReverseOrderedCmp("a", "a"))
	assert.Equal(t, 1, ReverseOrderedCmp("a", "b"))
	assert.Equal(t, -1, ReverseOrderedCmp("b", "a"))
}

func TestBoolCmp(t *testing.T) {
	assert.Equal(t, 0, BoolCmp(true, true))
	assert.Equal(t, 0, BoolCmp(false, false))
	assert.Equal(t, -1, BoolCmp(false, true))
	assert.Equal(t, 1, BoolCmp(true, false))
}

func TestComplex64Cmp(t *testing.T) {
	assert.Equal(t, 0, Complex64Cmp(complex(1, 2), complex(1, 2)))
	assert.Equal(t, 1, Complex64Cmp(complex(2, 2), complex(1, 2)))
	assert.Equal(t, -1, Complex64Cmp(complex(1, 1), complex(1, 2)))
}

func TestComplex128Cmp(t *testing.T) {
	assert.Equal(t, 0, Complex128Cmp(complex(1, 2), complex(1, 2)))
	assert.Equal(t, 1, Complex128Cmp(complex(2, 2), complex(1, 2)))
	assert.Equal(t, -1, Complex128Cmp(complex(1, 1), complex(1, 2)))
}
