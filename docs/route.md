# Route

gow 需要在代码中，手动书写路由来实现不同的访问动词和地址


## *感谢*

gow 中的 route 实现机制，融合了 gin 和 mux .

## 基本思路

一个http请求

```
用户 -> request /v1/user/1 -> response -> 用户
```

在 `gow` 中

```
用户 -> request -> 匹配路由 -> 找到处理程序(handler) -> 执行处理程序 -> response -> 用户
```


## 访问动词

```go
GET POST PUT DELETE OPTIONS HEAD PATCH Any
```

### 可自定义

```go
func (group *RouterGroup) Handle(httpMethod, path string, handlers ...HandlerFunc)
```

使用 `r.Handle` 匹配 GET请求

```go
package main

import "github.com/zituocn/gow"

func main() {
	r := gow.Default()
	r.Handle("GET", "/", func(c *gow.Context) {
		c.String("index")
	})
}
```

使用 `r.Handle` 同时匹配 GET和POST方法

```go
package main

import "github.com/zituocn/gow"

func main() {
	r := gow.Default()
	r.Handle("GET,POST", "/", func(c *gow.Context) {
		c.String("index")
	})
}
```

## 基础路由

```go
func main(){

    r:=gow.Default()

    // GET /
    r.GET("/",IndexHandler)

    // GET /user/1
    r.GET("/user/1",UserHandler) 

    // POST /user/1
    r.POST("/user/1",UserUpdateHandler)  

    // DELETE /user/1
    r.DELETE("/user/1",UserDeleteHandler)

    // PUT /user/1
    r.PUT("/user/1",UserPutHandler)

    // /article/100 
    r.GET("/article/{article_id}",ArticlDetailHandler)

    // /read-100.html
    r.GET("/read-{id}.html",ReadHandler)

    // 静态资源路由
    r.Static("/static","static")


    r.Run(8080)
}

```

## 正则(参数)路由


```go
func main() {
    r := gow.Default()
    r.SetIgnoreCase(false)
    r.GET("/{name}/{id}", UserHandler)
    r.Run(8080)
}

// UserHandler get user info
func UserHandler(c *gow.Context) {
    name := c.Param("name")
    id, _ := c.ParamInt64("id")
    c.JSON(gow.H{
        "name": name,
        "id":   id,
    })
}
```

```sh 
curl -i http://127.0.0.1:8080/wahaha/100
curl -i http://127.0.0.1:8080/test/101
```

使用 `c.Param()` 获取参数

```go
func (c *Context) Param(key string) string 
```

```go
func (c *Context) ParamInt(key string) (int, error) 
```

```go
func (c *Context) ParamInt64(key string) (int64, error) 
```


## 匹配所有

使用 `{match_all}` 匹配 前缀 `/v1/` 的所有请求

```go
package main

import (
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()
    r.SetIgnoreCase(false)
    r.GET("/v1/{match_all}", UserHandler)
    r.Run(8080)
}

// UserHandler get user info
func UserHandler(c *gow.Context) {
    matchAll := c.Param("match_all")
    c.JSON(gow.H{
        "match_all":matchAll,
    })
}

```

```sh
curl -i http://127.0.0.1:8080/v1/user/page

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 11 Jun 2021 06:22:11 GMT
Content-Length: 31

{
  "match_all": "user/page"
}
```

## 路由分组

使用 `r.Group()` 来完成

```go
func (group *RouterGroup) Group(path string, handlers ...HandlerFunc)
```


演示代码

```go

func main() {
    r := gow.Default()

    v1 := r.Group("/v1")
    {
        v1.GET("/user", handler)
        v1.DELETE("/user", handler)
        v1.POST("/user", handler)
        v1.GET("/user/page", handler)
    }

    v2 := r.Group("/v2")
    {
        v2.GET("/user", handler)
        v2.DELETE("/user", handler)
        v2.POST("/user", handler)
        v2.GET("/user/page", handler)
    }

    r.Run(8080)
}
```


## 路由大小写

默认情况下，路由时忽略字母大小写：

```go
func main(){
    r:=gow.Default()
    r.GET("/User/1",UserHandler) 
    r.Run(8080)
}
```

```sh
curl -i http://127.0.0.1:8080/user/1 

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 11 Jun 2021 04:47:05 GMT
Content-Length: 94

{
  "city": "成都市",
  "prov": "四川省",
  "uid": 1,
  "username": "新月却泽滨"
}
```

可以通过以下方法来设置忽略路由的大小写

```go
package main

func main(){
    r:=gow.Default()
    // 忽略大小写
    r.SetIgnoreCase(true)  
    r.GET("/User/1",UserHandler) 
    r.Run(8080)
}
```

以下URL不能访问：

```sh
curl -i http://127.0.0.1:8080/user/1

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
Date: Fri, 11 Jun 2021 04:49:22 GMT
Content-Length: 18

404 page not found
```




也可以在统一配置文件中的开关，来实现忽略与否 (可参考: [统一配置文件](https://github.com/zituocn/gow/blob/main/docs/config.md))

```sh
ignore_case = false # 不忽略
ignore_case = true  # 忽略
```

## 路由尾部斜杠 "/"

访问网站时，可能经常会忽略路由地址的 "/"，即有无 "/"，都能访问。但在API或某些特殊环境下，尾部斜杠不能忽略。

可以通过以下访问，选择是否忽略路由的尾部"/"

```go
package main

func main(){
    r:=gow.Default()

    // 设置是否忽略尾部 "/"
    r.SetIgnoreTrailingSlash(true)  
    
    // 不忽略尾部 "/"
    // r.SetIgnoreTrailingSlash(false)

    r.GET("/user/1",UserHandler) 
    r.Run(8080)
}
```

当用户使用以下方式时，都能命中路由

```sh
curl -i http://127.0.0.1:8080/user/1/
```

```sh
curl -i http://127.0.0.1:8080/user/1

```

gow 默认情况下，会忽略地址后的 "/"

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
