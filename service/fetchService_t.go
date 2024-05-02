package service

import (
	"RankWillServer/util"
	"context"
	"log"
	"sync"
	"time"
)

func (contest *Contest) HandleContestTest(ctx context.Context) error {
	ch := util.GetChanelFromCtxByKey(ctx, util.ContestChanelKey)
	contest.rankPages = make(map[int]*RankPage)
	getContestantNumErr := util.Retry(10, 500*time.Microsecond, func() error {
		return contest.QueryContestantNumByContestName(ctx)
	})
	if getContestantNumErr != nil {
		return getContestantNumErr
	}
	contest.pageNum = (contest.contestantNum-1)/25 + 1

	contest.pageNum = 50
	log.Println("pagenum :", contest.pageNum)
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
