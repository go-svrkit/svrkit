// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package zlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createLoggerAndPrint(t *testing.T, conf *Config, msg string) {
	var logger = conf.Build()
	assert.NotNil(t, logger)
	logger.Info(msg)
}

func TestConfig_Build(t *testing.T) {
	var conf = NewConfig()
	createLoggerAndPrint(t, conf, "hello")

	conf.Encoding = "json"
	conf.TimeLayout = "2006-01-02T15:04:05.000000Z0700"
	createLoggerAndPrint(t, conf, "hello")
}
