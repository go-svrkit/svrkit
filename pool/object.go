// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"bytes"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
)

func createObjPool[T any]() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return new(T)
		},
	}
}

func createObjPoolBy[T any](creator func() *T) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return creator()
		},
	}
}

type ObjectPool[T any] struct {
	pool *sync.Pool
}

func NewObjectPool[T any]() *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: createObjPool[T](),
	}
}

func NewObjectPoolWith[T any](creator func() *T) *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: createObjPoolBy[T](creator),
	}
}

func (a *ObjectPool[T]) Alloc() *T {
	return a.pool.Get().(*T)
}

func (a *ObjectPool[T]) Free(p *T) {
	a.pool.Put(p)
}

var (
	bufferPool = NewObjectPool[bytes.Buffer]()
	readerPool = NewObjectPool[bytes.Reader]()

	encoderPool = NewObjectPoolWith[msgpack.Encoder](createMPEncoder)
	decoderPool = NewObjectPoolWith[msgpack.Decoder](createMPDecoder)
)

func AllocBytesBuffer() *bytes.Buffer {
	return bufferPool.Alloc()
}

func FreeBytesBuffer(buf *bytes.Buffer) {
	// 太大的buffer应该直接交给GC，不要再回收了
	// See https://github.com/golang/go/issues/23199
	if buf.Cap() > 64<<10 {
		return
	}
	buf.Reset()
	bufferPool.Free(buf)
}

func AllocBytesReader() *bytes.Reader {
	return readerPool.Alloc()
}

func FreeBytesReader(rd *bytes.Reader) {
	rd.Reset(nil)
	readerPool.Free(rd)
}

func createMPEncoder() *msgpack.Encoder {
	var enc = msgpack.NewEncoder(nil)
	enc.SetCustomStructTag("json")
	enc.SetOmitEmpty(true)
	return enc
}

func createMPDecoder() *msgpack.Decoder {
	var dec = msgpack.NewDecoder(nil)
	dec.SetCustomStructTag("json")
	return dec
}

func AllocMPEncoder() *msgpack.Encoder {
	return encoderPool.Alloc()
}

func FreeMPEncoder(enc *msgpack.Encoder) {
	enc.Reset(nil)
	encoderPool.Free(enc)
}

func AllocMPDecoder() *msgpack.Decoder {
	return decoderPool.Alloc()
}

func FreeMPDecoder(dec *msgpack.Decoder) {
	dec.Reset(nil)
	decoderPool.Free(dec)
}
