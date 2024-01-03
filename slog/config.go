// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package slog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 配置参数
type Config struct {
	LocalTime  bool   `json:"localtime"`   // use local or UTC time
	Level      string `json:"level"`       // DEBUG, INFO, WARN, DPANIC, PANIC, FATAL
	Encoding   string `json:"encoding"`    // "json" or "console"
	Filename   string `json:"filename"`    //
	MaxSize    int    `json:"max_size"`    // defaults to 100 MB
	MaxBackups int    `json:"max_backups"` // maximum number of old log files to retain
	CallerSkip int    `json:"caller_skip"` // skip some number of extra stack frames (0 = don't skip)
}

func NewConfig() *Config {
	return &Config{
		Level:      "debug",
		Encoding:   "console",
		MaxSize:    200,
		MaxBackups: 10,
		LocalTime:  true,
		CallerSkip: 1,
	}
}

func (c *Config) Build() *zap.Logger {
	var core = CreateZapCoreBy(c)
	return zap.New(core,
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.PanicLevel))
}

// CreateZapCoreBy 根据配置创建zapcore.Core
func CreateZapCoreBy(c *Config) zapcore.Core {
	var encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if c.Encoding == "console" && IsTerminal(os.Stdout) {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	var enc zapcore.Encoder
	if c.Encoding == "json" {
		enc = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		enc = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var level = atomicLevel(c.Level)

	// application should redirect stderr to stdout first
	var w = NewWriter(c.Filename, "stdout", c.MaxSize, c.MaxBackups, c.LocalTime)
	var sink = zapcore.AddSync(w)
	return zapcore.NewCore(enc, sink, zap.NewAtomicLevelAt(level))
}

func atomicLevel(logLevel string) zapcore.Level {
	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		panic(err)
	}
	return level
}
