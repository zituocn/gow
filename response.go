/*

response.go

*/

package gow

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

// ResponseWriter interface
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// Status Returns the HTTP response status code of the current request.
	Status() int

	// Size Returns the number of bytes already written into the response http body.
	//  See Written()
	Size() int

	// WriteString Writes the string into the response body.
	WriteString(string) (int, error)

	// Written Returns true if the response body was already written.
	Written() bool

	// WriteHeaderNow Forces to write the http header (status code + headers).
	WriteHeaderNow()

	// Pusher get the http.Pusher for server push
	Pusher() http.Pusher
}

// responseWriter
type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

var _ ResponseWriter = &responseWriter{}

func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
}

// WriteHeader write header
func (w *responseWriter) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			debugPrint(fmt.Sprintf("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code))
		}
		w.status = code
	}
}

// WriteHeaderNow WriteHeaderNow
func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

// Write write []byte
func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

// WriteString write string
func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return
}

// Status return response status code
func (w *responseWriter) Status() int {
	return w.status
}

// Size return response size
func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify interface.
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flush interface.
func (w *responseWriter) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *responseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
