package qnet

import (
	"bytes"
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestPipeConn(t *testing.T, addr string) (local, remote net.Conn) {
	var bus = make(chan net.Conn, 1)
	go listenTestServer(t, addr, 1, bus)
	local, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	remote = <-bus
	assert.NotNil(t, remote)
	return
}

func TestReadLenData(t *testing.T) {
	tests := []struct {
		data   []byte
		expect []byte
		err    bool
	}{
		{[]byte{0x12}, nil, true},
		{[]byte{0x0, 0x0}, nil, true},
		{[]byte{0x12, 0x34}, nil, true},
		{[]byte{0xFF, 0xFF}, nil, true},
		{[]byte{0x0, 0x2}, nil, false},
		{[]byte{0x0, 0x6, 0x1, 0x2, 0x3, 0x4}, []byte{0x1, 0x2, 0x3, 0x4}, false},
	}
	for _, tc := range tests {
		out, err := ReadLenData(bytes.NewReader(tc.data), 0x0FFF)
		if tc.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.True(t, bytes.Equal(out, tc.expect))
		}
	}
}

func TestWriteLenData(t *testing.T) {
	tests := []struct {
		data   []byte
		expect []byte
		err    bool
	}{
		{nil, nil, false},
		{[]byte{0x1}, []byte{0x0, 0x3, 0x1}, false},
		{[]byte{0x1, 0x2, 0x3}, []byte{0x0, 0x5, 0x1, 0x2, 0x3}, false},
	}
	for _, tc := range tests {
		var buf bytes.Buffer
		err := WriteLenData(&buf, tc.data)
		var out = buf.Bytes()
		if tc.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.True(t, bytes.Equal(out, tc.expect))
		}
	}
}

func TestGetLocalIPList(t *testing.T) {
	var list = GetLocalIPList()
	assert.True(t, len(list) > 0)
}

func TestReadProtoMessage(t *testing.T) {
	local, remote := createTestPipeConn(t, "localhost:42334")
	defer func() {
		local.Close()
		remote.Close()
	}()

	var req = &PrebuildReq{PosX: 1234, PosZ: 5678}
	assert.Nil(t, WriteProtoMessage(local, req))

	var req2 PrebuildReq
	assert.Nil(t, ReadProtoMessage(remote, &req2))

	assert.True(t, reflect.DeepEqual(req, &req2))
}
