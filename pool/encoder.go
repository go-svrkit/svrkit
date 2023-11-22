package pool

import (
	"github.com/vmihailenco/msgpack/v5"
)

var (
	encoderPool = NewObjectPoolWith[msgpack.Encoder](createMPEncoder)
	decoderPool = NewObjectPoolWith[msgpack.Decoder](createMPDecoder)
)

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
