# middleware

`middleware` 支持路由分组，即不同的路由分组，可使用不同的 `middleware` .

## 如何使用中间件

```go
r:=gow.Default()
r.Use(...)
```
*或*

```go
r:=gow.New()
r.Use(...)
```

## 使用自带中间件

* Recovery：Recovery()
* 日志： Logger()
* 翻页：DataPager()

```go
r:=gow.Default()
r.Use(gow.DataPager())
...
```
当使用`gow.Default()`时，系统已经默认使用 `Logger` 和 `Recovery()` 两个中间件.


```go
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
```

你也可以添加使用自己设计的`日志` 和 `Recovery`中间件


## 自己设计一个

基本代码：

```go
func FuncName() gow.HandlerFunc {
	return func(c *gow.Context) {
        ...
        ...
        c.Next()  //此处不忘记
	}
}
```

看完整的源码:

```go
package main

import (
	"github.com/zituocn/gow"
)

func main() {
	r := gow.Default()

	v1 := r.Group("/v1")
	
	// /v1 下所有路由，使用APIAuth()
	v1.Use(APIAuth())


	user := v1.Group("/user")

	// /v1/user/ 下所有路由，使用UserAuth()
	user.Use(UserAuth())
	{
		// route: /v1/user/1 
		user.GET("/{uid}", UserHandler)
	}
	r.Run()

}

// APIAuth API接口鉴权
func APIAuth() gow.HandlerFunc {
	return func(c *gow.Context) {
		auth := c.GetHeader("auth")
		if auth != "123456" {
			c.ServerJSON(403, gow.H{
				"code": 403,
				"msg":  "没有API访问权限",
			})
			c.StopRun()
			return
		}
		c.Next()
	}
}

// UserAuth 用户鉴权
func UserAuth() gow.HandlerFunc {
	return func(c *gow.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.ServerJSON(403, gow.H{
				"code": 403,
				"msg":  "此用户无权限",
			})
			c.StopRun()
			return
		}
		c.Next()
	}
}

// UserHandler get user info
func UserHandler(c *gow.Context) {
	h := map[string]interface{}{
		"uid":      1,
		"username": "新月却泽滨",
		"city":     "成都市",
		"prov":     "四川省",
	}
	c.JSON(h)
}

```