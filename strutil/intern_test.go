package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternStr(t *testing.T) {
	var s = "hello,world"
	var interned = InternStr(s)
	assert.Equal(t, s, interned)
}

func TestInternBytes(t *testing.T) {
	var b = []byte("hello,world")
	var interned = InternBytes(b)
	assert.Equal(t, string(b), interned)
}
