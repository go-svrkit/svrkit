// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package qnet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/codec/testdata"
	"gopkg.in/svrkit.v1/strutil"
)

func TestMsgFlag_Has(t *testing.T) {
	var flag MsgFlag
	assert.False(t, flag.Has(FlagCompress))
	assert.False(t, flag.Has(FlagEncrypt))
	flag |= FlagCompress
	assert.True(t, flag.Has(FlagCompress))
	assert.False(t, flag.Has(FlagEncrypt))
	flag |= FlagEncrypt
	assert.True(t, flag.Has(FlagEncrypt))
}

func TestMsgFlag_Clear(t *testing.T) {
	var flag = FlagCompress | FlagEncrypt
	assert.True(t, flag.Has(FlagCompress))
	assert.True(t, flag.Has(FlagEncrypt))
	flag = flag.Clear(FlagEncrypt)
	assert.True(t, flag.Has(FlagCompress))
	assert.False(t, flag.Has(FlagEncrypt))
	flag = flag.Clear(FlagCompress)
	assert.False(t, flag.Has(FlagCompress))
}

func TestNetV1Header_Pack(t *testing.T) {
	var head = NewNetV1Header()
	head.Pack(12, 0xF, 123, 1234)
	head.SetCRC(12345678)
	assert.Equal(t, head.Len(), uint32(12))
	assert.Equal(t, head.Seq(), uint32(123))
	assert.Equal(t, head.Command(), uint32(1234))
	assert.Equal(t, head.Flag(), MsgFlag(0xF))
	assert.Equal(t, head.CRC(), uint32(12345678))
}

func TestCompress(t *testing.T) {
	tests := []struct {
		input string
	}{
		{""},
		{"hell world"},
		{"aaabbbcccdddeeefffggg"},
		{"a quick brown fox jumps over the lazy dog"},
		{"It was the best of times, it was the worst of times, it was the age of wisdom, it was the age of foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of Darkness, it was the spring of hope, it was the winter of despair, we had everything before us, we had nothing before us, we were all going direct to Heaven, we were all going direct the other way—in short, the period was so far like the present period, that some of its noisiest authorities insisted on its being received, for good or for evil, in the superlative degree of comparison only."},
	}
	for i, tc := range tests {
		var encoded bytes.Buffer
		if err := compress([]byte(tc.input), &encoded); err != nil {
			t.Fatalf("compress: %v", err)
		}
		var decoded bytes.Buffer
		if err := uncompress(encoded.Bytes(), &decoded); err != nil {
			t.Fatalf("uncompress: %v", err)
		}
		var out = string(decoded.Bytes())
		if out != tc.input {
			t.Logf("case %d: compress %s -> %s", i+1, tc.input, out)
		}
	}
}

func TestDecodeMsgFrom(t *testing.T) {
	var buf bytes.Buffer
	var msg = CreateNetMessageWith(&testdata.BuildReq{PosX: 11, PosY: 22})
	assert.Nil(t, EncodeMsgTo(msg, nil, &buf))
	var netMsg = AllocNetMessage()
	assert.Nil(t, DecodeMsgFrom(&buf, MaxPayloadSize, nil, netMsg))
	assert.True(t, len(netMsg.Data) > 0)
	var req testdata.BuildReq
	assert.Nil(t, netMsg.DecodeTo(&req))
	assert.Equal(t, req.PosX, int32(11))
	assert.Equal(t, req.PosY, int32(22))
}

func isMsgEqual(t *testing.T, a, b *NetMessage) bool {
	if !(a.Command == b.Command && a.Seq == b.Seq) {
		return false
	}
	if err := a.Encode(); err != nil {
		t.Fatalf("encode %v", err)
	}
	data1, data2 := a.Data, b.Data
	if len(data1) > 0 && len(data2) > 0 {
		return bytes.Equal(data1, data2)
	}
	return a.Body == nil && b.Body == nil
}

func testEncode(t *testing.T, size int) {
	var netMsg = AllocNetMessage()
	netMsg.Command = 1001
	netMsg.Seq = 1

	if size > 0 {
		netMsg.Data = strutil.RandBytes(size)
	}

	var buf bytes.Buffer
	if err := EncodeMsgTo(netMsg, nil, &buf); err != nil {
		t.Fatalf("%v", err)
	}
	var msg2 = AllocNetMessage()
	if err := DecodeMsgFrom(&buf, MaxPayloadSize, nil, msg2); err != nil {
		t.Fatalf("%v", err)
	}
	if !isMsgEqual(t, netMsg, msg2) {
		t.Fatalf("size %d not equal", size)
	}
}

func TestCodecEncode(t *testing.T) {
	testEncode(t, 0)
	testEncode(t, 64)
	testEncode(t, DefaultCompressThreshold)
	testEncode(t, MaxPayloadSize-MaxPayloadSize/1000)
	//testEncode(t, MaxPayloadSize)
}
