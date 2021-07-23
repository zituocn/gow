package logy

import (
	"time"
)

type multiWriter struct {
	ws []Writer
}

func (m *multiWriter) WriteLog(now time.Time, level int, b []byte) {
	for _, w := range m.ws {
		w.WriteLog(now, level, b)
	}
}

// MultiWriter return a new Writer
func MultiWriter(ws ...Writer) Writer {
	return &multiWriter{
		ws: ws,
	}
}
