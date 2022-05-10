/*

std.go

default std Logger

use:

	logy.Info("");
	logy.Infof("%s",v)
	logy.SetOutput(...)

*/

package logy

import (
	"bytes"
	"os"
)

var std *Logger

func init() {
	std = NewLogger(WithColor(NewWriter(os.Stdout)), StdFlags, LevelDebug)
	std.SetCallDepth(std.GetCallDepth() + 1)
}

//============ set /get ==============

func SetFlags(flag int) {
	std.SetFlags(flag)
}

func SetOutput(w Writer, prefix string) {
	std.SetOutput(w, prefix)
}

func SetCallDepth(depth int) {
	std.SetCallDepth(depth)
}

func GetWriterType() string {
	return std.GetWriterType()
}

// =============== logy.Info ==============

// Info info v
//	logy.Info("test")
func Info(v ...interface{}) {
	std.Info(getFormat(len(v)), v...)
}

// Debug debug v
func Debug(v ...interface{}) {
	std.Debug(getFormat(len(v)), v...)
}

// Error error v
func Error(v ...interface{}) {
	std.Error(getFormat(len(v)), v...)
}

// Warn warn v
func Warn(v ...interface{}) {
	std.Warn(getFormat(len(v)), v...)
}

// Fatal Fatal v
func Fatal(v ...interface{}) {
	std.Fatal(getFormat(len(v)), v...)
}

// Panic Panic v
func Panic(v ...interface{}) {
	std.Panic(getFormat(len(v)), v...)
}

// Notice notice v
func Notice(v ...interface{}) {
	std.Notice(getFormat(len(v)), v...)
}

// Infof need format
//	logy.Infof("user :%s",user.Username)
func Infof(format string, v ...interface{}) {
	std.Info(format, v...)
}

func Noticef(format string, v ...interface{}) {
	std.Notice(format, v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debug(format, v...)
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

func getFormat(length int) string {
	buffer := &bytes.Buffer{}
	for i := 0; i < length; i++ {
		buffer.WriteString("%v ")
	}
	return buffer.String()
}
