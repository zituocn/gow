# Request

gow 使用 ` gow.Context ` 中的方法来获取来自用户请求的各种数据。

## 获取 form && query参数

```go
package main

import (
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()
    r.GET("/v1/user", GetUser)
    r.Run()
}

// GetUser get user
func GetUser(c *gow.Context) {
    name := c.GetString("name")
    page, _ := c.GetInt("page")
    limit, _ := c.GetInt64("limit")
    is, _ := c.GetBool("is")
    score, _ := c.GetFloat64("score")

    h := gow.H{
        "name":  name,
        "page":  page,
        "limit": limit,
        "is":    is,
        "score": score,
    }
    c.JSON(h)
}

```


```sh
curl -i http://127.0.0.1:8080/v1/user?name=zituocn&page=2&limit=10&is=true&score=58.6

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 11 Jun 2021 08:02:28 GMT
Content-Length: 83

{
  "is": true,
  "limit": 10,
  "name": "zituocn",
  "page": 2,
  "score": 58.6
}

```

## 获取更多类型参数

* c.GetHeader("Accpet")
* c.UserAgent()
* c.Request.RequestURI
* c.GetIP()
* c.ClientIP()
* c.Referer()
* c.Host()
* c.IsWebsocket()
* c.IsWeChat()
* c.Body()


## Request body

```go

package main

import (
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()
    r.PUT("/v1/user", UpdateUser)
    r.Run()
}

type UserInfo struct {
    Username string `json:"username"`
    Mobile   string `json:"mobile"`
}

// UpdateUser set user
func UpdateUser(c *gow.Context) {
    userInfo := new(UserInfo)

     // 反序列化到 struct
    err := c.DecodeJSONBody(&userInfo) 
    if err != nil {
        c.JSON(gow.H{
            "code": 1,
            "msg":  err.Error(),
        })
        return
    }
    c.JSON(userInfo)
}

```


```sh
curl -i -X PUT -H "Content-Type: application/json" -d '{"username":"zituocn","mobile":"13999998888"}' http://127.0.0.1:8080/v1/user


HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 11 Jun 2021 08:19:27 GMT
Content-Length: 55

{
  "username": "zituocn",
  "mobile": "13999998888"
}
```

## Request Header

```go
func (c *Context) GetHeader(key string) string 
```

```go
// 获取header中传递的 token
token:=c.GetHeader("token")
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