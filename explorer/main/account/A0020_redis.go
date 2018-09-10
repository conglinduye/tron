package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var _redisCli *redis.Client

func getRedisClient() *redis.Client {
	return _redisCli
}

// default 127.0.0.1:6379
func initRedis(redisAddr []string) {
	redisOpt := &redis.Options{
		Addr:     redisAddr[0],
		Password: "",
		DB:       0,
	}
	_redisCli = redis.NewClient(redisOpt)

	pong, err := _redisCli.Ping().Result()
	fmt.Printf("redis ping ret:%v, err:%v\n", pong, err)
}

// redis error ...
var (
	ErrorRedisNilResult = fmt.Errorf("redis cmd result is nil")
)
