// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package zlog

import (
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LayoutISO8601 = "2006-01-02T15:04:05.000Z0700"

// Config 配置参数
type Config struct {
	Level      string `json:"level"`       // DEBUG, INFO, WARN, PANIC, FATAL
	Encoding   string `json:"encoding"`    // "json" or "console"
	TimeLayout string `json:"time-layout"` // layout to encoding time
	Filename   string `json:"filename"`    //
	MaxSize    int    `json:"max_size"`    // defaults to 100 MB
	MaxBackups int    `json:"max_backups"` // maximum number of old log files to retain
	CallerSkip int    `json:"caller_skip"` // skip some number of extra stack frames (0 = don't skip)
}

func NewConfig() *Config {
	return &Config{
		Level:      "debug",
		Encoding:   "console",
		TimeLayout: LayoutISO8601,
		MaxSize:    100,
		MaxBackups: 10,
		CallerSkip: 1,
	}
}

func (c *Config) initDefault() {
	if c.Level == "" {
		c.Level = "INFO"
	}
	if c.Encoding == "" {
		c.Encoding = "console"
	}
	if c.TimeLayout == "" {
		c.TimeLayout = LayoutISO8601
	}
	if c.MaxSize <= 0 {
		c.MaxSize = 100
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 10
	}
}

func (c *Config) Build() *zap.Logger {
	c.initDefault()
	var core = CreateZapCore(c)
	return zap.New(core,
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.PanicLevel))
}

// CreateZapCore 根据配置创建zapcore.Core
func CreateZapCore(c *Config) zapcore.Core {
	var encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(c.TimeLayout),
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if c.Encoding == "console" && IsTerminal(os.Stdout) {
		if runtime.GOOS == "windows" {
			encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		} else {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
	}
	var enc zapcore.Encoder
	if c.Encoding == "json" {
		enc = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		enc = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var level = atomicLevel(c.Level)

	// application should redirect stderr to stdout first
	var w = NewWriter(c.Filename, "stdout", c.MaxSize, c.MaxBackups)
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
