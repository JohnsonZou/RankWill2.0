package redis

import (
	"RankWillServer/backend/model"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func SetUserFetchRatingInfo(ctx context.Context, user *model.UserRankInfo) error {
	rdb := GetRedisClient(ctx)
	curContestSKey := BuildRedisFetchedContestantSKey(user.ContestId, user.Username)
	_, setErr := rdb.Set(curContestSKey, BuildRedisContestantSVal(user.Rating, user.AttendedContestsCount), 14*time.Hour).Result()
	if setErr != nil {
		log.Printf("[Error][SetUserFetchRatingInfo]err: %+v\n", setErr)
		return setErr
	}
	return nil
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

	userByte, _ := json.Marshal(user)
	_, ZAddErr := rdb.ZAdd(user.ContestName, redis.Z{Score: float64(user.Rank), Member: userByte}).Result()
	if ZAddErr != nil {
		log.Printf("[Error][SetUserPredictedRatingInfo]err: %+v\n", ZAddErr)
		return ZAddErr
	}
	_, setErr := rdb.Set(curContestSKey, BuildRedisContestantSVal(user.PredictedRating, user.AttendedContestsCount), 14*time.Hour).Result()
	if setErr != nil {
		log.Printf("[Error][SetUserPredictedRatingInfo]err: %+v\n", setErr)
		return setErr
	}
	return nil
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
func GetUserByRank(ctx context.Context, contestName string, l, r int) (res []model.UserRankInfo, err error) {
	rdb := GetRedisClient(ctx)

	members, err := rdb.ZRangeByScore(contestName, redis.ZRangeBy{
		Min: strconv.FormatInt(int64(l), 10),
		Max: strconv.FormatInt(int64(r), 10)}).Result()

	if err != nil {
		log.Printf("[Error][GetUserByRank]err: %+v\n", err)
		return nil, err
	}
	for _, v := range members {
		var u model.UserRankInfo
		err = json.Unmarshal([]byte(v), &u)
		if err != nil {
			log.Printf("[Error][GetUserByRank] unmarshal err: %+v\n", err)
			return nil, err
		}
		res = append(res, u)
	}
	return
}
