/*

gin route_group


*/

package gow

import (
	"net/http"
	"path"
	"strings"
)

var (
	httpMethods = map[string]bool{
		"GET":     true,
		"POST":    true,
		"DELETE":  true,
		"PUT":     true,
		"PATCH":   true,
		"OPTIONS": true,
		"HEAD":    true,
	}
)

// H map[string]interface{}
type H map[string]interface{}

type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

// IRoutes all router handler interface.
type IRoutes interface {
	Handle(string, string, ...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes

	StaticFile(string, string) IRoutes
	Static(string, string) IRoutes
	StaticFS(string, http.FileSystem) IRoutes
}

// RouterGroup is used internally to configure router.
type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

var _IRouter = &RouterGroup{}

// Use add middleware to the group
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

// Group create a new route group.
//	v1:=r.Group("/v1")
//	v1.GET("/user",...)
func (group *RouterGroup) Group(path string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(path),
		engine:   group.engine,
	}
}

// Handle registers a new request handle and middleware with the given path and method
//	r.Handle("GET","/",handler)
//	r.Handler("GET,POST","/",handler)
func (group *RouterGroup) Handle(method, path string, handlers ...HandlerFunc) IRoutes {
	methods := strings.Split(method, ",")
	for _, m := range methods {
		group.handle(m, path, handlers)
	}
	return group.returnObj()
}

// GET register GET handler
func (group *RouterGroup) GET(path string, handlers ...HandlerFunc) IRoutes {
	group.handle(http.MethodGet, path, handlers)
	return group.returnObj()
}

// POST register POST handler
func (group *RouterGroup) POST(path string, handlers ...HandlerFunc) IRoutes {
	group.handle(http.MethodPost, path, handlers)
	return group.returnObj()
}

// DELETE register DELETE handler
func (group *RouterGroup) DELETE(path string, handlers ...HandlerFunc) IRoutes {
	group.handle(http.MethodDelete, path, handlers)
	return group.returnObj()
}

// PATCH register PATCH handler
func (group *RouterGroup) PATCH(path string, handlers ...HandlerFunc) IRoutes {
	group.handle(http.MethodPatch, path, handlers)
	return group.returnObj()
}

// PUT register PUT handler
func (group *RouterGroup) PUT(path string, handlers ...HandlerFunc) IRoutes {
	group.handle(http.MethodPut, path, handlers)
	return group.returnObj()
}

// OPTIONS register OPTIONS handler
func (group *RouterGroup) OPTIONS(path string, handlers ...HandlerFunc) IRoutes {
	group.handle(http.MethodOptions, path, handlers)
	return group.returnObj()
}

// HEAD register HEAD handler
func (group *RouterGroup) HEAD(path string, handlers ...HandlerFunc) IRoutes {
	group.handle(http.MethodHead, path, handlers)
	return group.returnObj()
}

// Any register all HTTP methods.
func (group *RouterGroup) Any(path string, handlers ...HandlerFunc) IRoutes {
	group.handle(http.MethodGet, path, handlers)
	group.handle(http.MethodPost, path, handlers)
	group.handle(http.MethodPut, path, handlers)
	group.handle(http.MethodPatch, path, handlers)
	group.handle(http.MethodHead, path, handlers)
	group.handle(http.MethodOptions, path, handlers)
	group.handle(http.MethodDelete, path, handlers)
	return group.returnObj()
}

func (group *RouterGroup) handle(method, path string, handlers HandlersChain) {
	if method == "" {
		panic("need http method")
	}
	if path == "" {
		panic("need route path")
	}
	if len(handlers) == 0 {
		panic("need HandlerFunc func(ctx *Context)")
	}
	method = strings.ToUpper(method)
	if ok := httpMethods[method]; !ok {
		panic("http method " + method + " not supported")
	}
	var match matchValue
	if group.engine.Match(method, path, &match) {
		panic("method: " + method + " path: " + path + " has been defined")
	}
	absolutePath := group.calculateAbsolutePath(path)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(method, absolutePath, handlers)
}

// Static handler static dir
//	r.Static("/static","static")
func (group *RouterGroup) Static(path, root string) IRoutes {
	return group.StaticFS(path, Dir(root, false))
}

// StaticFS r.StaticFS(path,fs)
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := group.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "{static_file_path}")

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
	return group.returnObj()
}

func (group *RouterGroup) createStaticHandler(path string, fs http.FileSystem) HandlerFunc {
	absolutePath := group.calculateAbsolutePath(path)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		if _, noListing := fs.(*onlyFilesFS); noListing {
			c.Writer.WriteHeader(http.StatusNotFound)
		}
		file := c.Param("static_file_path")
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.handlers = group.engine.noRoute
			// Reset index
			c.index = -1
			return
		}
		_ = f.Close()

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

// StaticFile static file
//	r.StaticFile("/favicon.png","/static/img/favicon.png")
func (group *RouterGroup) StaticFile(path, filepath string) IRoutes {
	if strings.Contains(path, ":") || strings.Contains(path, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := func(c *Context) {
		c.File(filepath)
	}
	group.GET(path, handler)
	group.HEAD(path, handler)
	return group.returnObj()
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(path string) string {
	return joinPaths(group.basePath, path)
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}
