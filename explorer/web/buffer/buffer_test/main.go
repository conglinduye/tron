package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/wlcy/tron/explorer/lib/mysql"
)

func main() {
	mysql.Initialize("mine", "3306", "tron", "tron", "tron")

	initRedis([]string{"127.0.0.1:6379"})

	bb := getBlockBuffer()

	go bb.backgroundWorker()

	cnt := 0
	for cnt < 10 {

		fmt.Printf("nowblock:%v, confirmed blockID:%v\n\n\n\n", bb.maxBlockID, bb.maxConfirmedBlockID)
		time.Sleep(5 * time.Second)
		cnt++

		tsr := time.Now()
		rs := int64(0)
		re := int64(50)
		ret, _ := bb.GetBlocks(-1, rs, re)
		retLen := len(ret)
		fmt.Printf("\nload from buffer %v~ %v (%v), size:%v, ret[0].num:%v, ret[%v].num:%v, cost:%v\n\n", rs, re, re, len(ret), ret[0].Number, retLen, ret[retLen-1].Number, time.Since(tsr))
		var c, unc int
		var minCBlockID int64 = 900000000
		var maxCBlockID int64
		var maxUncBlockID int64
		var minUncBlockID int64 = 9000000000
		for _, block := range ret {
			if block.Confirmed {
				c++
				if maxCBlockID < block.Number {
					maxCBlockID = block.Number
				}
				if minCBlockID > block.Number {
					minCBlockID = block.Number
				}
			} else {
				unc++
				if maxUncBlockID < block.Number {
					maxUncBlockID = block.Number
				}
				if minUncBlockID > block.Number {
					minUncBlockID = block.Number
				}
			}
		}
		fmt.Printf("(min, max) confirmed block id:(%v,%v) count:%v;  (min, max) unconfirmed block id:(%v,%v) count:%v\n", minCBlockID, maxCBlockID, maxCBlockID-minCBlockID+1, minUncBlockID, maxUncBlockID, maxUncBlockID-minUncBlockID+1)

		tsr = time.Now()
		rs = int64(50)
		re = int64(100)
		ret, _ = bb.GetBlocks(-1, 50, 100)
		retLen = len(ret)
		fmt.Printf("\nload from buffer %v~ %v (%v), size:%v, ret[0].num:%v, ret[%v].num:%v, cost:%v\n\n", rs, re, re, len(ret), ret[0].Number, retLen, ret[retLen-1].Number, time.Since(tsr))

		minCBlockID, maxCBlockID, minUncBlockID, maxUncBlockID = 9000000000, 0, 9000000000, 0
		for _, block := range ret {
			if block.Confirmed {
				c++
				if maxCBlockID < block.Number {
					maxCBlockID = block.Number
				}
				if minCBlockID > block.Number {
					minCBlockID = block.Number
				}
			} else {
				unc++
				if maxUncBlockID < block.Number {
					maxUncBlockID = block.Number
				}
				if minUncBlockID > block.Number {
					minUncBlockID = block.Number
				}
			}
		}
		fmt.Printf("(min, max) confirmed block id:(%v,%v) count:%v;  (min, max) unconfirmed block id:(%v,%v) count:%v\n", minCBlockID, maxCBlockID, maxCBlockID-minCBlockID+1, minUncBlockID, maxUncBlockID, maxUncBlockID-maxUncBlockID+1)

	}
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
