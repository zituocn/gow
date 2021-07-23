/*
color.go

控制台颜色输出

*/

package logy

import (
	"sync"
	"time"
)

const (
	endColor = "\033[0m"
)

var (
	logColor = []string{
		LevelDebug:  "\033[1;34m",
		LevelInfo:   "\033[1;37m",
		LevelNotice: "\033[1;33m",
		LevelWarn:   "\033[1;32m",
		LevelError:  "\033[1;31m",
		LevelPanic:  "\033[1;31m",
		LevelFatal:  "\033[1;31m",
	}
)

type colorWriter struct {
	mu     *sync.Mutex
	b      []byte
	writer Writer
}

// WriteLog colorWrite WriteLog
func (w *colorWriter) WriteLog(now time.Time, level int, b []byte) {
	w.b = w.b[:0]
	w.b = append(w.b, logColor[level]...)
	w.b = append(w.b, b...)
	w.b = append(w.b, endColor...)
	w.writer.WriteLog(now, level, w.b)
}

// WithColor use colorWriter
func WithColor(w Writer) Writer {
	return &colorWriter{writer: w}
}
