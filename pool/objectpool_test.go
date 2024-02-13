package pool

import (
	"bytes"
	"testing"
)

func TestAllocBytesBuffer(t *testing.T) {
	var bufferList []*bytes.Buffer
	for i := 0; i < 100; i++ {
		var buf = AllocBytesBuffer()
		bufferList = append(bufferList, buf)
	}
	for _, buf := range bufferList {
		FreeBytesBuffer(buf)
	}
	bufferList = nil
}

func TestAllocBytesReader(t *testing.T) {
	var bufferList []*bytes.Reader
	for i := 0; i < 100; i++ {
		var buf = AllocBytesReader()
		bufferList = append(bufferList, buf)
	}
	for _, buf := range bufferList {
		FreeBytesReader(buf)
	}
	bufferList = nil
}
