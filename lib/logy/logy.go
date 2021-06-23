package logy

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	LogDate = 1 << iota
	LogTime
	LogMicroSeconds
	LogLongFile
	LogShortFile
	LogUTC
	LogModule
	LogLevel

	StdFlags = LogDate | LogMicroSeconds | LogShortFile | LogLevel
)

const (
	LevelTest = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

var (
	levels = []string{
		"[T]",
		"[D]",
		"[I]",
		"[N]",
		"[W]",
		"[E]",
		"[P]",
		"[F]",
	}
)

// Logger log struct
type Logger struct {
	flag      int
	level     int
	out       Writer
	callDepth int
	prefix    string
	pool      *sync.Pool
}

// NewLogger return a *Logger
func NewLogger(w Writer, flag int, level int) *Logger {
	return &Logger{
		flag:      flag,
		level:     level,
		out:       w,
		callDepth: 2,
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
	}
}

// Info logy info
func (log *Logger) Info(format string, v ...interface{}) {
	if LevelInfo < log.level {
		return
	}
	log.output(LevelInfo, fmt.Sprintf(format, v...))
}

func (log *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug < log.level {
		return
	}
	log.output(LevelDebug, fmt.Sprintf(format, v...))
}

func (log *Logger) Notice(format string, v ...interface{}) {
	if LevelNotice < log.level {
		return
	}
	log.output(LevelNotice, fmt.Sprintf(format, v...))
}

func (log *Logger) Error(format string, v ...interface{}) {
	if LevelError < log.level {
		return
	}
	log.output(LevelError, fmt.Sprintf(format, v...))
}

func (log *Logger) Warn(format string, v ...interface{}) {
	if LevelWarn < log.level {
		return
	}
	log.output(LevelWarn, fmt.Sprintf(format, v...))
}

func (log *Logger) Panic(format string, v ...interface{}) {
	if LevelPanic < log.level {
		return
	}
	s := fmt.Sprintf(format, v...)
	log.output(LevelWarn, s)
	panic(s)
}

func (log *Logger) Fatal(format string, v ...interface{}) {
	if LevelFatal < log.level {
		return
	}
	log.output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(-1)
}

//=================== get/set===============

func (log *Logger) SetCallDepth(depth int) {
	log.callDepth = depth
}

func (log *Logger) GetCallDepth() int {
	return log.callDepth
}

func (log *Logger) SetFlags(flag int) {
	log.flag = flag
}

func (log *Logger) SetPrefix(prefix string) {
	log.prefix = prefix
}

// GetWriterType return writer type
func (log *Logger) GetWriterType() string {
	return reflect.TypeOf(log.out).String()
}

// SetOutput set output
//	w => Writer interface
//
func (log *Logger) SetOutput(w Writer, prefix string) {
	log.out = w
	log.prefix = prefix
}

//====================== private ===================

func (log *Logger) writeHeader(buffer *bytes.Buffer, now time.Time, file string, line int, level int) {

	// prefix
	if log.prefix != "" {
		buffer.WriteByte('[')
		buffer.WriteString(log.prefix)
		buffer.WriteByte(']')
		buffer.WriteByte(' ')
	}

	// datetime
	if log.flag&(LogDate|LogTime|LogMicroSeconds) != 0 {

		if log.flag&LogDate != 0 {
			year, month, day := now.Date()
			formatWrite(buffer, year, 4)
			buffer.WriteByte('/')

			formatWrite(buffer, int(month), 2)
			buffer.WriteByte('/')

			formatWrite(buffer, day, 2)
			buffer.WriteByte(' ')
		}

		if log.flag&(LogTime|LogMicroSeconds) != 0 {
			hour, min, sec := now.Clock()
			formatWrite(buffer, hour, 2)
			buffer.WriteByte(':')

			formatWrite(buffer, min, 2)
			buffer.WriteByte(':')

			formatWrite(buffer, sec, 2)

			if log.flag&LogMicroSeconds != 0 {
				buffer.WriteByte('.')
				formatWrite(buffer, now.Nanosecond()/1e6, 3)
			}
			buffer.WriteByte(' ')
		}
	}

	// level
	if log.flag&LogLevel != 0 {
		buffer.WriteString(levels[level])
		buffer.WriteByte(' ')
	}

	// package
	if log.flag&LogModule != 0 {
		buffer.WriteByte('[')
		buffer.WriteString(moduleOf(file))
		buffer.WriteByte(']')
		buffer.WriteByte(' ')
	}

	// filename and line
	if log.flag&(LogShortFile|LogLongFile) != 0 {
		if log.flag&LogShortFile != 0 {
			i := strings.LastIndex(file, "/")
			file = file[i+1:]
		}
		buffer.WriteString(file)
		buffer.WriteByte(':')
		formatWrite(buffer, line, -1)
		buffer.WriteByte(':')
		buffer.WriteByte(' ')
	}
}

func (log *Logger) output(level int, msg string) {
	var (
		now  = time.Now()
		file string
		line int
	)

	if log.flag&(LogShortFile|LogLongFile) != 0 {
		ok := false
		_, file, line, ok = runtime.Caller(log.callDepth)
		if !ok {
			file = "???"
			line = 0
		}
	}

	buffer := log.pool.Get().(*bytes.Buffer)
	buffer.Reset()

	// write header
	log.writeHeader(buffer, now, file, line, level)

	// write msg
	buffer.WriteString(msg)

	if len(msg) > 0 && msg[len(msg)-1] != '\n' {
		buffer.WriteByte('\n')
	}

	log.out.WriteLog(now, level, buffer.Bytes())
	log.pool.Put(buffer)
}

func formatWrite(buffer *bytes.Buffer, i int, wid int) {
	var u = uint(i)
	if u == 0 && wid <= 1 {
		buffer.WriteByte('0')
		return
	}
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}

	for bp < len(b) {
		buffer.WriteByte(b[bp])
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
