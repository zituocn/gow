package logy

import (
	"io"
	"os"
	"time"
)

type mutexWriter struct {
	w io.Writer
}

// WriteLog write log
func (mw *mutexWriter) WriteLog(t time.Time, level int, s []byte) {
	mw.w.Write(s)
}

// NewMutexWriter returns a currently safe writer.
func NewMutexWriter(w io.Writer) Writer {
	return writer{w: w}
}

var std *Logger

func init() {
	std = NewLogger(NewMutexWriter(os.Stdout), LstdFlags, LevelDebug)
	std.SetCallDepth(std.CallDepth() + 1)
}

func Flags() int {
	return std.Flags()
}

func SetFlags(flag int) {
	std.SetFlags(flag)
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func SetOutput(w Writer, prefix string) {
	std.SetOutput(w, prefix)
}

func SetCallDepth(depth int) {
	std.SetCallDepth(depth)
}

func CallDepth() int {
	return std.CallDepth()
}

func Debug(v ...interface{}) {
	std.Debug("%v", v...)
}

func Info(v ...interface{}) {
	std.Info("%v", v...)
}

func Notice(v ...interface{}) {
	std.Notice("%v", v...)
}

func Warn(v ...interface{}) {
	std.Warn("%v", v...)
}

func Error(v ...interface{}) {
	std.Error("%v", v...)
}

func Panic(v ...interface{}) {
	std.Panic("%v", v...)
}

func Fatal(v ...interface{}) {
	std.Fatal("%v", v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debug(format, v...)
}

func Infof(format string, v ...interface{}) {
	std.Info(format, v...)
}

func Noticef(format string, v ...interface{}) {
	std.Notice(format, v...)
}

func Warnf(format string, v ...interface{}) {
	std.Warn(format, v...)
}

func Errorf(format string, v ...interface{}) {
	std.Error(format, v...)
}

func Panicf(format string, v ...interface{}) {
	std.Panic(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	std.Fatal(format, v...)
}
