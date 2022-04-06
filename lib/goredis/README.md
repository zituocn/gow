# goredis

使用 github.com/go-redis/redis/v8 完成的redis封装


## 使用

更多redis方法封装，请查看 package：

```go
https://github.com/go-redis/redis
```

官方文档：

```go
https://redis.uptrace.dev/
```

### 单实例的情况

```go
package main

import (
	"context"
	"github.com/zituocn/gow/lib/goredis"
	"log"
	"time"
)

// init redis
func init() {
	err := goredis.InitDefaultDB(&goredis.RedisConfig{
		Host:     "192.168.0.197",
		Port:     6379,
		Pool:     100,
		Password: "",
		Name: "test",
		DB:       0,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	ctx := context.Background()

	// get rdb
	rdb := goredis.GetRDB()

	// string set
	val, err := rdb.Set(ctx, "key", "redis test", 10*time.Minute).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("val = ", val)

	// string get
	result, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("result = ", result)

}

```

### 多个实例的情况

```go
package main

import (
	"context"
	"github.com/zituocn/gow/lib/goredis"
	"log"
	"time"
)

func init() {
	rcs := make([]*goredis.RedisConfig, 0)

	rcs = append(rcs, &goredis.RedisConfig{
		Name:     "test-1",
		Host:     "192.168.0.197",
		Port:     6379,
		Pool:     100,
		Password: "123456",
		DB:       0,
	})

	rcs = append(rcs, &goredis.RedisConfig{
		Name:     "test-2",
		Host:     "192.168.0.197",
		Port:     6379,
		Pool:     100,
		Password: "123456",
		DB:       13,
	})

	err := goredis.InitDB(rcs)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	ctx := context.Background()

	// get rdb by name
	rdb := goredis.GetRDBByName("test-2")

	// string set
	val, err := rdb.Set(ctx, "key", "redis test", 10*time.Minute).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("val = ", val)

	// string get
	result, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("result = ", result)
}

```