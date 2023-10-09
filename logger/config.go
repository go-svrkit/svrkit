// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

type Hooker interface {
	Name() string
	Fire(entry zapcore.Entry) error
}

// Config 配置参数
type Config struct {
	LocalTime  bool   `json:"localtime"`   // use local or UTC time
	Level      string `json:"level"`       // DEBUG, INFO, WARN, DPANIC, PANIC, FATAL
	Encoding   string `json:"encoding"`    // "json" or "console"
	Filename   string `json:"filename"`    //
	MaxSize    int    `json:"max_size"`    // defaults to 100 megabytes
	MaxBackups int    `json:"max_backups"` // maximum number of old log files to retain
}

func NewConfig() Config {
	return Config{
		Level:      "debug",
		Encoding:   "console",
		MaxSize:    200,
		MaxBackups: 10,
		LocalTime:  true,
	}
}

func (c Config) Build() *zap.Logger {
	var core = c.CreateLogCore()
	return zap.New(core,
		zap.AddCallerSkip(1),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.PanicLevel))
}

// application should redirect stderr to stdout first
func (c Config) CreateLogCore() zapcore.Core {
	var encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if c.Encoding == "console" && isTerminal(os.Stdout) {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	var enc zapcore.Encoder
	if c.Encoding == "json" {
		enc = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		enc = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var level = atomicLevel(c.Level)
	var w = c.CreateLogWriter(c.Filename, "stdout")
	return zapcore.NewCore(enc, zapcore.AddSync(w), zap.NewAtomicLevelAt(level))
}

func (c Config) CreateLogWriter(filename, console string) *Writer {
	return NewWriter(filename, console, c.MaxSize, c.MaxBackups, c.LocalTime)
}

func atomicLevel(logLevel string) zapcore.Level {
	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		panic(err)
	}
	return level
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
