// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package logger

import (
	"go.uber.org/zap"
)

var (
	_logger = NewConfig().Build() // core logger
	_sugar  = _logger.Sugar()     // sugared logger
)

func New(cfg Config) *zap.Logger {
	return cfg.Build()
}

func Default() *zap.Logger {
	return _logger
}

func Sugared() *zap.SugaredLogger {
	return _sugar
}

func Setup(cfg Config) {
	var l = New(cfg)
	_logger = l
	_sugar = l.Sugar()
}

func Debugf(format string, args ...interface{}) {
	_sugar.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	_sugar.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	_sugar.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	_sugar.Errorf(format, args...)
}

func DPanicf(format string, args ...interface{}) {
	_sugar.DPanicf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	_sugar.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	_sugar.Fatalf(format, args...)
}
