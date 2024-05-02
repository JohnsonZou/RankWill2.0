package redis

import (
	"RankWillServer/backend/model"
	"context"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

func SetUserFetchRatingInfo(ctx context.Context, user *model.UserRankInfo) error {
	rdb := GetRedisClient(ctx)
	curContestSKey := BuildRedisFetchedContestantSKey(user.ContestId, user.Username)
	_, setErr := rdb.Set(curContestSKey, BuildRedisContestantSVal(user.Rating, user.AttendedContestsCount), 14*time.Hour).Result()
	return setErr
}

func GetUserFetchRatingInfo(ctx context.Context, user *model.UserRankInfo) error {
	rdb := GetRedisClient(ctx)
	SKey := BuildRedisFetchedContestantSKey(user.ContestId, user.Username)
	val, err := rdb.Get(SKey).Result()
	if err == redis.Nil {
		return err
	}
	rating, attendedContestCount, err := ParseContestSVal(val)

	user.Rating = rating
	user.AttendedContestsCount = int(attendedContestCount)
	return err
}

func SetUserPredictedRatingInfo(ctx context.Context, user *model.UserRankInfo) error {
	rdb := GetRedisClient(ctx)
	curContestSKey := BuildRedisPredictedContestantSKey(user.ContestId, user.Username)
	_, setErr := rdb.Set(curContestSKey, BuildRedisContestantSVal(user.PredictedRating, user.AttendedContestsCount), 14*time.Hour).Result()
	return setErr
}
func GetUserPredictedRatingInfo(ctx context.Context, user *model.UserRankInfo) error {
	rdb := GetRedisClient(ctx)
	SKey := BuildRedisPredictedContestantSKey(user.ContestId, user.Username)
	val, err := rdb.Get(SKey).Result()
	if err == redis.Nil {
		return err
	}
	rating, attendedContestCount, err := ParseContestSVal(val)
	user.Rating = rating
	user.AttendedContestsCount = int(attendedContestCount)
	return err
}

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
