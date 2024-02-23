package redis

import (
	"douyinApi/config"

	"github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

var RedisDb *redis.Client
var Rs *redsync.Redsync

func init() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.Get().Redis.Addr,
		Password: config.Get().Redis.Password,
		DB:       int(config.Get().Redis.Db),
	})
	_, err := RedisDb.Ping().Result()
	if err != nil {
		panic(err)
	}

	// 创建redsync的客户端连接池
	pool := goredis.NewPool(RedisDb)

	// 创建redsync实例
	Rs = redsync.New(pool)

}
