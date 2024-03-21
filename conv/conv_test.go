// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package conv

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randBytes(length int) []byte {
	if length <= 0 {
		return nil
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		ch := uint8(rand.Int31() % 0xFF)
		result[i] = ch
	}
	return result
}

func TestBytesAsString(t *testing.T) {
	var rawbytes = randBytes(1024)
	var s = BytesAsStr(rawbytes)
	assert.Equal(t, string(rawbytes), s)
}

func TestStringAsBytes(t *testing.T) {
	var text = string(randBytes(1024))
	var b = StrAsBytes(text)
	assert.True(t, bytes.Equal([]byte(text), b))
}

func BenchmarkBytesToString(b *testing.B) {
	b.StopTimer()
	var rawbytes = randBytes(2048)
	b.StartTimer()
	var text string
	for i := 0; i < 100000; i++ {
		text = string(rawbytes)
	}
	text = text[:0]
}

func BenchmarkBytesAsString(b *testing.B) {
	b.StopTimer()
	var rawbytes = randBytes(2048)
	b.StartTimer()
	var text string
	for i := 0; i < 100000; i++ {
		text = BytesAsStr(rawbytes)
	}
	text = text[:0]
}

func BenchmarkStringToBytes(b *testing.B) {
	b.StopTimer()
	var text = string(randBytes(2048))
	b.StartTimer()
	var rawbytes []byte
	for i := 0; i < 100000; i++ {
		rawbytes = []byte(text)
	}
	rawbytes = rawbytes[:0]
}

func BenchmarkStringAsBytes(b *testing.B) {
	b.StopTimer()
	var text = string(randBytes(2048))
	b.StartTimer()
	var rawbytes []byte
	for i := 0; i < 100000; i++ {
		rawbytes = StrAsBytes(text)
	}
	rawbytes = rawbytes[:0]
}

func TestBoolToInt(t *testing.T) {
	var n8 = BoolTo[uint8](true)
	assert.Equal(t, uint8(1), n8)
	var n = BoolTo[int](true)
	assert.Equal(t, int(1), n)
	n = BoolTo[int](false)
	assert.Equal(t, int(0), n)
}

func TestIntToBool(t *testing.T) {
	assert.Equal(t, true, IntToBool(1))
	assert.Equal(t, true, IntToBool(int64(2)))
	assert.Equal(t, true, IntToBool(int32(3)))
	assert.Equal(t, true, IntToBool(int16(4)))
	assert.Equal(t, true, IntToBool(int8(5)))
	assert.Equal(t, false, IntToBool(0))
	assert.Equal(t, false, IntToBool(int32(0)))
	assert.Equal(t, false, IntToBool(int16(0)))
	assert.Equal(t, false, IntToBool(int8(0)))
}
