package service

import (
	"RankWillServer/util"
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func (contest *Contest) HandleContest(ctx context.Context) error {
	ch := util.GetChanelFromCtxByKey(ctx, util.ContestChanelKey)
	getContestantNumErr := util.Retry(10, 500*time.Microsecond, func() error {
		return contest.I18NQueryContestantNumByContestName(ctx)
	})
	if getContestantNumErr != nil {
		return getContestantNumErr
	}
	contest.pageNum = (contest.contestantNum-1)/25 + 1
	for i := 1; i <= int(contest.pageNum); i++ {
		ch <- i
	}
	wg := sync.WaitGroup{}
	for i := 0; i < maxGoroutineNum; i++ {
		go func(ctx context.Context, contest *Contest) {
			wg.Add(1)
			ch := util.GetChanelFromCtxByKey(ctx, util.ContestChanelKey)
			for {
				pageNum, ok := <-ch
				if !ok {
					break
				}
				page, queryRankPageErr := contest.I18NQueryContestRankByPage(ctx, pageNum)
				if queryRankPageErr != nil {
					//to fix !!! fatal
					log.Fatalf("[Error][QueryPage]Contest Name: %s, Page: %d", contest.contestName, pageNum)
				}
				contest.rankPages[pageNum] = page
				for _, u := range page.Total_Rank {
					if handleUserInfoErr := u.handleUserRankInfo(ctx); handleUserInfoErr != nil {
						log.Fatalf("[Error][QueryUser]Contest Name: %s, Page: %d, Username: %s", contest.contestName, pageNum, u.Username)
					}
				}
			}
			wg.Done()
		}(ctx, contest)
	}
	wg.Wait()
	return nil
}
func buildRedisContestantSKey(contestID int, uname string) string {
	return strconv.Itoa(contestID) + "###" + uname
}
func buildRedisContestantSVal(rating float64, attendedContestCount int) string {
	return strconv.FormatFloat(rating, 'f', -1, 64) + "#" + strconv.Itoa(attendedContestCount)
}
func parseContestSVal(key string) (rating float64, attendedContestCount int64, err error) {
	strArr := strings.Split(key, "#")
	rating, err = strconv.ParseFloat(strArr[0], 64)
	if err != nil {
		return
	}
	attendedContestCount, err = strconv.ParseInt(strArr[1], 10, 32)
	return
}
func (user *userRankInfo) handleUserRankInfo(ctx context.Context) error {
	//!!! to do
	rdb := util.GetRedisClient(ctx)
	curContestSKey := buildRedisContestantSKey(user.ContestId, user.Username)
	exist, err := rdb.Exists(curContestSKey).Result()
	if err != nil {
		return err
	}
	if exist > 0 {
		return nil
	}
	lastContestSKey := buildRedisContestantSKey(user.ContestId-1, user.Username)
	lastContestVal, getErr := rdb.Get(lastContestSKey).Result()
	var curContestRating float64
	var curContestAC int
	if getErr == nil {
		rating, attendContestCount, parseErr := parseContestSVal(lastContestVal)
		if parseErr != nil {
			log.Printf("[Redis]Contest Val Parse fail. err: %v", parseErr)
		}
		curContestAC = int(attendContestCount) + 1
		curContestRating = rating
	} else if getErr != redis.Nil {
		return getErr
	} else {
		info, queryRatingErr := user.QueryUserCurrentRating(ctx)
		if queryRatingErr != nil {
			return queryRatingErr
		}
		curContestRating = info.Rating
		curContestAC = info.AttendedContestsCount
	}

	_, setErr := rdb.SetNX(curContestSKey, buildRedisContestantSVal(curContestRating, curContestAC), 14*time.Hour).Result()
	return setErr
}
