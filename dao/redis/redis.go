package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"web_app/settings"
)

var REDIS *redis.Client

func Init(redisCfg *settings.RedisConfig) (err error) {
	fmt.Printf("%s,%s,%i", redisCfg.Host, redisCfg.Password, redisCfg.DB)
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Host,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Error("redis connect ping failed , err :", zap.Error(err))
		return
	} else {
		zap.L().Info("redis connect ping response:", zap.String("pong", pong))
		REDIS = client
		fmt.Println("redis连接成功")
	}
	REDIS.Set(context.Background(), "name", "mai", 0)
	return
}
