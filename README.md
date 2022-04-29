# gow
gow is a golang HTTP web framework

![gow logo](docs/logo.png)


> 借鉴和参考的项目：gin/beego/mux

## 项目地址

[https://github.com/zituocn/gow](https://github.com/zituocn/gow)

## 特性

* 类 `gin` 的 `Context` 封装、路由分组和 middleware，可快速入门
* 使用 `regexp` 实现路由完全匹配，支持大小写忽略
* 统一的配置入口(ini格式)，也可实现自己喜欢的配置方式
* 支持服务器端渲染HTML页面，可自由扩展HTML模板函数
* 可以自由选择封装在lib的sdk，如 mysql redis nsq rpc mem-cache oauth pay 等

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

## 一些演示代码

可直接运行

* [github.com/zituocn/gow-demo](https://github.com/zituocn/gow-demo)

---

## 2. 更多文档

* [路由详解 && 路由参数 && 路由分组](https://github.com/zituocn/gow/blob/main/docs/route.md)
* [中间件(middleware) 使用](https://github.com/zituocn/gow/blob/main/docs/middleware.md)
* [获取请求值](https://github.com/zituocn/gow/blob/main/docs/request.md)
* [输出值 && JSON / XML / JSONP / YAML](https://github.com/zituocn/gow/blob/main/docs/response.md)
* [统一配置文件](https://github.com/zituocn/gow/blob/main/docs/config.md)
* [做一个网站 && HTML模板使用指南](https://github.com/zituocn/gow/blob/main/docs/website.md)
* [HTML模板函数](https://github.com/zituocn/gow/blob/main/docs/html.md)
* [文件的上传及下载](https://github.com/zituocn/gow/blob/main/docs/upload.md)
* [实现反向代理(new)](https://github.com/zituocn/gow/blob/main/docs/proxy.md)
* [lib 库介绍：logy mysql config ](https://github.com/zituocn/gow/blob/main/docs/lib.md)

## 3. 感谢

* [beego](https://github.com/beego/beego) -> 参考了1.x中的HTML模板设计
* [gin](https://github.com/gin-gonic/gin) -> 参考了 `engine` 和 `Context` 设计
* [mux](https://github.com/gorilla/mux)   -> 参考了 路由设计
* [gorm](https://github.com/go-gorm/gorm) -> 推荐使用 gorm
* [gini](https://github.com/gkzy/gini)     -> 用来操作 `ini` 格式的配置文件

## 4. License

MIT License. See the LICENSE file for details.



