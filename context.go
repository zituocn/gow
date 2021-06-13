package gow

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	abortIndex int8 = math.MaxInt8 / 2

	ContentJSON              = "application/json; charset=utf-8"
	ContentHTML              = "text/html; charset=utf-8"
	ContentJavaScript        = "application/javascript; charset=utf-8"
	ContentXML               = "application/xml; charset=utf-8"
	ContentXML2              = "text/xml; charset=utf-8"
	ContentPlain             = "text/plain; charset=utf-8"
	ContentPOSTForm          = "application/x-www-form-urlencoded"
	ContentMultipartPOSTForm = "multipart/form-data"
	ContentPROTOBUF          = "application/x-protobuf"
	ContentMSGPACK           = "application/x-msgpack"
	ContentMSGPACK2          = "application/msgpack"
	ContentYAML              = "application/x-yaml; charset=utf-8"
	ContentDownload          = "application/octet-stream; charset=utf-8"
)

// Context gow context
type Context struct {
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int8
	fullPath string

	engine *Engine
	params *Params

	mu   sync.RWMutex
	Keys map[string]interface{}

	// Data html template render Data
	Data map[interface{}]interface{}

	Pager *Pager

	Errors errorMsgs

	sameSite http.SameSite
}

//reset reset Context
func (c *Context) reset() {
	c.Writer = &c.writermem
	c.Params = c.Params[0:0]
	c.handlers = nil
	c.index = -1
	c.fullPath = ""
	c.Keys = nil
	c.Errors = c.Errors[0:0]
	c.Data = make(map[interface{}]interface{}, 0)
	//c.Pager = nil
	*c.params = (*c.params)[0:0]
}

// Handler returns the main handler.
func (c *Context) Handler() HandlerFunc {
	return c.handlers.Last()
}

// HandlerName last handler name
func (c *Context) HandlerName() string {
	return nameOfFunction(c.handlers.Last())
}

// HandlerNames return []string
func (c *Context) HandlerNames() []string {
	hn := make([]string, 0, len(c.handlers))
	for _, val := range c.handlers {
		hn = append(hn, nameOfFunction(val))
	}
	return hn
}

// Next c.Next method
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// IsProd return running in production mode
func (c *Context) IsProd() bool {
	return c.engine.RunMode == ProdMode
}

// IsAborted return is abort
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort abort handler
func (c *Context) Abort() {
	c.index = abortIndex
}

// StopRun stop run handler
func (c *Context) StopRun() {
	panic(stopRun)
}

// AbortWithStatus abort and write status code
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Writer.WriteHeaderNow()
	c.Abort()
}

// Error set error to c.Errors
func (c *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

// AbortWithError abort and error
func (c *Context) AbortWithError(code int, err error) *Error {
	c.AbortWithStatus(code)
	return c.Error(err)
}

/*
 Header
*/

// Header set response header
//	c.Header("Server","gow")
func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

/*
INPUT DATA
REQUEST
*/

// GetIP return k8s Cluster ip
//	default 10.10.10.2
func (c *Context) GetIP() (ip string) {
	//header传递传递的IP
	ip = c.GetHeader("ip")
	if ip == "" {
		ip = c.GetHeader("X-Original-Forwarded-For")
	}
	if ip == "" {
		ip = c.GetHeader("Remote-Host")
	}
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	if ip == "" {
		ip = "10.10.10.2"
	}

	ips := strings.Split(ip, ",")
	if len(ips) > 0 {
		ip = ips[0]
	}

	return strings.TrimSpace(ip)
}

// ClientIP return client ip
func (c *Context) ClientIP() (ip string) {
	addr := c.Request.RemoteAddr
	str := strings.Split(addr, ":")
	if len(str) > 1 {
		ip = str[0]
	}
	return
}

// DecodeJSONBody Unmarshal request body to v
//	return error
func (c *Context) DecodeJSONBody(v interface{}) error {
	body := c.Body()
	return json.Unmarshal(body, &v)
}

// Body return request body -> []byte
func (c *Context) Body() []byte {
	if c.Request.Body == nil {
		return []byte{}
	}
	var body []byte
	body, _ = ioutil.ReadAll(c.Request.Body)

	c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}

// UserAgent return request user agent
func (c *Context) UserAgent() string {
	return c.GetHeader("User-Agent")
}

// IsAjax return is ajax request
//	return X-Requested-With==XMLHttpRequest
func (c *Context) IsAjax() bool {
	return c.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

// Referer return request referer
func (c *Context) Referer() string {
	return c.Request.Referer()
}

// Host return request host
func (c *Context) Host() string {
	return c.Request.Host
}

// IsWebsocket return is websocket request
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.GetHeader("Connection")), "upgrade") &&
		strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
		return true
	}
	return false
}

// IsWeChat return is wechat request
func (c *Context) IsWeChat() bool {
	return strings.Contains(strings.ToLower(c.UserAgent()), strings.ToLower("MicroMessenger"))
}

// Param return the value of the URL param.
func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}

// ParamInt  return the value of the URL param
func (c *Context) ParamInt(key string) (int, error) {
	v := c.Param(key)
	return strconv.Atoi(v)
}

// ParamInt64  return the value of the URL param
func (c *Context) ParamInt64(key string) (int64, error) {
	v := c.Param(key)
	return strconv.ParseInt(v, 10, 64)
}


// Query return query string
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// Form return request.FormValue key
func (c *Context) Form(key string) string {
	return c.Request.FormValue(key)
}

// input
func (c *Context) input() url.Values {
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	return c.Request.Form
}

func (c *Context) formValue(key string) string {
	if v := c.Form(key); v != "" {
		return v
	}
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	return c.Request.Form.Get(key)
}


func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

// GetString 按key返回字串值，可以设置default值
func (c *Context) GetString(key string, def ...string) string {
	if v := c.formValue(key); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetStrings return []string
func (c *Context) GetStrings(key string, def ...[]string) []string {
	var defaultDef []string
	if len(def) > 0 {
		defaultDef = def[0]
	}

	if v := c.input(); v == nil {
		return defaultDef
	} else if kv := v[key]; len(kv) > 0 {
		return kv
	}
	return defaultDef
}

// GetInt return int
func (c *Context) GetInt(key string, def ...int) (int, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(v)
}

// GetInt8 GetInt8
//	-128~127
func (c *Context) GetInt8(key string, def ...int8) (int8, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(v, 10, 8)
	return int8(i64), err
}

//GetUint8 GetUint8
//	0~255
func (c *Context) GetUint8(key string, def ...uint8) (uint8, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseUint(v, 10, 8)
	return uint8(i64), err
}

// GetInt16 GetInt16
//	-32768~32767
func (c *Context) GetInt16(key string, def ...int16) (int16, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(v, 10, 16)
	return int16(i64), err
}

// GetUint16 GetUint16
//	0~65535
func (c *Context) GetUint16(key string, def ...uint16) (uint16, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseUint(v, 10, 16)
	return uint16(i64), err
}

//GetInt32 GetInt32
//	-2147483648~2147483647
func (c *Context) GetInt32(key string, def ...int32) (int32, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(v, 10, 32)
	return int32(i64), err
}

// GetUint32 GetUint32
//	0~4294967295
func (c *Context) GetUint32(key string, def ...uint32) (uint32, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseUint(v, 10, 32)
	return uint32(i64), err
}

// GetInt64 GetInt64
//	-9223372036854775808~9223372036854775807
func (c *Context) GetInt64(key string, def ...int64) (int64, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseInt(v, 10, 64)
}

// GetUint64 GetUint64
//	0~18446744073709551615
func (c *Context) GetUint64(key string, def ...uint64) (uint64, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseUint(v, 10, 64)
	return i64, err
}

// GetFloat64 GetFloat64
func (c *Context) GetFloat64(key string, def ...float64) (float64, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseFloat(v, 64)
}

// GetBool GetBool
func (c *Context) GetBool(key string, def ...bool) (bool, error) {
	v := c.formValue(key)
	if len(v) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseBool(v)
}

/*
Response
*/

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// Redirect http redirect
//	c.Redirect(301,url)
//	c.Redirect(302,url)
func (c *Context) Redirect(code int, url string) {
	c.Writer.WriteHeader(code)
	http.Redirect(c.Writer, c.Request, url, code)
}

// ServerString response text message
//	c.ServerString(200,"success")
//	c.ServerString(404,"page not found")
func (c *Context) ServerString(code int, msg string) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Writer.Header().Set("Content-Type", ContentPlain)
	c.Status(code)
	c.Writer.Write([]byte(msg))
}

// String response text message
func (c *Context) String(msg string) {
	c.ServerString(http.StatusOK, msg)
}

// ServerYAML response yaml data
//	c.ServerYAML(200,yamlData)
func (c *Context) ServerYAML(code int, data interface{}) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Header("Content-Type", ContentYAML)
	c.Status(code)

	bs, err := yaml.Marshal(data)
	if err != nil {
		c.Header("Content-Type", "")
		c.ServerString(http.StatusServiceUnavailable, err.Error())
	}
	c.Writer.Write(bs)
}

// YAML response yaml data
func (c *Context) YAML(data interface{}) {
	c.ServerYAML(http.StatusOK, data)
}

// ServerJSON response JSON data
//	c.ServerJSON(200,"success")
//	c.ServerJSON(404,structData)
//	c.ServerJSON(404,mapData)
func (c *Context) ServerJSON(code int, data interface{}) {
	if code < 0 {
		code = http.StatusOK
	}

	c.Header("Content-Type", ContentJSON)
	c.Status(code)

	encoder := json.NewEncoder(c.Writer)

	if c.engine.RunMode == DevMode {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(data); err != nil {
		c.Header("Content-Type", "")
		c.ServerString(http.StatusServiceUnavailable, err.Error())
	}
}

// JSON response JSON data
func (c *Context) JSON(data interface{}) {
	c.ServerJSON(http.StatusOK, data)
}

// ServerJSONP write data by jsonp format
func (c *Context) ServerJSONP(code int, callback string, data interface{}) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Header("Content-Type", ContentJavaScript)
	c.Status(code)

	bytes, err := json.Marshal(data)
	if err != nil {
		c.Header("Content-Type", "")
		c.ServerString(http.StatusServiceUnavailable, err.Error())
	}
	c.Writer.Write([]byte(callback + "("))
	c.Writer.Write(bytes)
	c.Writer.Write([]byte(");"))
}

// JSONP write date by jsonp format
func (c *Context) JSONP(callback string, data interface{}) {
	c.ServerJSONP(http.StatusOK, callback, data)
}

// ServerXML response xml data
func (c *Context) ServerXML(code int, data interface{}) {
	if code < 0 {
		code = http.StatusOK
	}
	c.Header("Content-Type", ContentXML)
	c.Status(code)
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		c.Header("Content-Type", "")
		c.ServerString(http.StatusServiceUnavailable, err.Error())
	}
}

// XML  response xml data
func (c *Context) XML(data interface{}) {
	c.ServerXML(http.StatusOK, data)
}

// Render html render
func (c *Context) Render(code int, name string, data interface{}) {
	c.Writer.WriteHeader(code)
	if !bodyAllowedForStatus(code) {
		c.engine.Render.WriteContentType(c.Writer)
		c.writermem.WriteHeader(code)
		return
	}

	if err := c.engine.Render.Render(c.Writer, name, data); err != nil {
		debugPrint("html render error: #s", err.Error())
	}
}

// ServerHTML html page render
//	c.ServerHTML(200,"index.html")
//	c.ServerHTML(200,"admin/login.html")
//	c.ServerHTML(404,"error.html")
func (c *Context) ServerHTML(code int, name string, data ...interface{}) {
	if !c.engine.AutoRender {
		c.ServerString(http.StatusNotFound, string(default404Body))
		return
	}
	var v interface{}
	if len(data) > 0 {
		v = data[0]
	} else {
		v = c.Data
	}
	c.Render(code, name, v)
}

// HTML html page render
//	c.HTML("index.html")
//	c.HTML("login.html",data)
func (c *Context) HTML(name string, data ...interface{}) {
	if len(data) > 0 {
		v := data[0]
		c.ServerHTML(http.StatusOK, name, v)
		return
	}
	c.ServerHTML(http.StatusOK, name)
}

/*
COOKIE
*/

// SetCookie set cookie
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: c.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// GetCookie get cookie
func (c *Context) GetCookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

/*
UPLOAD
*/

// GetFile get single file from request
func (c *Context) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
			return nil, nil, err
		}
	}
	return c.Request.FormFile(key)
}

// GetFiles get files from request
func (c *Context) GetFiles(key string) ([]*multipart.FileHeader, error) {
	if files, ok := c.Request.MultipartForm.File[key]; ok {
		return files, nil
	}
	return nil, http.ErrMissingFile
}

// SaveToFile upload the file and save it on the server.
//	c.SaveToFile("file","./upload/1.jpg")
func (c *Context) SaveToFile(fromFile, toFile string) error {
	file, _, err := c.Request.FormFile(fromFile)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(toFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

/*
DOWNLOAD
*/

// FileAttachment writes the specified file into the body stream in an efficient way
// On the client side, the file will typically be downloaded with the given filename
func (c *Context) FileAttachment(filepath, filename string) {
	c.Header("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(c.Writer, c.Request, filepath)
}

// Download download data
func (c *Context) Download(data []byte) {
	c.Header("Content-Type", ContentDownload)
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(data)
}

// DownLoadFile download data to filename
//	c.DownLoadFile(data,"table.xlsx")
func (c *Context) DownLoadFile(data []byte, filename string) {
	c.Header("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Download(data)
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
