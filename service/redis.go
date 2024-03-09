package service

import (
	"RankWillServer/util"
	"context"
	"github.com/go-redis/redis"
	"time"
)

func (user *userRankInfo) setUserFetchRatingInfo(ctx context.Context) error {
	rdb := util.GetRedisClient(ctx)
	curContestSKey := util.BuildRedisFetchedContestantSKey(user.ContestId, user.Username)
	_, setErr := rdb.Set(curContestSKey, util.BuildRedisContestantSVal(user.Rating, user.AttendedContestsCount), 14*time.Hour).Result()
	return setErr
}

func (user *userRankInfo) getUserFetchRatingInfo(ctx context.Context) error {
	rdb := util.GetRedisClient(ctx)
	SKey := util.BuildRedisFetchedContestantSKey(user.ContestId, user.Username)

	val, err := rdb.Get(SKey).Result()
	if err == redis.Nil {
		return err
	}
	rating, attendedContestCount, err := util.ParseContestSVal(val)

	user.Rating = rating
	user.AttendedContestsCount = int(attendedContestCount)
	return err
}

func (user *userRankInfo) setUserPredictedRatingInfo(ctx context.Context) error {
	rdb := util.GetRedisClient(ctx)
	curContestSKey := util.BuildRedisPredictedContestantSKey(user.ContestId, user.Username)
	_, setErr := rdb.SetNX(curContestSKey, util.BuildRedisContestantSVal(user.Rating, user.AttendedContestsCount), 14*time.Hour).Result()
	return setErr
}
func (user *userRankInfo) getUserPredictedRatingInfo(ctx context.Context) error {
	rdb := util.GetRedisClient(ctx)
	SKey := util.BuildRedisPredictedContestantSKey(user.ContestId, user.Username)
	val, err := rdb.Get(SKey).Result()
	if err == redis.Nil {
		return err
	}
	rating, attendedContestCount, err := util.ParseContestSVal(val)
	user.Rating = rating
	user.AttendedContestsCount = int(attendedContestCount)
	return err
}
