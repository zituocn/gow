# Response

gow 使用 ` gow.Context ` 中的方法来向用户响应不同的数据

```sh
JSON XML Text JSONP YAML File
```

## 输出方式

``` go
// GetUser response user info
func GetUser(c *gow.Context) {

    // 200 状态码
    c.String("hello gow")

    // 带状态码
    // c.ServerString(401,"hello gow")
}

```

## 更多方法

不带状态码，默认 200

* c.JSON()
* c.XML()
* c.JSONP()
* c.YAML()
* c.File("readme.txt")

带状态码

* c.ServerJSON()
* c.ServerXML()
* c.ServerJSONP()
* c.ServerYMAL()
* c.ServerString()

## Response Header

```go
func (c *Context) Header(key, value string)
```

```go
c.Header("server", "gow")
c.Header("X-Powered-By", "gow")
c.Header("X-Gow-Version", "v0.1.0")
```

## Redirect

```go
func (c *Context) Redirect(code int, url string)
```

```go
func GetUser(c *gow.Context) {

// 永久重定向
c.Redirect(301, "https://github.com/zituocn/gow")

// 临时重定向    
// c.Redirect(302,"https://github.com/zituocn/gow")
}
```

## 更多文档

* [路由详解 && 路由参数 && 路由分组](https://github.com/zituocn/gow/blob/main/docs/route.md)
* [中间件(middleware) 使用](https://github.com/zituocn/gow/blob/main/docs/middleware.md)
* [获取请求值](https://github.com/zituocn/gow/blob/main/docs/request.md)
* [输出值 && JSON / XML / JSONP / YAML](https://github.com/zituocn/gow/blob/main/docs/response.md)
* [统一配置文件](https://github.com/zituocn/gow/blob/main/docs/config.md)
* [做一个网站 && HTML模板使用指南](https://github.com/zituocn/gow/blob/main/docs/website.md)
* [HTML模板函数](https://github.com/zituocn/gow/blob/main/docs/html.md)
* [文件的上传及下载](https://github.com/zituocn/gow/blob/main/docs/upload.md)
* [lib 库介绍：logx mysql config ](https://github.com/zituocn/logx)