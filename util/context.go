package util

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
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
func SetTestMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, TestModeKey, "t")
}
func IsTestMode(ctx context.Context) bool {
	t := ctx.Value(TestModeKey)
	_, ok := t.(string)
	return ok
}
func SetHttpClient(ctx context.Context, cli *http.Client) context.Context {
	return context.WithValue(ctx, httpClientKey, cli)
}
func InitMainScheduler(ctx context.Context) context.Context {
	s := gocron.NewScheduler(time.Local)
	s.StartAsync()
	return context.WithValue(ctx, MainSchedulerKey, s)
}
func GetMainMQChanel(ctx context.Context) *amqp.Channel {
	ch := ctx.Value(MainMQChanelKey)
	res, ok := ch.(*amqp.Channel)
	if ok {
		return res
	}

	return nil
}
func GetMainScheduler(ctx context.Context) *gocron.Scheduler {
	cli := ctx.Value(MainSchedulerKey)
	res, ok := cli.(*gocron.Scheduler)
	if ok {
		return res
	}
	return nil
}
func InitRedisClient(ctx context.Context) (context.Context, error) {
	if ctx.Value(RedisClientKey) != nil {
		return ctx, nil
	}
	viper.SetConfigName("redisConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		return ctx, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("addr"),
		Password: viper.GetString("password"),
		DB:       viper.GetInt("db"),
	})
	if _, err := rdb.Ping().Result(); err != nil {
		log.Println(err)
		return ctx, err
	}
	return context.WithValue(ctx, RedisClientKey, rdb), nil
}
func GetRedisClient(ctx context.Context) *redis.Client {
	cli := ctx.Value(RedisClientKey)
	res, ok := cli.(*redis.Client)
	if ok {
		return res
	}
	return nil
}

func InitTimer(ctx context.Context) context.Context {
	cntMap := make(map[string]int)
	return context.WithValue(context.WithValue(ctx, StartTimeKey, time.Now().UnixMilli()), QueryCounterKey, cntMap)
}

func AddCounterAndGetSpeed(ctx context.Context) (float64, int) {
	deltaT := time.Now().UnixMilli() - ctx.Value(StartTimeKey).(int64)
	cnt := ctx.Value(QueryCounterKey).(map[string]int)
	cnt[QueryCounterKey]++
	return 1000 * float64(cnt[QueryCounterKey]) / float64(deltaT), cnt[QueryCounterKey]
}
