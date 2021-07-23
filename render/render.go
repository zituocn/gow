package render

import "net/http"

// Render interface is to be implemented by HTML
type Render interface {
	Render(w http.ResponseWriter, name string, data interface{}) error
	WriteContentType(w http.ResponseWriter)
}

// Delims template delims
type Delims struct {
	Left  string
	Right string
}

// writeContentType
func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
