package service

import (
	"RankWillServer/util"
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func (contest *Contest) HandleContest(ctx context.Context) error {
	ch := util.GetChanelFromCtxByKey(ctx, util.ContestChanelKey)
	contest.rankPages = make(map[int]*RankPage)
	getContestantNumErr := util.Retry(10, 500*time.Microsecond, func() error {
		return contest.QueryContestantNumByContestName(ctx)
	})
	if getContestantNumErr != nil {
		return getContestantNumErr
	}
	contest.pageNum = (contest.contestantNum-1)/25 + 1

	for i := 1; i <= contest.pageNum; i++ {
		ch <- i
	}
	wg := sync.WaitGroup{}
	for i := 0; i < maxGoroutineNum; i++ {
		wg.Add(1)
		go func(ctx context.Context, contest *Contest) {
			chGet := util.GetChanelFromCtxByKey(ctx, util.ContestChanelKey)
			for {
				ok := true
				var pageNum int
				select {
				case pageNum = <-chGet:
				default:
					ok = false
				}
				if !ok {
					break
				}
				page, queryRankPageErr := contest.QueryContestRankByPage(ctx, pageNum)
				if queryRankPageErr != nil {
					log.Printf("[Error][QueryPage]Contest Name: %s, Page: %d", contest.TitleSlug, pageNum)
				} else {
					log.Println("[Success]Finish page ", pageNum)
				}

				contest.Lock.Lock()
				contest.rankPages[pageNum] = page
				contest.Lock.Unlock()
			}
			wg.Done()
		}(ctx, contest)
	}
	wg.Wait()
	return contest.HandlePages(ctx)
}
func (contest *Contest) HandlePages(ctx context.Context) error {
	totalQueue := util.Queue{}
	reqQueue := make(chan userRankInfo, 100)
	for _, p := range contest.rankPages {
		for _, u := range p.TotalRank {
			//bugfix
			totalQueue.Push(u)
		}
	}
	waitTime := 0
	for !totalQueue.Empty() {
		if waitTime > 0 {
			log.Println("[Sleep] for ", waitTime, " seconds.")
		}
		time.Sleep(time.Duration(waitTime) * time.Second)
		waitTime = 0
		wg := sync.WaitGroup{}
		for i := 0; i < maxGoroutineNum; i++ {
			if totalQueue.Empty() {
				break
			}
			reqQueue <- totalQueue.Pop().(userRankInfo)
		}
		for i := 0; i < maxGoroutineNum; i++ {
			wg.Add(1)
			go func() {
				select {
				case u := <-reqQueue:
					err := u.handleUserRankInfo(ctx)
					if err != nil {
						totalQueue.Push(u)
						waitTime++
					}
				default:
					break
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
	return nil
}

func (user *userRankInfo) handleUserRankInfo(ctx context.Context) error {
	err := user.getUserFetchRatingInfo(ctx)
	if err != redis.Nil {
		log.Println("[Fetched]", user.UserSlug)
		return nil
	}
	var uLast userRankInfo
	uLast = *user
	uLast.ContestId -= 1
	err = uLast.getUserPredictedRatingInfo(ctx)
	if err == nil {
		user.AttendedContestsCount = uLast.AttendedContestsCount
		user.Rating = uLast.Rating
	} else if err != redis.Nil {
		return err
	} else {
		err = user.QueryUserCurrentRating(ctx)
		if err != nil {
			return err
		}
		if user.AttendedContestsCount == 0 {
			user.Rating = defaultUserRating //default
		}
	}
	return user.setUserFetchRatingInfo(ctx)
}

func FindAndRegisterContest(ctx context.Context) error {
	rdb := util.GetRedisClient(ctx)
	contests, err := CNQueryUpComingContest(ctx)
	if err != nil {
		panic(err)
	}
	for _, c := range contests {

		ContestKey := util.BuildRedisContestKey(c.TitleSlug)
		exist, existErr := rdb.Exists(ContestKey).Result()
		if existErr != nil {
			log.Printf("[Redis][FindAndRegisterContest] check exist fail!!!,err :%v", existErr.Error())
			return err
		}
		if exist <= 0 {
			rdb.Set(ContestKey, "key", 0)
			err = util.SendMsgToDelayQueueByUnixTime(ctx, c.TitleSlug, c.StartTime)
			if err != nil {
				log.Printf("[Redis][FindAndRegisterContest] send msg to mq fail!!!,err :%v", err.Error())
			}

		}
	}
	return err
}
