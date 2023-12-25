package service

import (
	"net/http"
	"sync"
	"time"
)

func (contest *Contest) HandleContest(client *http.Client, ch chan int) error {
	getContestantNumErr := Retry(10, 500*time.Microsecond, func() error {
		return contest.I18NQueryContestantNumByContestName(client)
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
		go func() {
			wg.Add(1)
			contest.handlePages(client, ch)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}
func (contest *Contest) handlePages(client *http.Client, ch chan int) error {
	for {
		pageNum, ok := <-ch
		if !ok {
			break
		}
		page, queryRankPageErr := contest.I18NQueryContestRankByPage(client, int16(pageNum))
		if queryRankPageErr != nil {
			return queryRankPageErr
		}
		contest.rankPages[pageNum] = page
		for _, u := range page.Total_Rank {
			if handleUserInfoErr := u.handleUserRankInfo(client); handleUserInfoErr != nil {
				return handleUserInfoErr
			}
		}
	}
	return nil
}
func (user userRankInfo) handleUserRankInfo(client *http.Client) error {
	//!!! to do
	_, queryRatingErr := user.QueryAndSaveUserCurrentRating(client)
	return queryRatingErr
}
