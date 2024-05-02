package util

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func InitRedisClient(ctx context.Context) (context.Context, error) {
	if ctx.Value(RedisClientKey) != nil {
		return ctx, nil
	}
	viper.SetConfigName("redis_config")
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

func BuildRedisContestKey(contestName string) string {
	return "ContestKey###" + contestName
}
func BuildRedisFetchedContestantSKey(contestID int, uname string) string {
	return strconv.Itoa(contestID) + "###" + uname
}
func BuildRedisContestantSVal(rating float64, attendedContestCount int) string {
	return strconv.FormatFloat(rating, 'f', -1, 64) + "#" + strconv.Itoa(attendedContestCount)
}

func BuildRedisPredictedContestantSKey(contestID int, uname string) string {
	return strconv.Itoa(contestID) + "######" + uname
}
func ParseContestSVal(key string) (rating float64, attendedContestCount int64, err error) {
	strArr := strings.Split(key, "#")
	rating, err = strconv.ParseFloat(strArr[0], 64)
	if err != nil {
		return
	}
	attendedContestCount, err = strconv.ParseInt(strArr[1], 10, 32)
	return
}
