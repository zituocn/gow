package logy

import (
	"io"
	"time"
)

// Writer  interface
type Writer interface {
	WriteLog(time.Time, int, []byte)
}

// NewWriter return a Writer interface
func NewWriter(w io.Writer) Writer {
	return writer{w: w}
}


type writer struct {
	w io.Writer
}

func (wr writer) WriteLog(t time.Time, level int, p []byte) {
	wr.w.Write(p)
}
