package rawredis

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

// DSN redis connection string
var DSN string

var _redis *redis.Client
var _once sync.Once
var _dsn string

// GetRedisClient 获取redis连接单例
func GetRedisClient() *redis.Client {
	_once.Do(func() {
		_dsn = DSN
		var err error
		_redis, err = InitRedis(_dsn)
		if nil != err || nil == _redis {
			panic(fmt.Sprintf("connect to redis server:[%v] failed:%v", _dsn, err))
		}
	})
	return _redis
}

// InitRedis 创建连接指定redis-server的redis-client
func InitRedis(redisAddr string) (*redis.Client, error) {
	redisOpt := &redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	}
	ret := redis.NewClient(redisOpt)

	pong, err := ret.Ping().Result()
	if nil != err || 0 == len(pong) {
		return nil, err
	}
	return ret, err
}
