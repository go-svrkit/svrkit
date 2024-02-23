package qnet

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer_WriteRead(t *testing.T) {
	var b = NewBuffer(nil, nil)
	b.WriteBool(true)
	assert.True(t, b.MustReadBool())

	b.WriteInt8(math.MinInt8)
	assert.Equal(t, b.MustReadInt8(), int8(math.MinInt8))
	b.WriteUint8(math.MaxUint8)
	assert.Equal(t, b.MustReadUint8(), uint8(math.MaxUint8))

	b.WriteUint16(math.MaxUint16)
	assert.Equal(t, b.MustReadUint16(), uint16(math.MaxUint16))
	b.WriteInt16(math.MinInt16)
	assert.Equal(t, b.MustReadInt16(), int16(math.MinInt16))

	b.WriteUint32(math.MaxUint32)
	assert.Equal(t, b.MustReadUint32(), uint32(math.MaxUint32))
	b.WriteInt32(math.MinInt32)
	assert.Equal(t, b.MustReadInt32(), int32(math.MinInt32))

	b.WriteUint64(math.MaxUint64)
	assert.Equal(t, b.MustReadUint64(), uint64(math.MaxUint64))
	b.WriteInt64(math.MinInt64)
	assert.Equal(t, b.MustReadInt64(), int64(math.MinInt64))

	b.WriteFloat32(math.MaxFloat32)
	assert.Equal(t, b.MustReadFloat32(), float32(math.MaxFloat32))
	b.WriteFloat64(math.MaxFloat64)
	assert.Equal(t, b.MustReadFloat64(), math.MaxFloat64)

	b.WriteString("hello")
	s, err := b.ReadNString(5)
	assert.Nil(t, err)
	assert.Equal(t, s, "hello")

	b.WriteBool(true)
	b.WriteUint8(math.MaxUint8)
	b.WriteUint16(math.MaxUint16)
	b.WriteUint32(math.MaxUint32)
	b.WriteUint64(math.MaxUint64)
	b.WriteString("world")

	assert.Equal(t, b.MustReadBool(), true)
	assert.Equal(t, b.MustReadUint8(), uint8(math.MaxUint8))
	assert.Equal(t, b.MustReadUint16(), uint16(math.MaxUint16))
	assert.Equal(t, b.MustReadUint32(), uint32(math.MaxUint32))
	assert.Equal(t, b.MustReadUint64(), uint64(math.MaxUint64))
	s, err = b.ReadNString(5)
	assert.Nil(t, err)
	assert.Equal(t, s, "world")
}

func TestBuffer_Peek(t *testing.T) {
	var b = NewBuffer(nil, nil)
	b.WriteBool(true)
	b.WriteUint8(math.MaxUint8)
	b.WriteUint16(math.MaxUint16)
	b.WriteUint32(math.MaxUint32)
	b.WriteUint64(math.MaxUint64)
	b.WriteString("world")

	{
		v, er := b.PeekBool()
		assert.Nil(t, er)
		assert.True(t, v)
		assert.Equal(t, b.MustReadBool(), true)
	}
	{
		v, er := b.PeekUint8()
		assert.Nil(t, er)
		assert.Equal(t, v, uint8(math.MaxUint8))
		assert.Equal(t, b.MustReadUint8(), uint8(math.MaxUint8))
	}
	{
		v, er := b.PeekUint16()
		assert.Nil(t, er)
		assert.Equal(t, v, uint16(math.MaxUint16))
		assert.Equal(t, b.MustReadUint16(), uint16(math.MaxUint16))
	}
	{
		v, er := b.PeekUint32()
		assert.Nil(t, er)
		assert.Equal(t, v, uint32(math.MaxUint32))
		assert.Equal(t, b.MustReadUint32(), uint32(math.MaxUint32))
	}
	{
		v, er := b.PeekUint64()
		assert.Nil(t, er)
		assert.Equal(t, v, uint64(math.MaxUint64))
		assert.Equal(t, b.MustReadUint64(), uint64(math.MaxUint64))
	}
}
