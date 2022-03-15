package redis

import (
	"github.com/J-Y-Zhang/cloud-storage/file/config"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var (
	pool      *redis.Pool
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			// 1.打开连接
			conn, err := redis.Dial("tcp", config.RedisHost + ":" + config.RedisPort)
			if err != nil {
				log.Print("连接redis失败, 错误信息 ", err)
				return nil, err
			}
			// 2.访问验证
			// 我没有设置密码, 暂时不需要密码

			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Millisecond {
				return nil
			}

			_, err := conn.Do("PING")
			return err
		},
		MaxIdle:         50,
		MaxActive:       30,
		IdleTimeout:     300 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}


func RedisPool() *redis.Pool {
	if pool != nil {
		return pool
	}

	pool = newRedisPool()
	return pool
}
