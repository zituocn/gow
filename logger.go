/*

gin logger.go

*/

package gow

import (
	"github.com/mattn/go-isatty"
	"github.com/zituocn/logx"
	"net/http"
	"os"
	"time"
)

type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var (
	consoleColorMode = autoColor
)

type LogFormatterParams struct {
	Request *http.Request

	//AppName
	AppName string

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
	// user agent string
	UserAgent string
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the logx.
func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}

// Logger default logger
func Logger() HandlerFunc {
	out := DefaultWriter
	var notLogged []string
	isTerm := true
	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}
	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) {

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request:    c.Request,
				isTerm:     isTerm,
				Keys:       c.Keys,
				TimeStamp:  time.Now(),
				Method:     c.Request.Method,
				StatusCode: c.Writer.Status(),
				ClientIP:   c.GetIP(),
				BodySize:   c.Writer.Size(),
				UserAgent:  c.UserAgent(),
			}
			param.Latency = param.TimeStamp.Sub(start)

			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path

			var statusColor, methodColor, resetColor string

			if param.IsOutputColor() && !c.IsProd() {
				statusColor = param.StatusCodeColor()
				methodColor = param.MethodColor()
				resetColor = param.ResetColor()
			}

			if param.Latency > time.Minute {
				// Truncate in a golang < 1.8 safe way
				param.Latency = param.Latency - param.Latency%time.Second
			}

			logx.Infof("%s %3d %s| %13v | %15s |%s %-7s %s %s%s",
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.ErrorMessage,
			)
		}
	}
}

// LoggerLong write useragent
func LoggerLong() HandlerFunc {
	out := DefaultWriter
	var notLogged []string
	isTerm := true
	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}
	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) {

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request:    c.Request,
				isTerm:     isTerm,
				Keys:       c.Keys,
				TimeStamp:  time.Now(),
				Method:     c.Request.Method,
				StatusCode: c.Writer.Status(),
				ClientIP:   c.GetIP(),
				BodySize:   c.Writer.Size(),
				UserAgent:  c.UserAgent(),
			}
			param.Latency = param.TimeStamp.Sub(start)

			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path

			var statusColor, methodColor, resetColor string

			if param.IsOutputColor() && !c.IsProd() {
				statusColor = param.StatusCodeColor()
				methodColor = param.MethodColor()
				resetColor = param.ResetColor()
			}

			if param.Latency > time.Minute {
				// Truncate in a golang < 1.8 safe way
				param.Latency = param.Latency - param.Latency%time.Second
			}

			logx.Infof("%s %3d %s| %13v | %15s |%s %-7s %s %s %s",
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.UserAgent,
			)
		}
	}
}
