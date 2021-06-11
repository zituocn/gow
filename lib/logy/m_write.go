package logy

import "time"

type multiWriter struct {
	ws []Writer
}

func (mw *multiWriter) WriteLog(t time.Time, level int, b []byte) {
	for _, w := range mw.ws {
		w.WriteLog(t, level, b)
	}
}

// MultiWriter return a Writer interface
func MultiWriter(wr ...Writer) Writer {
	return &multiWriter{ws: wr}
}
