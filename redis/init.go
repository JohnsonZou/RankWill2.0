package redis

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

const KeyRedisClient = "redis_client"

func InitRedisClient(ctx context.Context) (context.Context, error) {
	if ctx.Value(KeyRedisClient) != nil {
		return ctx, nil
	}
	vp := viper.New()
	vp.SetConfigName("redis_config")
	vp.SetConfigType("yaml")
	dir, _ := os.Getwd()
	vp.AddConfigPath(dir + "\\config\\")

	if err := vp.ReadInConfig(); err != nil {
		return ctx, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     vp.GetString("addr"),
		Password: vp.GetString("password"),
		DB:       vp.GetInt("db"),
	})
	if _, err := rdb.Ping().Result(); err != nil {
		log.Println(err)
		return ctx, err
	}
	return context.WithValue(ctx, KeyRedisClient, rdb), nil
}
func GetRedisClient(ctx context.Context) *redis.Client {
	cli := ctx.Value(KeyRedisClient)
	res, ok := cli.(*redis.Client)
	if ok {
		return res
	}
	return nil
}
