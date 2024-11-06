// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package conv

import (
	"bytes"
	"math"
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

type allInt struct {
	a int8
	b uint8
	c int16
	d uint16
	e int32
	f uint32
	g int64
	h uint64
}

func setAllInt(val *allInt, i int) {
	val.a = int8(i)
	val.b = uint8(i)
	val.c = int16(i)
	val.d = uint16(i)
	val.e = int32(i)
	val.f = uint32(i)
	val.g = int64(i)
	val.h = uint64(i)
}

func assertAllInt[T Integer](t *testing.T, expect T, val *allInt) {
	assert.Equal(t, expect, ConvTo[T](val.a))
	assert.Equal(t, expect, ConvTo[T](val.b))
	assert.Equal(t, expect, ConvTo[T](val.c))
	assert.Equal(t, expect, ConvTo[T](val.d))
	assert.Equal(t, expect, ConvTo[T](val.e))
	assert.Equal(t, expect, ConvTo[T](val.f))
	assert.Equal(t, expect, ConvTo[T](val.g))
	assert.Equal(t, expect, ConvTo[T](val.h))
}

func TestConvTo(t *testing.T) {
	var all = new(allInt)
	setAllInt(all, 123)
	assertAllInt(t, int8(123), all)
	assertAllInt(t, uint8(123), all)
	assertAllInt(t, int16(123), all)
	assertAllInt(t, uint16(123), all)
	assertAllInt(t, int32(123), all)
	assertAllInt(t, uint32(123), all)
	assertAllInt(t, int(123), all)
	assertAllInt(t, uint(123), all)
	assertAllInt(t, int64(123), all)
	assertAllInt(t, uint64(123), all)

	assert.Equal(t, int64(math.MinInt32), ConvTo[int64](math.MinInt32))

	assert.Equal(t, int16(1), ConvTo[int16](true))
	assert.Equal(t, int16(0), ConvTo[int16](false))

	assert.Equal(t, int32(-3), ConvTo[int32](-3.14))
	assert.Equal(t, int32(6), ConvTo[int32](6.18))

	assert.Equal(t, int32(123456789), ConvTo[int32]("123456789"))
	assert.Equal(t, int32(-123456789), ConvTo[int32]("-123456789"))
}
