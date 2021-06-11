package gow

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	DefaultWriter io.Writer = os.Stdout
	DefaultErrorWriter io.Writer = os.Stderr
)

func IsDebugging() bool {
	return true
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(DefaultWriter, "[gow-debug] "+format, values...)
	}
}
