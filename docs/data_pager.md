# DataPager

DataPager() 是一个gow自带的翻页处理 middleware.


## DataPager的结构

```go
// Pager pager struct
type Pager struct {
    Page      int64 `json:"page"`           // 当前码
    Limit     int64 `json:"-"`              // 每页条条
    Offset    int64 `json:"-"`       
    Count     int64 `json:"count"`          // 数据总数
    PageCount int64 `json:"page_count"`     // 总的页数
}
```

## 演示代码

```go
package main

import (
    "fmt"
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()
    r.SetAppConfig(gow.GetAppConfig())

    v1 := r.Group("/v1")

    // use gow.DataPager()
    v1.Use(gow.DataPager())
    {
        v1.GET("/user/page", DataHandler)
    }
    r.Run()
}

// DataHandler handler api
func DataHandler(c *gow.Context) {
    // 取数据
    data, count := PageList(c.Pager.Page, c.Pager.Limit)

    // 设置总条数
    c.Pager.Count = count

    // 输出
    c.DataJSON(data, c.Pager)
}

// User struct
type User struct {
    Uid      int64  `json:"uid"`
    UserName string `json:"username"`
}

// PageList 返回假的数据
func PageList(offset, limit int64) (data []*User, count int64) {
    data = make([]*User, 0)
    for i := 0; int64(i) < limit; i++ {
        data = append(data, &User{
            Uid:      int64(i),
            UserName: fmt.Sprintf("用户名-%d", i),
        })
    }
    count = 10248
    return
}

```

```sh
curl -i http://127.0.0.1:8080/v1/user/page?page=2&limit=5

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 12 Jun 2021 06:39:26 GMT
Content-Length: 525

{
  "code": 0,
  "msg": "success",
  "time": 1623479966,
  "body": {
    "pager": {
      "page": 2,
      "count": 10248,
      "page_count": 2050
    },
    "data": [
      {
        "uid": 0,
        "username": "用户名-0"
      },
      {
        "uid": 1,
        "username": "用户名-1"
      },
      {
        "uid": 2,
        "username": "用户名-2"
      },
      {
        "uid": 3,
        "username": "用户名-3"
      },
      {
        "uid": 4,
        "username": "用户名-4"
      }
    ]
  }
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