// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestReadHeadBody(t *testing.T) {

}

func TestProcessHeaderFlags(t *testing.T) {

}

func TestEncodeMsgTo(t *testing.T) {

}

func TestDecodeMsgFrom(t *testing.T) {

}
