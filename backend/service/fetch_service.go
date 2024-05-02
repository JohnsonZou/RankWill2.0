package service

import (
	"RankWillServer/backend/model"
	"RankWillServer/backend/mq"
	"RankWillServer/backend/util"
	myredis "RankWillServer/redis"
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func FetchContest(ctx context.Context, contest *model.Contest) error {
	ch := util.GetChanelFromCtxByKey(ctx, util.ContestChanelKey)
	contest.RankPages = make(map[int]*model.RankPage)
	getContestantNumErr := util.Retry(10, 500*time.Microsecond, func() error {
		return QueryContestantNumByContestName(ctx, contest)
	})
	if getContestantNumErr != nil {
		return getContestantNumErr
	}
	contest.PageNum = (contest.ContestantNum-1)/25 + 1

	if util.IsTestMode(ctx) {
		contest.PageNum = testPageNum
	}

	for i := 1; i <= contest.PageNum; i++ {
		ch <- i
	}
	wg := sync.WaitGroup{}
	for i := 0; i < maxGoroutineNum; i++ {
		wg.Add(1)
		go func(ctx context.Context, contest *model.Contest) {
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
				page, queryRankPageErr := QueryContestRankByPage(ctx, contest, pageNum)
				if queryRankPageErr != nil {
					log.Printf("[Error][QueryPage]Contest Name: %s, Page: %d", contest.TitleSlug, pageNum)
				} else {
					log.Println("[Success]Finish page ", pageNum)
				}
				contest.Lock.Lock()
				contest.RankPages[pageNum] = page
				contest.Lock.Unlock()
			}
			wg.Done()
		}(ctx, contest)
	}
	wg.Wait()
	return HandlePages(ctx, contest)
}
func HandlePages(ctx context.Context, contest *model.Contest) error {
	totalQueue := util.Queue{}
	reqQueue := make(chan model.UserRankInfo, 100)
	for _, p := range contest.RankPages {
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
			reqQueue <- totalQueue.Pop().(model.UserRankInfo)
		}
		for i := 0; i < maxGoroutineNum; i++ {
			wg.Add(1)
			go func() {
				select {
				case u := <-reqQueue:
					err := handleUserRankInfo(ctx, &u)
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

func handleUserRankInfo(ctx context.Context, user *model.UserRankInfo) error {
	err := myredis.GetUserFetchRatingInfo(ctx, user)
	if err != redis.Nil {
		log.Println("[Fetched]", user.UserSlug)
		return nil
	}
	var uLast model.UserRankInfo
	uLast = *user
	uLast.ContestId -= 1
	err = myredis.GetUserPredictedRatingInfo(ctx, &uLast)
	if err == nil {
		user.AttendedContestsCount = uLast.AttendedContestsCount
		user.Rating = uLast.Rating
	} else if err != redis.Nil {
		return err
	} else {
		err = QueryUserCurrentRating(ctx, user)
		if err != nil {
			return err
		}
		if user.AttendedContestsCount == 0 {
			user.Rating = defaultUserRating //default
		}
	}
	return myredis.SetUserFetchRatingInfo(ctx, user)
}

func FindAndRegisterContest(ctx context.Context) error {
	rdb := myredis.GetRedisClient(ctx)
	contests, err := CNQueryUpComingContest(ctx)
	if err != nil {
		panic(err)
	}
	for _, c := range contests {
		ContestKey := myredis.BuildRedisContestKey(c.TitleSlug)
		log.Println(ContestKey)
		exist, existErr := rdb.Exists(ContestKey).Result()
		if existErr != nil {
			log.Printf("[Redis][FindAndRegisterContest] check exist fail!!!,err :%v", existErr.Error())
			return err
		}
		if exist <= 0 {
			rdb.Set(ContestKey, "key", 0)

			var gap int64 = 1080
			var delay int64 = 30

			for t := c.StartTime + delay + gap; ; t += gap {

				lastTime := t >= c.StartTime+delay+c.Duration

				var msgByte []byte
				if !lastTime {
					msgByte, _ = json.Marshal(model.MQMessage{
						ContestName:   c.TitleSlug,
						PredictNeeded: false,
					})
				} else {
					msgByte, _ = json.Marshal(model.MQMessage{
						ContestName:   c.TitleSlug,
						PredictNeeded: true,
					})
				}
				err = mq.SendMsgToDelayQueueByUnixTime(ctx, string(msgByte), t)
				if err != nil {
					log.Printf("[Redis][FindAndRegisterContest] send msg to mq fail!!!,err :%v", err.Error())
				}
				if lastTime {
					break
				}
			}
		} else {
			log.Println("[Redis][FindAndRegisterContest] contest has been send to mq before, ", ContestKey)
		}
	}
	return err
}
