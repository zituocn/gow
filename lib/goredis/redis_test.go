package goredis

import (
	"context"
	"log"
	"testing"
	"time"
)

func init() {
	err := InitDefaultDB(&RedisConfig{
		Name:     "test",
		Host:     "127.0.0.1",
		Port:     6379,
		Pool:     1000,
		Password: "2236236",
		DB:       3,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func Benchmark_test(b *testing.B) {
	c := context.Background()
	rdb := GetRDB()
	for i := 0; i < b.N; i++ {
		rdb.Set(c, "key:123", "test value", 10*time.Second)
	}
}
