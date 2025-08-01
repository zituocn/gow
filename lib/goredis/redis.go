/*
go redis的简单封装
可实现多个实例的连接
使用库： github.com/redis/go-redis/v9
sam
2022-04-01
*/

package goredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zituocn/logx"
	"time"
)

var (
	dbs           map[string]*redis.Client
	defaultDBName string
	ctx           = context.Background()
)

// init 初始化dbs的map，后续不再重置
func init() {
	dbs = make(map[string]*redis.Client)
}

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
	newRedis(db)
	return
}

// GetDBNames 返回所有链接
func GetDBNames() map[string]*redis.Client {
	return dbs
}

// GetRDB returns a *redis.Client
func GetRDB() *redis.Client {
	m, ok := dbs[defaultDBName]
	if !ok {
		logx.Panic("[redis] 未init，请参照使用说明")
	}
	return m
}

// InitDB init multiple rdb to map
func InitDB(list []*RedisConfig) (err error) {
	if len(list) == 0 {
		err = fmt.Errorf("[redis] 没有需要init的DB")
		return
	}
	for _, item := range list {
		newRedis(item)
	}

	return
}

// GetRDBByName get rdb by name
func GetRDBByName(name string) *redis.Client {
	m, ok := dbs[name]
	if !ok {
		logx.Panic("[redis] 未init，请参照使用说明")
	}
	return m
}

// GetDBS get dbs from map
func GetDBS() map[string]*redis.Client {
	return dbs
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
		logx.Panicf("[redis]-[%s] 配置信息获取失败", rc.Name)
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
		DialTimeout:  5 * time.Second,
		MaxRetries:   -1,
		MinIdleConns: 10,
	}

	rdb = redis.NewClient(opt)

	// COMMAND ping
	for _, err := rdb.Ping(ctx).Result(); err != nil; {
		logx.Errorf("[redis]-%s 连接异常: %v", rc.string(), err)
		time.Sleep(5 * time.Second)
	}

	dbs[rc.Name] = rdb
}
