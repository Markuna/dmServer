package api

import (
	"douyinApi/config"
	"douyinApi/redis"
	"log"
	"time"
)

const access_token_key = "douyin:access:token"
const access_token_key_lock = access_token_key + ":lock"
const expire_time = 60 * time.Minute
const loop_time = 10 * time.Minute

func init() {
	go func() {
		for {
			// redis 取锁
			mutex := redis.Rs.NewMutex(access_token_key_lock)
			// 分布式加锁
			if err := mutex.Lock(); err != nil {
				log.Println(err)
				// 睡眠10分钟
				time.Sleep(loop_time)
				continue
			}
			// 拿锁成功的话, 检查过期时间是否小于轮询时间
			if getAccessTokenTTL() < loop_time+time.Minute {
				// 获取新refresh_token 更新到redis
				updateRedisAccessToken()
				time.Sleep(2 * time.Second)
			}
			// 释放互斥锁
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
			// 睡眠10分钟
			time.Sleep(loop_time)
		}
	}()
}

func GetAccessToken() string {
	value := redis.RedisDb.Get(access_token_key).Val()
	if len(value) == 0 || value == "" {
		return updateRedisAccessToken()
	}
	return value
}

func getAccessTokenTTL() time.Duration {
	ttl := redis.RedisDb.TTL(access_token_key).Val()
	return ttl
}

func updateRedisAccessToken() string {
	req := &AccessTokenReq{
		AppId:      config.Get().Douyin.AppId,
		Secret:     config.Get().Douyin.AppSecret,
		Grant_type: "client_credential",
	}
	resp := AccessTokenApi(req)
	if resp != nil {
		newToken := resp.AccessToken
		redis.RedisDb.Set(access_token_key, newToken, expire_time).Err()
		return newToken
	}
	return ""
}
