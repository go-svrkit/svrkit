// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package zlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	assert.NotNil(t, Logger())
	assert.NotNil(t, Sugared())
	assert.Equal(t, Level(), zapcore.DebugLevel)
}

func TestSetLogger(t *testing.T) {
	var cfg = NewConfig()
	var core = CreateZapCore(cfg)
	assert.NotNil(t, core)
	SetLoggerWith(core, 1)

	Infof("hello")
}

func TestSync(t *testing.T) {
	Infof("hello")
	assert.Nil(t, Sync())
}
