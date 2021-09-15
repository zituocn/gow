/*
gow.go

sam
2021/6/7
*/

package gow

import (
	"fmt"
	"github.com/zituocn/gow/lib/logy"
	"github.com/zituocn/gow/render"
	"html/template"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	defaultMultipartMemory = 32 << 20 // 32 MB
	defaultMode            = "dev"
	DevMode                = "dev"
	ProdMode               = "prod"
	defaultViews           = "views"
	defaultStatic          = "static"
	defaultHttpAddr        = "127.0.0.1:8080"
)

var (
	default404Body = []byte("404 page not found")
)

// HandlerFunc gow handler
//	func(ctx *Context)
type HandlerFunc func(ctx *Context)

// HandlersChain []HandlerFunc
type HandlersChain []HandlerFunc

// Last HandlerFunc
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// Engine gow Engine
type Engine struct {
	AppName  string
	RunMode  string
	httpAddr string

	AutoRender bool
	delims     render.Delims
	FuncMap    template.FuncMap
	Render     render.Render
	viewsPath  string
	staticPath string
	sessionOn  bool
	gzipOn     bool

	// ignoreCase if true ignore case on the route
	//	default: true
	ignoreCase bool

	// default:true
	// ignoreTrailingSlash if true ignore the default trailing slash one the route
	ignoreTrailingSlash bool

	RouterGroup
	HandleMethodNotAllowed bool
	MaxMultipartMemory     int64

	allNoMethod HandlersChain
	allNoRoute  HandlersChain

	noRoute  HandlersChain
	noMethod HandlersChain

	pool      sync.Pool
	trees     []methodTree
	maxParams uint16
}

// New returns a new blank Engine instance without any middleware attached.
//	gzipOn false
//	sessionOn false
//	ignoreCase true
//	ignoreTrailingSlash true
//	AutoRender false
//	RunMode = "dev"
func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		AppName:             "gow",
		RunMode:             "dev",
		httpAddr:            defaultHttpAddr,
		AutoRender:          false,
		FuncMap:             template.FuncMap{},
		delims:              render.Delims{Left: "{{", Right: "}}"},
		viewsPath:           defaultViews,
		staticPath:          defaultStatic,
		gzipOn:              false,
		sessionOn:           false,
		ignoreCase:          true,
		ignoreTrailingSlash: true,
		MaxMultipartMemory:  defaultMultipartMemory,

		trees: make([]methodTree, 0, 9),
	}
	engine.RouterGroup.engine = engine

	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}

	return engine
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	engine.NoRoute(default404Handler)
	return engine
}

// SetAppConfig set engine config
func (engine *Engine) SetAppConfig(app *AppConfig) {
	if app != nil {
		engine.AppName = app.AppName
		engine.RunMode = app.RunMode
		engine.viewsPath = app.Views
		engine.delims = render.Delims{Left: app.TemplateLeft, Right: app.TemplateRight}
		engine.AutoRender = app.AutoRender
		engine.httpAddr = app.HTTPAddr
		engine.sessionOn = app.SessionOn
		engine.gzipOn = app.GzipOn
		engine.ignoreCase = app.IgnoreCase
	}
}

// AddFuncMap add fn func to template func map
func (engine *Engine) AddFuncMap(key string, fn interface{}) {
	engine.FuncMap[key] = fn
}

// SetViews set engine.viewPath = path
func (engine *Engine) SetViews(path string) {
	engine.viewsPath = path
}

// SetGzipOn set gzip on
func (engine *Engine) SetGzipOn(on bool) {
	engine.gzipOn = on
}

// SetIgnoreCase set ignore case on the route
func (engine *Engine) SetIgnoreCase(ignore bool) {
	engine.ignoreCase = ignore
}

// SetIgnoreTrailingSlash set ignore trailing slash one the route
func (engine *Engine) SetIgnoreTrailingSlash(ignore bool) {
	engine.ignoreTrailingSlash = ignore
}

// Run start http service
func (engine *Engine) Run(args ...interface{}) (err error) {
	defer func() { logy.Error(err) }()

	engine.useMiddleware()

	if engine.AutoRender {
		engine.Render = render.HTMLRender{}.NewHTMLRender(engine.viewsPath, engine.FuncMap, engine.delims, engine.AutoRender, engine.RunMode)
	}
	if engine.RunMode == DevMode {
		fmt.Printf("%s\n", logo)
	}

	address := engine.getAddress(args...)
	logy.Infof("[%s] [%s] Listening and serving HTTP on http://%s\n", engine.AppName, engine.RunMode, address)
	err = http.ListenAndServe(address, engine)
	return
}

// RunTLS  start https service
func (engine *Engine) RunTLS(certFile, keyFile string, args ...interface{}) (err error) {
	defer func() { logy.Error(err) }()

	engine.useMiddleware()

	if engine.AutoRender {
		engine.Render = render.HTMLRender{}.NewHTMLRender(engine.viewsPath, engine.FuncMap, engine.delims, engine.AutoRender, engine.RunMode)
	}

	if engine.RunMode == DevMode {
		fmt.Printf("%s\n", logo)
	}

	address := engine.getAddress(args...)
	logy.Infof("[%s] [%s] Listening and serving HTTP on https://%s\n", engine.AppName, engine.RunMode, address)
	err = http.ListenAndServeTLS(address, certFile, keyFile, engine)
	return
}

// RunUnix start unix:/ service
func (engine *Engine) RunUnix(file string) (err error) {
	defer func() { logy.Error(err) }()

	engine.useMiddleware()

	listener, err := net.Listen("unix", file)
	if err != nil {
		return
	}
	defer listener.Close()
	defer os.Remove(file)

	logy.Infof("[%s] [%s] Listening and serving HTTP on unix:/%s\n", engine.AppName, engine.RunMode, file)

	err = http.Serve(listener, engine)
	return
}

// RunFd start fd service
func (engine *Engine) RunFd(fd int) (err error) {
	defer func() { logy.Error(err) }()

	if engine.RunMode == DevMode {
		fmt.Printf("%s\n", logo)
	}

	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd@%d", fd))
	listener, err := net.FileListener(f)
	if err != nil {
		return
	}
	defer listener.Close()
	err = http.Serve(listener, engine)

	return
}

// ServeHTTP implements the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()
	engine.handleHTTPRequest(c)
	engine.pool.Put(c)
}

// Match routes from engine.trees
func (engine *Engine) Match(method, path string, match *matchValue) bool {
	path = cleanPath(path, engine.ignoreTrailingSlash)
	for _, tree := range engine.trees {
		routes := tree.get(method)
		for _, route := range routes {
			if route.Match(path, match) {
				return true
			}
		}
	}
	return false
}

// useMiddleware use configuration-related middleware
func (engine *Engine) useMiddleware() {
	// init session and use middleware
	if engine.sessionOn {
		InitSession()
		engine.Use(Session())
	}

	// use gzip middleware
	if engine.gzipOn {
		engine.Use(Gzip(DefaultCompression))
	}
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	var match matchValue
	if engine.Match(c.Request.Method, c.Request.URL.Path, &match) {
		//TODO: cache the c.params
		if match.params != nil {
			c.Params = *match.params
		}
		c.handlers = match.handlers
		c.fullPath = match.fullPath
		c.Next()
		c.writermem.WriteHeaderNow()
		return
	}

	c.handlers = engine.allNoRoute

	serveError(c, http.StatusNotFound, default404Body)
}

// Use middleware
//	ex: r.Use(Auth())
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}

// NoRoute set 404 handler
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

// RouteMapInfo route map info
type RouteMapInfo struct {
	Method  string
	Path    string
	Handler string
}

// PrintRouteMap print route map
func (engine *Engine) PrintRouteMap() {
	routeMap := engine.RouteMap()
	for _, item := range routeMap {
		fmt.Printf(" %6s  %20s %5s %s \n", item.Method, item.Path, " ", item.Handler)
	}
}

// RouteMap return  gow route
func (engine *Engine) RouteMap() []*RouteMapInfo {
	rm := make([]*RouteMapInfo, 0)
	for _, t := range engine.trees {
		for _, r := range t.routes {
			rm = append(rm, &RouteMapInfo{
				Method:  t.method,
				Path:    r.path,
				Handler: nameOfFunction(r.handlers[len(r.handlers)-1]),
			})
		}
	}
	return rm
}

// addRoute add route to engine.trees
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	mt := methodTree{
		method: method,
	}
	rc := &routeConfig{
		ignoreCase:          engine.ignoreCase,
		ignoreTrailingSlash: engine.ignoreTrailingSlash,
	}
	mt.routes = mt.addRoute(path, handlers, rc)
	engine.trees = append(engine.trees, mt)
}

func (engine *Engine) allocateContext() *Context {
	v := make(Params, 0, engine.maxParams)
	return &Context{engine: engine, params: &v}
}

var (
	mimePlain = []string{ContentPlain}
	mimeHTML  = []string{ContentHTML}
)

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = mimePlain
		_, err := c.Writer.Write(defaultMessage)
		if err != nil {
			debugPrint("cannot write message to writer during serve error: %v", err)
		}
		return
	}
	c.writermem.WriteHeaderNow()
}

func default404Handler(c *Context) {
	c.writermem.status = http.StatusNotFound
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == http.StatusNotFound {
		c.writermem.Header()["Content-Type"] = mimeHTML
		default404Page = strings.ReplaceAll(default404Page, "{version}", version)
		_, err := c.Writer.Write([]byte(default404Page))
		if err != nil {
			debugPrint("cannot write message to writer during serve error: %v", err)
		}
		return
	}

	c.writermem.WriteHeaderNow()
}
