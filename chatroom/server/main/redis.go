package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

//定义全局redis pool
var pool *redis.Pool

func initRedisPoll(address string, maxIdle, maxActive int, idelTimeOut time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     //最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
		MaxActive:   maxActive,   //最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
		IdleTimeout: idelTimeOut, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
