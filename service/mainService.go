package service

import (
	"net/http"
	"strconv"
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

		}
	}
	return nil
}
func (user userRankInfo) handleUserRankInfo(client *http.Client) error {
	if user.DataRegion == "US" {
		uinfo, queryRatingErr := I18NQueryUserCurrentRating(client, user.Username)
		if queryRatingErr != nil {
			return queryRatingErr
		}
		//redis cli get
		redisStrKey := strconv.Itoa(int(user.ContestId)) + "#" + user.Username
		redisStrVal := strconv.FormatFloat(uinfo.Rating, 'f', 3, 64)

	}
	if user.DataRegion == "CN" {
		uinfo, err := user.QueryUserCurrentRating(client)
	}
	return nil
}
