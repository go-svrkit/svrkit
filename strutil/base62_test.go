// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package strutil

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBase62String(t *testing.T) {
	assert.Equal(t, "a", EncodeBase62String(0))
	assert.NotEqual(t, "a", EncodeBase62String(123456789))
	assert.NotEqual(t, "", EncodeBase62String(123456789))
	assert.NotEqual(t, "", EncodeBase62String(9223372036854775807))
}

func TestDecodeBase62String(t *testing.T) {
	assert.Equal(t, int64(0), DecodeBase62String("a"))
	var nn = []int64{0, math.MaxInt16, 1234567890, math.MaxInt32, math.MaxInt64}
	for _, n := range nn {
		var shorten = EncodeBase62String(n)
		var got = DecodeBase62String(shorten)
		assert.Equal(t, n, got)
	}
}

func BenchmarkEncodeBase62String(b *testing.B) {
	b.StopTimer()
	var id = rand.Int63()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		EncodeBase62String(id)
	}
}

func BenchmarkDecodeBase62String(b *testing.B) {
	b.StopTimer()
	var id = rand.Int63()
	var shorten = EncodeBase62String(id)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DecodeBase62String(shorten)
	}
}
