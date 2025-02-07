// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package qlog

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter_Write(t *testing.T) {
	var w = NewWriter("", "stdout", 100, 0, 1)
	_, err := w.Write([]byte("hello"))
	assert.Nil(t, err)
}

func TestIsTerminal(t *testing.T) {
	assert.True(t, IsTerminal(os.Stdout))
	assert.True(t, IsTerminal(os.Stderr))
	var sb strings.Builder
	assert.False(t, IsTerminal(&sb))
}

func TestAppendFileLog(t *testing.T) {
	var filename = "test-qlog-writer.log"
	defer os.Remove(filename)

	err := AppendFileLog(filename, "hello")
	assert.Nil(t, err)
	err = AppendFileLog(filename, "world")
	assert.Nil(t, err)
}
