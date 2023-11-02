// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
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
			filename = fmt.Sprintf("%s.log", filepath.Base(s))
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

// AppendFileLog append log to file
func AppendFileLog(filename, format string, a ...interface{}) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, format, a...)
	return err
}

// AppendAppErrorLog append app error log to file, mkdir `logs` first
func AppendAppErrorLog(format string, a ...interface{}) error {
	appname := filepath.Base(os.Args[0])
	if i := strings.LastIndex(appname, "."); i > 0 {
		appname = appname[:i]
	}
	if i := strings.LastIndex(appname, "_"); i > 0 {
		appname = appname[i+1:]
	}
	var filename = fmt.Sprintf("logs/%s_error.log", appname)
	return AppendFileLog(filename, format, a...)
}

func IsTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
