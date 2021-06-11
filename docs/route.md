# Route

gow 需要在代码中，手动配置路由来实现不同的访问动词和地址

## 基本思路

一个http请求

```
用户 -> request /v1/user/1 -> response -> 用户
```

在 `gow` 中

```
用户 -> request -> 匹配路由 -> 找到处理程序(handler) -> 执行处理程序 -> response -> 用户
```


## 标准动词

```go
GET POST PUT DELETE OPTIONS HEAD PATCH Any
```

### 可自定义

```go
func (group *RouterGroup) Handle(httpMethod, path string, handlers ...HandlerFunc)
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

默认情况下，路由已忽略字母大小写：

```go
func main(){
    r:=gow.Default()
    r.GET("/User/1",UserHandler) 
    r.Run(8080)
}
```

以下访问也能匹配

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


可以通过以下方法来设置不忽略大小写

```go
package main

func main(){
    r:=gow.Default()
    // 不忽略大小写
    r.SetIgnoreCase(false)  
    r.GET("/User/1",UserHandler) 
    r.Run(8080)
}
```

```sh
curl -i http://127.0.0.1:8080/user/1

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
Date: Fri, 11 Jun 2021 04:49:22 GMT
Content-Length: 18

404 page not found

```

也可以在统一配置文件中设置开关 (可参考: [统一配置文件]())

```sh
ignore_case = false
```

