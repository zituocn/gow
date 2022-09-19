# redis 使用

此库已弃用，不再维护建议使用：

```
github.com/zituocn/gow/lib/goredis
```


## 初始化链接

> 比如放在 init()方法中，初始化一次

```go
    
func init(){
	conf := RDSConfig{
		Host:      "192.168.0.197",
		Port:      6379,
		Password:  "123456",
		MaxIdle:   50,
		MaxActive: 10000,
		DB:        1,
	}

	err := InitRDSClient(&conf)
	if err != nil {
		fmt.Println("连接失败:", err)
	}
}

```

## 操作redis

```go

	//get rds common
	redis := GetRDSCommon()
	key := "token:1"
        
    // SetString
	_, err = redis.SetString(key, "abcdef1xd0r1jdkf")
	if err != nil {
		t.Fatal(err)
	}

```
