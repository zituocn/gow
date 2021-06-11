package logy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	Ldate         = 1 << iota // the date in the local time zone: 2006-01-02
	Ltime                     // the time in the local time zone: 15:04:05
	Lmicroseconds             // microsecond resolution: 01:23:23.123.  assumes Ltime.
	Llongfile                 // full file name and line number: /a/b/c/d.go:23
	Lshortfile                // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                      // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmodule                   // module name
	Llevel                    // the level of the logy

	LstdFlags = Ldate | Lmicroseconds | Lshortfile | Llevel // initial values for the standard logger
)

const (
	levelTest = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

var levels = []string{
	"[T]",
	"[D]",
	"[I]",
	"[N]",
	"[W]",
	"[E]",
	"[P]",
	"[F]",
}

// Writer writer interface
type Writer interface {
	WriteLog(t time.Time, level int, b []byte)
}

// writer 实现
type writer struct {
	w io.Writer
}

func (wr writer) WriteLog(t time.Time, level int, b []byte) {
	wr.w.Write(b)
}

// NewWriter return a Writer.
func NewWriter(w io.Writer) Writer {
	return writer{w: w}
}

//===========================logger======================

type LogData struct {
	Prefix  string    `json:"prefix"`  // log prefix
	Level   int       `json:"level"`   // log level
	Msg     string    `json:"msg"`     // msg
	Created time.Time `json:"created"` // time
}

type Logger struct {
	pool      *sync.Pool
	flag      int
	level     int
	out       Writer
	callDepth int
	prefix    string
}

// NewLogger return a new Logger
func NewLogger(w Writer, flag int, level int) *Logger {
	return &Logger{
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
		flag:      flag,
		level:     level,
		out:       w,
		callDepth: 2,
	}
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Prefix() string {
	return l.prefix
}

func (l *Logger) Flags() int {
	return l.flag
}

func (l *Logger) SetFlags(flag int) {
	l.flag = flag
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) SetOutput(w Writer, prefix string) {
	l.out = w
	l.prefix = prefix
}

func (l *Logger) SetCallDepth(depth int) {
	l.callDepth = depth
}

func (l *Logger) CallDepth() int {
	return l.callDepth
}

func (l *Logger) formatHeader(buf *bytes.Buffer, t time.Time, file string, line int, lvl int) {
	if l.prefix != "" {
		buf.WriteByte('[')
		buf.WriteString(l.prefix)
		buf.WriteByte(']')
		buf.WriteByte(' ')
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			buf.WriteByte('/')
			itoa(buf, int(month), 2)
			buf.WriteByte('/')
			itoa(buf, day, 2)
			buf.WriteByte(' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			buf.WriteByte(':')
			itoa(buf, min, 2)
			buf.WriteByte(':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				buf.WriteByte('.')
				itoa(buf, t.Nanosecond()/1e6, 3)
			}
			buf.WriteByte(' ')
		}
	}
	if l.flag&Llevel != 0 {
		buf.WriteString(levels[lvl])
		buf.WriteByte(' ')
	}
	if l.flag&Lmodule != 0 {
		buf.WriteByte('[')
		buf.WriteString(moduleOf(file))
		buf.WriteByte(']')
		buf.WriteByte(' ')
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		buf.WriteString(file)
		buf.WriteByte(':')
		itoa(buf, line, -1)
		buf.WriteString(": ")
	}
}

func (l *Logger) output(lvl int, s string) {
	now := time.Now()
	var file string
	var line int
	if l.flag&(Lshortfile|Llongfile) != 0 {
		var ok bool
		_, file, line, ok = runtime.Caller(l.callDepth)
		if !ok {
			file = "???"
			line = 0
		}
	}
	buf := l.pool.Get().(*bytes.Buffer)
	buf.Reset()
	l.formatHeader(buf, now, file, line, lvl)
	buf.WriteString(s)
	if len(s) > 0 && s[len(s)-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.out.WriteLog(now, lvl, buf.Bytes())
	l.pool.Put(buf)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug < l.level {
		return
	}
	l.output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(format string, v ...interface{}) {
	if LevelInfo < l.level {
		return
	}
	l.output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Notice(format string, v ...interface{}) {
	if LevelNotice < l.level {
		return
	}
	l.output(LevelNotice, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if LevelWarn < l.level {
		return
	}
	l.output(LevelWarn, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	if LevelError < l.level {
		return
	}
	l.output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(format string, v ...interface{}) {
	if LevelPanic < l.level {
		return
	}
	s := fmt.Sprintf(format, v...)
	l.output(LevelPanic, s)
	panic(s)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	if LevelFatal < l.level {
		return
	}
	l.output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func itoa(buf *bytes.Buffer, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		buf.WriteByte('0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}

	// avoid slicing b to avoid an allocation.
	for bp < len(b) {
		buf.WriteByte(b[bp])
		bp++
	}
}

func moduleOf(file string) string {
	pos := strings.LastIndex(file, "/")
	if pos != -1 {
		pos1 := strings.LastIndex(file[:pos], "/src/")
		if pos1 != -1 {
			return file[pos1+5 : pos]
		}
	}
	return "UNKNOWN"
}
