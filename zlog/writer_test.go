// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package zlog

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter_Write(t *testing.T) {
	var w = NewWriter("", "stdout", 100, 1)
	_, err := w.Write([]byte("hello"))
	assert.Nil(t, err)
}

func TestIsTerminal(t *testing.T) {
	assert.True(t, IsTerminal(os.Stdout))
	assert.True(t, IsTerminal(os.Stderr))
	var sb strings.Builder
	assert.False(t, IsTerminal(&sb))
}
