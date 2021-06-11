# gow 
gow is a golang HTTP web framework

> 借鉴和参考的项目：gin/beego/mux


### 项目地址

[https://github.comm/zituocn/gow](https://github.comm/zituocn/gow)


## 1. 快速开始

```sh
# 创建一个hello的项目
mkdir hello

cd hello

# 使用go mod
go mod init

# 安装gow

go get github.com/zituocn/gow
```

### 1.1 创建 main.go

```go
package main

import (
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()

    r.GET("/", func(c *gow.Context) {
        c.JSON(gow.H{
            "code": 0,
            "msg":  "success",
        })
    })
    
    //default :8080
    r.Run()
}
```
也可以写成这样

```go
package main

import (
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()
    r.GET("/", IndexHandler)
    //default :8080
    r.Run()  
}


// IndexHandler response h
func IndexHandler(c *gow.Context) {
    h := map[string]interface{}{
        "project": "gow",
        "website": "https://github.com/zituocn/gow",
    }
    c.JSON(h)
}

```

### 1.2 运行

```sh
go run main.go
```

运行结果

```sh
Listening and serving HTTP on http://127.0.0.1:8080
```

### 1.3 访问

*curl访问*

```sh
curl -i http://127.0.0.1:8080
```

请求结果

```sh
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 08 Jun 2021 08:51:25 GMT
Content-Length: 67

{
  "project": "gow",
  "website": "https://github.com/zituocn/gow"
}
```

浏览器访问

```sh
在浏览器访问：http://127.0.0.1:8080
```
---

## 2. 更多文档

* [路由详解 && 路由参数 && 路由分组]()
* [中间件(middleware) 使用]()
* [获取请求值]()
* [输出值 && JSON / XML / JSONP / YAML]()
* [统一配置文件]()
* [做一个网站 && HTML模板使用指南]()
* [HTML模板函数]()
* [文件的上传及下载]()
* [lib 库介绍：logy mysql config ]()



## 3. 感谢

* [beego](https://github.com/beego/beego)
* [gin](https://github.com/gin-gonic/gin)
* [mux](https://github.com/gorilla/mux)
* [gorm](https://github.com/go-gorm/gorm)
* [gini](https://github.com/gkzy/gini)

## 4. License

MIT License. See the LICENSE file for details.



