// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"math/rand"
	"testing"
)

func TestBase62Encode(t *testing.T) {
	for i := 0; i < 10000; i++ {
		var id = rand.Int63()
		var shorten = EncodeBase62String(id)
		var n = DecodeBase62String(shorten)
		if n != id {
			t.Fatalf("base62 not equal: %d != %d, %s", id, n, shorten)
		}
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
