package module

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
	"github.com/wlcy/tron/explorer/lib/mysql"
)

func TestB(*testing.T) {
	mysql.Initialize("mine", "3306", "tron", "tron", "tron")

	initRedis([]string{"127.0.0.1:6379"})

	bb := getBlockBuffer()

	bb.getNowConfirmedBlock()
	for {
		if bb.getNowBlock() {
			break
		}
	}
	fmt.Printf("nowblock:%v, confirmed blockID:%v\n", bb.maxBlockID, bb.maxConfirmedBlockID)
}

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
