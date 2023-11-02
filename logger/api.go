// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_logger = NewConfig().Build() // core logger
	_sugar  = _logger.Sugar()     // sugared logger
)

func Sugared() *zap.SugaredLogger {
	return _sugar
}

func Default() *zap.Logger {
	return _logger
}

func SetDefault(log *zap.Logger) {
	_logger = log
	_sugar = log.Sugar()
}

func InitBy(cfg *Config) {
	var log = cfg.Build()
	SetDefault(log)
}

func InitWith(callerSkip int, core zapcore.Core) {
	var log = zap.New(core,
		zap.AddCallerSkip(callerSkip),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.PanicLevel))
	SetDefault(log)
}

func Level() zapcore.Level {
	return _logger.Level()
}

func Sync() error {
	return _sugar.Sync()
}

func With(fields ...zap.Field) *zap.Logger {
	return _logger.With(fields...)
}

func Debug(args ...interface{}) {
	_sugar.Debug(args...)
}

func Debugw(msg string, kvPairs ...interface{}) {
	_sugar.Debugw(msg, kvPairs...)
}

func Debugf(format string, args ...interface{}) {
	_sugar.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	_sugar.Debugln(args...)
}

func Info(args ...interface{}) {
	_sugar.Info(args...)
}

func Infow(msg string, kvPairs ...interface{}) {
	_sugar.Infow(msg, kvPairs...)
}

func Infof(format string, args ...interface{}) {
	_sugar.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	_sugar.Infoln(args...)
}

func Warn(args ...interface{}) {
	_sugar.Warn(args...)
}

func Warnw(msg string, kvPairs ...interface{}) {
	_sugar.Warnw(msg, kvPairs...)
}

func Warnf(format string, args ...interface{}) {
	_sugar.Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	_sugar.Warnln(args...)
}

func Error(args ...interface{}) {
	_sugar.Error(args...)
}

func Errorw(msg string, kvPairs ...interface{}) {
	_sugar.Errorw(msg, kvPairs...)
}

func Errorf(format string, args ...interface{}) {
	_sugar.Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	_sugar.Errorln(args...)
}

func DPanic(args ...interface{}) {
	_sugar.DPanic(args...)
}

func DPanicw(msg string, kvPairs ...interface{}) {
	_sugar.DPanicw(msg, kvPairs...)
}

func DPanicf(format string, args ...interface{}) {
	_sugar.DPanicf(format, args...)
}

func DPanicln(args ...interface{}) {
	_sugar.DPanicln(args...)
}

func Panic(args ...interface{}) {
	_sugar.Panic(args...)
}

func Panicw(msg string, kvPairs ...interface{}) {
	_sugar.Panicw(msg, kvPairs...)
}

func Panicf(format string, args ...interface{}) {
	_sugar.Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	_sugar.Panicln(args...)
}

func Fatal(args ...interface{}) {
	_sugar.Fatal(args...)
}

func Fatalw(msg string, kvPairs ...interface{}) {
	_sugar.Fatalw(msg, kvPairs...)
}

func Fatalf(format string, args ...interface{}) {
	_sugar.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	_sugar.Fatalln(args...)
}
