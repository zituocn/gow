package logy

import (
	"sync"
	"time"
)

const colorEnd = "\033[0m"

// 加粗的颜色输出
var logColor = []string{
	LevelDebug:  "\033[1;36m",
	LevelInfo:   "\033[1;37m",
	LevelNotice: "\033[1;33m",
	LevelWarn:   "\033[1;35m",
	LevelError:  "\033[1;31m",
	LevelPanic:  "\033[1;31m",
	LevelFatal:  "\033[1;31m",
}

type colorWriter struct {
	m sync.Mutex
	b []byte
	w Writer
}

func (cw *colorWriter) WriteLog(t time.Time, level int, b []byte) {
	cw.m.Lock()
	cw.b = cw.b[:0]
	cw.b = append(cw.b, logColor[level]...)
	cw.b = append(cw.b, b...)
	cw.b = append(cw.b, colorEnd...)
	cw.w.WriteLog(t, level, cw.b)
	cw.m.Unlock()
}

// WithColor 指定某一个实现的writer使用颜色
func WithColor(w Writer) Writer {
	return &colorWriter{w: w}
}
