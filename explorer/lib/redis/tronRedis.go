// tronRedis
// tron 封装 redis client
// 基于 gopkg.in/redis.v4, 直接使用其命令接口以及返回类型
package redis

import (
	"fmt"

	src "gopkg.in/redis.v4"
)

// TronRedis 封装对象
// 集成 gopkg.in/redis.v4 redis.Client 对象
// 开放 redis 连接信息
type TronRedis struct {
	*src.Client        // 集成 redis.Client
	Addr        string // redis 连接地址
	Password    string // redis 密码
	DB          int    // 连接的 redis DB ID
	PoolSize    int    // 连接池大小 不填默认为10
}

// NewClient 创建一个连接指定 Redis 的客户端接口
// 接口可以在多个 goruntine 间共享
func NewClient(addr, password string, db, poolSize int) *TronRedis {

	redisOptions := &src.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	}

	client := src.NewClient(redisOptions)
	ret := &TronRedis{
		client,
		addr,
		password,
		db,
		poolSize,
	}

	return ret
}

//返回redis的信息
func (r *TronRedis) String() string {
	return fmt.Sprintf("redis info : host[%v] pass[%v] DB[%v] ", r.Addr, r.Password, r.DB)
}
