// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package codec

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/codec/testdata"
)

func TestHasValidSuffix(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Req", true},
		{"FooReq", true},
		{"FooReq-", false},
		{"Ack", true},
		{"FooAck", true},
		{"FooAck-", false},
		{"Ntf", true},
		{"Ntf-", false},
	}
	for i, tt := range tests {
		if got := HasValidSuffix(tt.input); got != tt.want {
			t.Errorf("case %d IsRequestMessage() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestIsReqMessage(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Req", true},
		{"FooReq", true},
		{"FooReq-", false},
	}
	for i, tt := range tests {
		if got := IsReqMessage(tt.input); got != tt.want {
			t.Errorf("case %d IsRequestMessage() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestIsAckMessage(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Ack", true},
		{"FooAck", true},
		{"FooAck-", false},
	}
	for i, tt := range tests {
		if got := IsAckMessage(tt.input); got != tt.want {
			t.Errorf("case %d IsAckMessage() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestGetPairingAckName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Req", ""},
		{"FooReq", "FooAck"},
		{"FooReqq", ""},
	}
	for i, tt := range tests {
		if got := GetPairingAckName(tt.input); got != tt.want {
			t.Errorf("case %d GetPairingAckName() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestGetPairingAckNameOf(t *testing.T) {
	defer Clear()

	assert.Nil(t, Register("testdata.BuildReq"))
	assert.Nil(t, Register("testdata.BuildAck"))

	tests := []struct {
		input uint32
		want  string
	}{
		{0, ""},
		{GetMessageId("testdata.BuildReq"), "testdata.BuildAck"},
	}
	for i, tt := range tests {
		if got := GetPairingAckNameOf(tt.input); got != tt.want {
			t.Errorf("case %d GetPairingAckNameOf() = %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestRegister(t *testing.T) {
	defer Clear()

	assert.Nil(t, Register("testdata.BuildReq"))

	var name = "testdata.BuildReq"
	var hash = GetMessageId(name)
	assert.Greater(t, hash, uint32(0))
	assert.Equal(t, GetMessageFullName(hash), name)
	var req testdata.BuildReq
	assert.Equal(t, GetMessageType(hash).String(), reflect.TypeOf(req).String())
	assert.Equal(t, hash, GetMessageIdOf(&req))
}

func TestCreateMessage(t *testing.T) {
	defer Clear()

	var name = "testdata.BuildReq"
	assert.Nil(t, Register(name))

	var hash = GetMessageId(name)
	var msg = CreateMessageByID(hash)
	assert.NotNil(t, msg)
	req, ok := msg.(*testdata.BuildReq)
	assert.True(t, ok)
	assert.NotNil(t, req)

	var msg2 = CreateMessageByName(name)
	assert.NotNil(t, msg2)
	req, ok = msg.(*testdata.BuildReq)
	assert.True(t, ok)
	assert.NotNil(t, req)
}

func TestCreatePairingAck(t *testing.T) {
	defer Clear()

	assert.Nil(t, Register("testdata.BuildReq"))
	assert.Nil(t, Register("testdata.BuildAck"))

	var name = "testdata.BuildReq"
	var msg = CreatePairingAck(name)
	assert.NotNil(t, msg)
	ack, ok := msg.(*testdata.BuildAck)
	assert.True(t, ok)
	assert.NotNil(t, ack)
}
