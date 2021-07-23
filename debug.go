package gow

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var (
	DefaultWriter      io.Writer = os.Stdout
	DefaultErrorWriter io.Writer = os.Stderr
)

func IsDebugging() bool {
	return true
}

// debugPrint debugPrint
//	debugPrint(123,456,"abc")
func debugPrint(values ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString("[gow-debug] ")
	for i := 0; i < len(values); i++ {
		buffer.WriteString("%v ")
	}
	fmt.Fprintf(DefaultWriter, buffer.String(), values...)
}
