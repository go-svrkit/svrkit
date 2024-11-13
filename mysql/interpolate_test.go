package mysql

import (
	"bytes"
	"testing"
	"time"
)

func Test_escapeBytesBackslash(t *testing.T) {
	var escapeBytes = []byte{
		'\x00', '\n', '\r', '\x1a', '\'', '"', '\\',
		' ', 'a', 'B', 'c', 'D',
	}
	const expectString = "\\0\\n\\r\\Z\\'\\\"\\\\ aBcD"

	var buf bytes.Buffer
	escapeBytesBackslash(&buf, escapeBytes)
	var result = buf.String()

	if result != expectString {
		t.Fatalf("[%s] != [%s]", result, expectString)
	}
}

func Test_interpolateParams(t *testing.T) {
	var epoch = time.Date(2000, 1, 1, 12, 0, 0, 0, time.Local)
	var buf bytes.Buffer
	var args = []interface{}{
		int64(42),
		"hello",
		[]byte("world"),
		true,
		3.14159,
		epoch,
	}
	err := InterpolateParams(&buf, "SELECT ?+?+?+?+?+?", args)
	if err != nil {
		t.Errorf("interpolateParams: %v", err)
		return
	}
	var result = buf.String()
	const expected = `SELECT 42+'hello'+_binary'world'+1+3.14159+'2000-01-01 12:00:00'`
	if result != expected {
		t.Errorf("Expected: %q\nGot: %q", expected, result)
	}
}
