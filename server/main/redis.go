package main

import (
	"time"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func initPool(address string, maxIdle int, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle : maxIdle,						//最大空闲连接数
		MaxActive : maxActive,					//和数据库的最大连接数，0表示无限制
		IdleTimeout : idleTimeout,				//最大空闲时间
		Dial : func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}