package util

import (
	"context"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func GetHttpClient(ctx context.Context) *http.Client {
	cli := ctx.Value(httpClientKey)
	res, ok := cli.(*http.Client)
	if ok {
		return res
	}
	return nil
}

func GetChanelFromCtxByKey(ctx context.Context, key string) chan int {
	ch := ctx.Value(key)
	res, ok := ch.(chan int)
	if ok {
		return res
	}
	return nil
}

func SetHttpClient(ctx context.Context, cli *http.Client) context.Context {
	return context.WithValue(ctx, httpClientKey, cli)
}

func InitRedisClient(ctx context.Context) error {
	if ctx.Value(redisClientKey) != nil {
		return nil
	}
	viper.SetConfigFile("../redisConfig.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("addr"),
		Password: viper.GetString("password"),
		DB:       viper.GetInt("db"),
	})
	if _, err := rdb.Ping().Result(); err != nil {
		return err
	}
	ctx = context.WithValue(ctx, redisClientKey, rdb)
	return nil
}
func GetRedisClient(ctx context.Context) *redis.Client {
	cli := ctx.Value(redisClientKey)
	res, ok := cli.(*redis.Client)
	if ok {
		return res
	}
	return nil
}
