# Website 

gow 支持 HTML 模板服务器端渲染，使用 `golang原生` 的模板语法

## *感谢*
gow 的 HTML模板 代码，参照了 `beego` 中的处理方法


## 渲染 HTML 页面

### 主要方法

```go
// HTML 200 状态码
func (c *Context) HTML(name string, data ...interface{})

// ServerHTML 可自定义状态码
func (c *Context) ServerHTML(code int, name string, data ...interface{})
```


### 演示代码

项目结构

```sh
PROJECT_NAME
├──conf
     app.conf
├──static
      ├── img
            ├──111.png
      ├──js
      ├──css
├──views
    ├──index.html
    ├──article
        ├──detail.html
├──main.go
```

```sh
vi conf/app.conf

app_name = test         # 应用名称
run_mode = dev          # 运行模式，支持：`dev` `prod` `test` 三种模式
http_addr = 8080        # http 监听地址
auto_render = true      # 是否自动渲染HTML模板
views = views           # HTML模板目录
template_left = {{      # golang 模板左符号
template_left = }}      # golang 模板右符号 
session_on = false      # 是否打开session
gzip_on = true          # 是否打开gzip 
ignore_case = true      # 是否忽略大小写
```

main.go

```go
package main

import (
    "github.com/zituocn/gow"
    "time"
)

func main() {
    r := gow.Default()

    // 使用自动配置，读取了 conf/app.conf中的配置
    r.SetAppConfig(gow.GetAppConfig())

    // 设置资源资源的路由
    r.Static("/static","static")
    r.GET("/",IndexHandler)
    r.Run()
}


// IndexHandler index page
func IndexHandler(c *gow.Context){
    c.Data["datetime"] = time.Now()
    c.Data["title"] = "这是一个gow的HTML页面"
    c.HTML("index.html")
}
```

views/index.html

```sh
<html>
<head>
    <title>gow html page</title>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width,initial-scale=1.0,minimum-scale=1.0,maximum-scale=1.0,user-scalable=0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
</head>
<body>
<div style="height: 100vh;margin:2rem auto;text-align: center;">
    <h2>{{.title}}</h2>
    <img src="/static/img/111.png" style="max-width: 100%;;">
    <br />
    <img src="/static/img/222.png" style="max-width:100%;margin-bottom: 2rem;">
    <p style="color:#999;">
        <i>{{.datetime}}</i>
    </p>
</div>
</body>
</html>

```

浏览器访问

```sh
http://127.0.0.1:8080/
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