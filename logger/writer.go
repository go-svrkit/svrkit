// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Writer write log to file and console
type Writer struct {
	lumberjack.Logger

	console *os.File
}

func NewWriter(filename, console string, maxSize, maxBackup int, localtime bool) *Writer {
	if filename == "" {
		if s, err := os.Executable(); err == nil {
			filename = fmt.Sprintf("logs/%s.log", filepath.Base(s))
		}
	}
	w := &Writer{
		Logger: lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackup,
			LocalTime:  localtime,
		},
	}
	switch console {
	case "stderr":
		w.console = os.Stderr
	case "stdout":
		w.console = os.Stdout
	}
	return w
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if w.console != nil {
		w.console.Write(p)
	}
	return w.Logger.Write(p)
}
