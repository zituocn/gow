# Config

gow 可以通过统一的配置文件，来实现一些自动化配置参数，自动载入。


## 自动化配置载入

项目结构

```sh
PROJECT_NAME
├──conf
    ├──app.conf         # 默认的配置文件 
    ├──prod.app.conf    # 环境变量 GOW_RUN_MODE="prod" 时的配置文件
    ├──dev.app.conf     # 环境变量 GOW_RUN_MODE="dev" 时的配置文件
├──static
      ├── img
            ├──111.jpg
            ├──222.jpg
            ├──333.jpg
      ├──js
      ├──css
├──views
    ├──index.html
    ├──article
        ├──detail.html
├──main.go
```


### 配置文件

```sh
vi conf/app.conf
```

```sh
app_name = test         # 应用名称
run_mode = dev          # 运行模式，支持：`dev` `prod` `test` 三种模式
http_addr = 9090        # http 监听地址
auto_render = true      # 是否自动渲染HTML模板
views = views           # HTML模板目录
template_left = {{      # golang 模板左符号
template_left = }}      # golang 模板右符号 
session_on = false      # 是否打开session
gzip_on = true          # 是否打开gzip 
ignore_case = true      # 是否忽略大小写
```


```sh
vi conf/prod.app.conf
```

```sh
app_name = test         # 应用名称
run_mode = prod         # 运行模式，支持：`dev` `prod` `test` 三种模式
http_addr = 8080        # http 监听地址
auto_render = true      # 是否自动渲染HTML模板
views = views           # HTML模板目录
template_left = {{      # golang 模板左符号
template_left = }}      # golang 模板右符号 
session_on = false      # 是否打开session
gzip_on = true          # 是否打开gzip 
ignore_case = true      # 是否忽略大小写
```


可根据 环境变量 `GOW_RUN_MODE=prod` 实现不同运行模式，自动载入不同的配置文件


```sh
env GOW_RUN_MODE=prod go run main.go
```
```sh
env GOW_RUN_MODE=test go run main.go
```

```sh
env GOW_RUN_MODE=dev go run main.go
```

### 使用方法

* 使用 ` r.GetAppConfig() ` 读取配置
* 使用 ` r.SetAppConfig() ` 载入配置

```go
package main

import (
    "github.com/zituocn/gow"
)

func main(){
    r:=gow.Default()
    r.SetAppConfig(gow.GetAppConfig())
    r.Run()
}
```

*提示* 

你也可以通过 `r.SetAppConfig()` 来实现自己的配置格式载入。

## 更多Config使用指南


### 一个演示

```sh
vi conf/app.conf
```

```sh
app_name = test
run_mode = dev
http_addr = 9090
auto_render = true
views = views
template_left = {{
template_left = }}
session_on = false
gzip_on = true
ignore_case = true

[redis]
host = "192.168.0.197"
port = 6379
db = 0
password = "123456"

[system]
version = v0.1.0
is_redirect = false
pageLimit = 100
```

```go
package main

import (
    "github.com/zituocn/gow"
    "github.com/zituocn/gow/lib/config"
)

func main() {
    r := gow.Default()
    r.GET("/v1/config", ConfigHandler)
    r.Run()
}

// ConfigHandler get config
func ConfigHandler(c *gow.Context) {
    host := config.GetString("redis::host")
    port, _ := config.GetInt("redis::port")
    db, _ := config.GetInt64("redis::db")
    password := config.GetString("redis::password")
    isRedirect, _ := config.GetBool("system::is_redirect")
    
    c.JSON(gow.H{
        "host":        host,
        "port":        port,
        "db":          db,
        "password":    password,
        "is_redirect": isRedirect,
    })
}

```

```sh
curl -i http://127.0.0.1:8080//v1/config 

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 11 Jun 2021 07:44:16 GMT
Content-Length: 105

{
  "db": 0,
  "host": "192.168.0.197",
  "is_redirect": false,
  "password": "123456",
  "port": 6379
}
```

### 更多方法

```go

// 不带默认值
config.GetString()
config.GetInt()
config.GetInt64()
config.GetFloat()
config.GetBool()


// 带默认值 
config.DefaultString()
config.DefaultInt()
config.DefaultInt64()
config.DefaultFloat()
config.DefaultBool()
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
* [lib 库介绍：logy mysql config ](https://github.com/zituocn/gow/blob/main/docs/lib.md)
