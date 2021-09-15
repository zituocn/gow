package gow

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// nameOfFunction return func name
//	use reflect
func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// cleanPath clear route path
func cleanPath(p string, ignore bool) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	if ignore && len(np) > 1 && np[len(np)-1] == '/' {
		np = np[0 : len(np)-1]
	}
	return np
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

// getAddress return HTTP addr
//	getAddress(9090)
//	getAddress(":9090")
//	getAddress("127.0.0.1:9090")
func (engine *Engine) getAddress(args ...interface{}) string {
	var (
		host string
		port int
	)
	switch len(args) {
	case 0:
		host, port = getHostAndPort(engine.httpAddr)
	case 1:
		switch arg := args[0].(type) {
		case string:
			host, port = getHostAndPort(args[0].(string))
		case int:
			port = arg
		}
	case 2:
		if arg, ok := args[0].(string); ok {
			host = arg
		}
		if arg, ok := args[1].(int); ok {
			port = arg
		}
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	return addr
}

// getHostAndPort split like 127.0.0.1:8080
func getHostAndPort(addr string) (host string, port int) {
	addrs := strings.Split(addr, ":")
	if len(addrs) == 1 {
		host = ""
		port, _ = strconv.Atoi(addrs[0])
	} else if len(addrs) >= 2 {
		host = addrs[0]
		port, _ = strconv.Atoi(addrs[1])
	}
	return
}
