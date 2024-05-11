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
	CNQueue := util.Queue{}
	USQueue := util.Queue{}
	for _, p := range contest.RankPages {
		for i := range p.TotalRank {
			p.TotalRank[i].ContestName = contest.TitleSlug
			if p.TotalRank[i].DataRegion == "CN" {
				CNQueue.Push(&p.TotalRank[i])
			} else {
				USQueue.Push(&p.TotalRank[i])
			}
		}
	}
	//during the multi-goroutine request
	//delay mainly due to the difference of data region
	//the time sonsumption of a round of multi-goroutine request
	//mainly depends on the slowest request in a set of requests
	//then, reqs for different data region should be divided
	HandleTotalQueue(ctx, &CNQueue)
	HandleTotalQueue(ctx, &USQueue)
	return nil
}

func HandleTotalQueue(ctx context.Context, queue *util.Queue) {
	reqQueue := make(chan *model.UserRankInfo, 100)
	waitTime := 0
	dynamicRoutineNum := 1
	for !queue.Empty() {
		if waitTime > 0 {
			log.Println("[Sleep] for ", waitTime, " seconds.")
			time.Sleep(time.Duration(waitTime) * time.Second)
			dynamicRoutineNum = 1
		} else {
			//dynamicRoutineNum may exceed maxGoroutineNum
			//the 'max' is just a kinda limitation
			if dynamicRoutineNum < maxGoroutineNum {
				dynamicRoutineNum *= 2
			} else if dynamicRoutineNum < 2*maxGoroutineNum {
				dynamicRoutineNum++
			}
		}
		waitTime = 0
		wg := sync.WaitGroup{}
		for i := 0; i < dynamicRoutineNum; i++ {
			if queue.Empty() {
				break
			}
			reqQueue <- queue.Pop().(*model.UserRankInfo)
		}
		for i := 0; i < dynamicRoutineNum; i++ {
			wg.Add(1)
			go func() {
				select {
				case u := <-reqQueue:
					err := handleUserRankInfo(ctx, u)
					if err != nil {
						log.Printf("[Error] Handle user rank info, err: %+v\n", err)
						queue.Push(u)
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
}

func handleUserRankInfo(ctx context.Context, user *model.UserRankInfo) error {
	//fetched
	err := myredis.GetUserFetchRatingInfo(ctx, user)
	if err != redis.Nil {
		log.Println("[Fetched]", user.UserSlug)
		return nil
	}

	//for the weekly-contest after biweekly-contest
	//rating data should inherit the predicted result of biweekly-contest
	var uLast model.UserRankInfo = *user
	uLast.ContestId -= 1
	err = myredis.GetUserPredictedRatingInfo(ctx, &uLast)
	if err == nil {
		//inherit the predicted result of biweekly-contest
		user.AttendedContestsCount = uLast.AttendedContestsCount
		user.Rating = uLast.Rating
	} else if err != redis.Nil {
		//unexpected err
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
			AddContestIntoMQ(ctx, c)
		} else {
			log.Println("[Redis][FindAndRegisterContest] contest has been send to mq before, ", ContestKey)
		}
	}
	return err
}

func AddContestIntoMQ(ctx context.Context, c *model.Contest) {

	//gap is the actual time gap between each total query for the dashboard
	var gap int64 = 1080
	var delay int64 = 30

	for t := c.StartTime + delay + gap; ; t += gap {

		//the last query should end with predict service
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
		err := mq.SendMsgToDelayQueueByUnixTime(ctx, string(msgByte), t)
		if err != nil {
			log.Printf("[Redis][FindAndRegisterContest] send msg to mq fail!!!,err :%v", err.Error())
		}

		//important!!!
		if lastTime {
			break
		}
	}
}
