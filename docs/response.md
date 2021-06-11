# Response

gow 支持多种常见格式的Response


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


## Redirect

```go
func (c *Context) Redirect(code int, url string)
```
```go
func GetUser(c *gow.Context) {

    // 永久重定向
    c.Redirect(301,"https://github.com/zituocn/gow")  

    // 临时重定向    
    // c.Redirect(302,"https://github.com/zituocn/gow")
}
```