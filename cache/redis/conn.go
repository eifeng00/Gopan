package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool       *redis.Pool
	redisHost  = "127.0.0.1:6379"
	redispPass = "123456"
)

// newRedisPool : 创建redis连接池
func newRedisPool() *redis.Pool {
	// docker run -p 6379:6379 -v /data/redis/:/data --name fileserver_redis -d redis --requirepass "123456" --appendonly yes
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			//1. open connect
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Printf(err.Error())
				return nil, err
			}
			//2. auth
			if _, err := c.Do("AUTH", redispPass); err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}
