# Upload && Download

gow 支持文件的上传及下载


## 文件上传

### 获取单上文件 


```go
func (c *Context) GetFile(key string) (multipart.File, *multipart.FileHeader, error) 

```

### 多个文件上传

```go
func (c *Context) GetFiles(key string) ([]*multipart.FileHeader, error)
```

### 上传并保存到服务器

```go
func (c *Context) SaveToFile(fromFile, toFile string) error 
```

演示代码

```go
package main

import (
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()
    r.POST("/v1/upload", UploadHandler)
    r.Run()
}

// UploadHandler 文件上传
//  上传文件后，在服务器保存为 test.jpg
func UploadHandler(c *gow.Context) {
    err:=c.SaveToFile("file","test.png")
    if err!=nil{
        c.JSON(gow.H{
            "code":1,
            "msg":err.Error(),
        })
        return
    }

    c.JSON(gow.H{
        "code":0,
        "msg":"success",
    })
}

```

```sh
curl http://127.0.0.1:8080/v1/upload -F "file=@1.png" -v

*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> POST /v1/upload HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.64.1
> Accept: */*
> Content-Length: 66401
> Content-Type: multipart/form-data; boundary=------------------------9c2cc5ed0cbd9621
> Expect: 100-continue
>
< HTTP/1.1 100 Continue
* We are completely uploaded and fine
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Fri, 11 Jun 2021 09:11:57 GMT
< Content-Length: 36
<
{
  "code": 0,
  "msg": "success"
}
* Connection #0 to host 127.0.0.1 left intact
* Closing connection 0

```

## 下载文件


```go

// FileAttachment 下载服务器上已经存在的文件
func (c *Context) FileAttachment(filepath, filename string)
```


```go

// DownLoadFile 下载 []byte 到文件，也可以实现拉取远程文件后再下载
func (c *Context) DownLoadFile(data []byte, filename string) 
```

演示代码

```go

package main

import (
    "github.com/zituocn/gow"
)

func main() {
    r := gow.Default()
    r.GET("/v1/down1", Down1Handler)
    r.GET("/v1/down2", Down2Handler)
    r.Run()
}


// Down1Handler 文件下载
func Down1Handler(c *gow.Context) {
    // 下载服务器已存在的文件
    c.FileAttachment("test.png", "new.png")
}

// Down2Handler 文件下载
func Down2Handler(c *gow.Context) {

    // 下载 []byte 到文件
    //  也可以实现拉取远程文件再下载
    c.DownLoadFile([]byte("hello gow!"), "hello.txt")
}

```

浏览器访问:

```
http://127.0.0.1:8080/v1/down1
http://127.0.0.1:8080/v1/down2
```