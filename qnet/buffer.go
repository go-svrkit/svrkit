// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"math"
	"unsafe"

	"gopkg.in/svrkit.v1/pool"
)

type Buffer struct {
	buf  *bytes.Buffer
	pool *pool.ObjectPool[Buffer]
}

func NewBuffer(b []byte, pool *pool.ObjectPool[Buffer]) *Buffer {
	return &Buffer{
		buf:  bytes.NewBuffer(b),
		pool: pool,
	}
}

func (b *Buffer) Free() {
	if b.pool != nil {
		b.pool.Put(b)
	}
}

func (b *Buffer) Reset() {
	if b.buf != nil {
		b.buf.Reset()
	}
}

func (b *Buffer) Bytes() []byte {
	if b.buf != nil {
		return b.buf.Bytes()
	}
	return nil
}

func (b *Buffer) linit() {
	if b.buf == nil {
		b.buf = new(bytes.Buffer)
	}
}

func (b *Buffer) WriteBool(v bool) {
	var c byte = 0
	if v {
		c = 1
	}
	b.linit()
	b.buf.WriteByte(c)
}

func (b *Buffer) WriteInt8(n int8) {
	b.linit()
	b.buf.WriteByte(byte(n))
}

func (b *Buffer) WriteUint8(n uint8) {
	b.linit()
	b.buf.WriteByte(n)
}

func (b *Buffer) WriteUint16(n uint16) {
	b.linit()
	var tmp [2]byte
	binary.LittleEndian.PutUint16(tmp[:], n)
	b.buf.Write(tmp[:])
}

func (b *Buffer) WriteInt16(n int16) {
	b.WriteUint16(uint16(n))
}

func (b *Buffer) WriteUint32(n uint32) {
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], n)
	b.linit()
	b.buf.Write(tmp[:])
}

func (b *Buffer) WriteInt32(n int32) {
	b.WriteUint32(uint32(n))
}

func (b *Buffer) WriteUint64(n uint64) {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], n)
	b.linit()
	b.buf.Write(tmp[:])
}

func (b *Buffer) WriteInt64(n int64) {
	b.WriteUint64(uint64(n))
}

func (b *Buffer) WriteFloat32(f float32) {
	var n = math.Float32bits(f)
	b.WriteUint32(n)
}

func (b *Buffer) WriteFloat64(f float64) {
	var n = math.Float64bits(f)
	b.WriteUint64(n)
}

func (b *Buffer) WriteString(s string) {
	b.linit()
	b.buf.WriteString(s)
}
func (b *Buffer) WriteBytes(buf []byte) {
	b.linit()
	b.buf.Write(buf)
}

func (b *Buffer) ReadBool() (bool, error) {
	var c, err = b.ReadUint8()
	return c != 0, err
}

func (b *Buffer) MustReadBool() bool {
	v, err := b.ReadBool()
	if err != nil {
		log.Panicf("MustReadBool: %v", err)
	}
	return v
}

func (b *Buffer) ReadUint8() (uint8, error) {
	return b.buf.ReadByte()
}

func (b *Buffer) MustReadUint8() uint8 {
	v, err := b.ReadUint8()
	if err != nil {
		log.Panicf("MustReadUint8: %v", err)
	}
	return v
}

func (b *Buffer) ReadInt8() (int8, error) {
	var c, err = b.ReadUint8()
	return int8(c), err
}

func (b *Buffer) MustReadInt8() int8 {
	v, err := b.ReadInt8()
	if err != nil {
		log.Panicf("MustReadInt8: %v", err)
	}
	return v
}

func (b *Buffer) ReadUint16() (uint16, error) {
	lo, err := b.ReadUint8()
	if err != nil {
		return 0, err
	}
	hi, err := b.ReadUint8()
	if err != nil {
		return 0, err
	}
	return uint16(lo) | uint16(hi)<<8, nil
}

func (b *Buffer) MustReadUint16() uint16 {
	v, err := b.ReadUint16()
	if err != nil {
		log.Panicf("MustReadUint16: %v", err)
	}
	return v
}

func (b *Buffer) ReadInt16() (int16, error) {
	var n, err = b.ReadUint16()
	return int16(n), err
}

func (b *Buffer) MustReadInt16() int16 {
	v, err := b.ReadInt16()
	if err != nil {
		log.Panicf("MustReadInt16: %v", err)
	}
	return v
}

func (b *Buffer) ReadUint32() (n uint32, err error) {
	var tmp [4]byte
	if _, err = b.buf.Read(tmp[:]); err != nil {
		return
	}
	n = binary.LittleEndian.Uint32(tmp[:])
	return
}

func (b *Buffer) MustReadUint32() uint32 {
	v, err := b.ReadUint32()
	if err != nil {
		log.Panicf("MustReadUint32: %v", err)
	}
	return v
}

func (b *Buffer) ReadInt32() (int32, error) {
	var n, err = b.ReadUint32()
	return int32(n), err
}

func (b *Buffer) MustReadInt32() int32 {
	v, err := b.ReadInt32()
	if err != nil {
		log.Panicf("MustReadInt32: %v", err)
	}
	return v
}

func (b *Buffer) ReadUint64() (n uint64, err error) {
	var tmp [8]byte
	if _, err = b.buf.Read(tmp[:]); err != nil {
		return
	}
	n = binary.LittleEndian.Uint64(tmp[:])
	return
}

func (b *Buffer) MustReadUint64() uint64 {
	v, err := b.ReadUint64()
	if err != nil {
		log.Panicf("MustReadUint64: %v", err)
	}
	return v
}

func (b *Buffer) ReadInt64() (int64, error) {
	var n, err = b.ReadUint64()
	return int64(n), err
}

func (b *Buffer) MustReadInt64() int64 {
	v, err := b.ReadInt64()
	if err != nil {
		log.Panicf("MustReadInt64: %v", err)
	}
	return v
}

func (b *Buffer) ReadFloat32() (float32, error) {
	var n, err = b.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(n), nil
}

func (b *Buffer) MustReadFloat32() float32 {
	v, err := b.ReadFloat32()
	if err != nil {
		log.Panicf("MustReadFloat32: %v", err)
	}
	return v
}

func (b *Buffer) ReadFloat64() (float64, error) {
	var n, err = b.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(n), nil
}

func (b *Buffer) MustReadFloat64() float64 {
	v, err := b.ReadFloat64()
	if err != nil {
		log.Panicf("MustReadFloat64: %v", err)
	}
	return v
}

func (b *Buffer) ReadNBytes(n int) (r []byte, err error) {
	var data = b.buf.Bytes()
	if n > 0 && len(data) >= n {
		r = make([]byte, n)
		_, err = b.buf.Read(r)
		return
	}
	return nil, io.EOF
}

func (b *Buffer) ReadNString(n int) (string, error) {
	buf, err := b.ReadNBytes(n)
	if err == nil {
		return *(*string)(unsafe.Pointer(&buf)), nil
	}
	return "", err
}

func (b *Buffer) PeekBool() (bool, error) {
	n, err := b.PeekUint8()
	return n > 0, err
}

func (b *Buffer) PeekUint8() (uint8, error) {
	var data = b.buf.Bytes()
	if len(data) >= 1 {
		return data[0], nil
	}
	return 0, io.EOF
}

func (b *Buffer) PeekInt8() (int8, error) {
	n, err := b.PeekUint8()
	return int8(n), err
}

func (b *Buffer) PeekUint16() (uint16, error) {
	var data = b.buf.Bytes()
	if len(data) >= 2 {
		n := binary.LittleEndian.Uint16(data[:2])
		return n, nil
	}
	return 0, io.EOF
}

func (b *Buffer) PeekInt16() (int16, error) {
	n, err := b.PeekUint16()
	return int16(n), err
}

func (b *Buffer) PeekUint32() (uint32, error) {
	var data = b.buf.Bytes()
	if len(data) >= 4 {
		n := binary.LittleEndian.Uint32(data[:4])
		return n, nil
	}
	return 0, io.EOF
}

func (b *Buffer) PeekInt32() (int32, error) {
	n, er := b.PeekUint32()
	return int32(n), er
}

func (b *Buffer) PeekUint64() (uint64, error) {
	var data = b.buf.Bytes()
	if len(data) >= 8 {
		n := binary.LittleEndian.Uint64(data[:8])
		return n, nil
	}
	return 0, io.EOF
}

func (b *Buffer) PeekInt64() (int64, error) {
	n, ok := b.PeekUint64()
	return int64(n), ok
}

func (b *Buffer) PeekFloat32() (float32, error) {
	n, err := b.PeekUint32()
	if err == nil {
		return math.Float32frombits(n), nil
	}
	return 0, err
}

func (b *Buffer) PeekFloat64() (float64, error) {
	n, err := b.PeekUint64()
	if err == nil {
		return math.Float64frombits(n), nil
	}
	return 0, err
}
