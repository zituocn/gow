package redis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

// RDSConfig redis 配置结构
type RDSConfig struct {

	// Host 主机
	Host string

	// Port 端口
	Port int

	// Password 密码，无密码传""
	Password  string
	MaxIdle   int
	MaxActive int

	// DB 使用的DB
	DB int
}

var cmd = new(RDSCommon)

// InitRDSClient use RDSConfig init redis pool
func InitRDSClient(rdc *RDSConfig) (err error) {
	if rdc == nil {
		err = fmt.Errorf("[RDS] 没有需要init的redis")
		return
	}
	if rdc.Host == "" || rdc.Port == 0 {
		err = fmt.Errorf("[RDS] 没有配置主机或端口")
		return
	}

	cmd.NewClient(newRedisPools(rdc))

	return
}

// ====================private====================

// newRedisPools
func newRedisPools(rdc *RDSConfig) *redis.Pool {
	p := &redis.Pool{
		MaxIdle:     rdc.MaxIdle,
		MaxActive:   rdc.MaxActive,
		IdleTimeout: 30 * time.Second,
		Wait:        true,
		Dial: func() (conn redis.Conn, err error) {
			return setDialog(rdc)
		},
	}

	rc := p.Get()
	defer rc.Close()
	_, err := rc.Do("PING")
	if err != nil {
		panic(fmt.Sprintf("[RDS] redis 初始化失败 %v", err))
		return nil
	}
	return p
}

// setDialog
func setDialog(rdc *RDSConfig) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", rdc.Host, rdc.Port))
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("连接redis错误")
	}

	if len(rdc.Password) != 0 {
		if _, err := conn.Do("AUTH", rdc.Password); err != nil {
			conn.Close()
		}
	}
	if _, err := conn.Do("SELECT", rdc.DB); err != nil {
		conn.Close()
	}

	return conn, nil
}
