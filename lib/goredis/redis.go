/*
go redis的简单封装
可实现多个实例的连接
使用库： github.com/go-redis/redis/v8
sam
2022-04-01
*/

package goredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zituocn/gow/lib/logy"
	"time"
)

var (
	dbs           map[string]*redis.Client
	defaultDBName string
	ctx           = context.Background()
)

// RedisConfig redis 连接配置
type RedisConfig struct {
	Name     string //连接名
	DB       int    //db序号
	Host     string //主机
	Port     int    //端口
	Password string //密码
	Pool     int    //连接池大小
}

// InitDefaultDB init a rdb to map
func InitDefaultDB(db *RedisConfig) (err error) {
	if db == nil {
		err = fmt.Errorf("[redis] 没有需要init的连接")
		return
	}
	defaultDBName = db.Name
	dbs = make(map[string]*redis.Client, 1)
	newRedis(db)
	return
}

// GetRDB returns a *redis.Client
func GetRDB() *redis.Client {
	m, ok := dbs[defaultDBName]
	if !ok {
		logy.Panic("[redis] 未init，请参照使用说明")
	}
	return m
}

// InitDB init multiple rdb to map
func InitDB(list []*RedisConfig) (err error) {
	if len(list) == 0 {
		err = fmt.Errorf("[redis] 没有需要init的DB")
		return
	}
	dbs = make(map[string]*redis.Client, len(list))
	for _, item := range list {
		newRedis(item)
	}

	return
}

// GetRDBByName get rdb by name
func GetRDBByName(name string) *redis.Client {
	m, ok := dbs[name]
	if !ok {
		logy.Panic("[redis] 未init，请参照使用说明")
	}
	return m
}

/*
private
*/

func (m *RedisConfig) string() string {
	return fmt.Sprintf("redis://%s:%s@%s:%d/%d", m.Name, m.Password, m.Host, m.Port, m.DB)
}

// newRedis use redisConfig make dbs
func newRedis(rc *RedisConfig) {
	var (
		rdb *redis.Client
	)
	if rc.Host == "" || rc.Port == 0 || rc.Name == "" {
		logy.Panicf("[redis]-[%s] 配置信息获取失败", rc.Name)
		return
	}
	if rc.DB < 0 {
		rc.DB = 0
	}
	if rc.Pool < 0 {
		rc.Pool = 10
	}
	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", rc.Host, rc.Port),
		Password:     rc.Password,
		DB:           rc.DB,
		PoolSize:     rc.Pool,
		IdleTimeout:  30 * time.Second,
		DialTimeout:  5 * time.Second,
		MaxRetries:   -1,
		MinIdleConns: 10,
	}

	rdb = redis.NewClient(opt)

	// COMMAND ping
	for _, err := rdb.Ping(ctx).Result(); err != nil; {
		logy.Errorf("[redis]-%s 连接异常: %v", rc.string(), err)
		time.Sleep(5 * time.Second)
	}

	dbs[rc.Name] = rdb
}
