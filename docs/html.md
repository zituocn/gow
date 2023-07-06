# HTML template func

gow 支持在HTML模板上，使用内置或自定义的模板函数(template func)


## 使用演示

```go
package main

import (
    "github.com/zituocn/gow"
    "time"
)

func main() {
    r := gow.Default()
    r.SetAppConfig(gow.GetAppConfig())

    r.GET("/", IndexHandler)

    r.Run()
}

var (
    label = `
        <span style="color:#ff3300;">这是一个label标签</span>
    `
)

func IndexHandler(c *gow.Context) {

    // 模板变量传递
    c.Data["label"] = label
    c.Data["now"] = time.Now()
    c.Data["timestamp"] = time.Now().Unix()
    
    c.HTML("index.html")
}

```

index.html

```sh
<!doctype html>
<html>
<head>
    <title>gow html page</title>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width,initial-scale=1.0,minimum-scale=1.0,maximum-scale=1.0">
    <style type="text/css">
        .box{margin:2rem auto;text-align: center;width:88%;}
    </style>
</head>
<body>
<div class="box">
    <p>
        {{str2html .label}}
    </p>
    <p>

        日期：{{date .now}}
        <br />
        时间：{{time .now}}
        <br />
        <br />
        时间戳日期：{{int_date .timestamp}}
        <br />
        自定义格式化：{{int_datetime_format .timestamp "YYYY-MM-DD HH:mm:ss"}}
        <br />
        时间戳时间：{{int_datetime .timestamp}}

        <br />

        JS: {{assets_js "/static/js/main.js"}}
        <br />
        CSS: {{assets_css "/static/css/index.css"}}
        <br />
    </p>
</div>
</body>
</html>
```

执行结果

```sh

curl -i http://127.0.0.1:8080/

HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Date: Sat, 12 Jun 2021 07:03:47 GMT
Content-Length: 898

<!doctype html>
<html>
<head>
    <title>gow html page</title>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width,initial-scale=1.0,minimum-scale=1.0,maximum-scale=1.0">
    <style type="text/css">
        .box{margin:2rem auto;text-align: center;width:88%;}
    </style>
</head>
<body>
<div class="box">
    <p>

        <span style="color:#ff3300;">这是一个label标签</span>

    </p>
    <p>

        日期：2021-06-12
        <br />
        时间：15:03:47
        <br />
        <br />
        时间戳日期：2021-06-12
        <br />
        自定义格式化：2021-06-12 15:03:47
        <br />
        时间戳时间：2021-06-12 15:03:47

        <br />

        JS: <script src="/static/js/main.js"></script>
        <br />
        CSS: <link href="/static/css/index.css" rel="stylesheet" />
        <br />
    </p>
</div>
</body>
</html>%


```


## 自定义模板函数的使用

```go

// AddFuncMap add fn func to template func map
func (engine *Engine) AddFuncMap(key string, fn interface{})
```


```go
package main

import (
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()

    // 添加自定义func
    r.AddFuncMap("hi", hello)
    r.SetAppConfig(gow.GetAppConfig())
    r.GET("/", IndexHandler)
    r.Run()
}

// hello 自定义函数
func hello(str string) string {
    return "hello : " + str
}

func IndexHandler(c *gow.Context) {
    c.HTML("index.html")
}

```


在模板文件上使用

```sh
   {{hi "gow"}}
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