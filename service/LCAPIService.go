package service

import (
	"RankWillServer/util"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func buildCNUserRatingGraphQLQueryPostBody(userName string) string {
	return CN_UserRatingGraphQLQueryPostBodyPrefix +
		userName + CN_UserRatingGraphQLQueryPostBodySuffix
}
func buildI18NUserRatingGraphQLQueryPostBody(userName string) string {
	return I18N_UserRatingGraphQLQueryPostBodyPrefix +
		userName + I18N_UserRatingGraphQLQueryPostBodySuffix
}
func (user *userRankInfo) buildUserRatingGraphQLQueryPostRequest() (*http.Request, error) {
	if user.DataRegion == "CN" {
		return util.GenPostReq(CN_LCGraphQLURL, buildCNUserRatingGraphQLQueryPostBody(user.Username))
	} else {
		return util.GenPostReq(I18N_LCGraphQLURL, buildI18NUserRatingGraphQLQueryPostBody(user.Username))
	}
}
func (user *userRankInfo) QueryUserCurrentRating(ctx context.Context) (err error) {
	t0 := time.Now().UnixMilli()

	//!!!

	if user.DataRegion != "CN" {
		user.AttendedContestsCount = 0
		user.Rating = 1500.0
		return nil
	}
	req, genReqErr := user.buildUserRatingGraphQLQueryPostRequest()
	defer func() {
		t1 := time.Now().UnixMilli()

		if err == nil {
			s, c := util.AddCounterAndGetSpeed(ctx)
			log.Println("[QueryUserCurrentRating] success, cost ", t1-t0, " ms, current speed: ", s, ", cur total queries :", c)
		}

	}()
	if genReqErr != nil {
		//to fix
		log.Fatalf("%v\n", genReqErr.Error())
		return genReqErr
	}
	var res *http.Response
	var queryResult LCUserRatingGraphQLResult

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res, err = client.Do(req)
	defer func() {
		util.CloseResponseBody(res)
	}()

	if res.StatusCode != 200 {
		log.Printf("[!]Encounter http %d err code\n", res.StatusCode)
		log.Println(util.ReadCloserToString(res.Body))
		err = errStatusCodeNot200
		return
	}
	if err != nil {
		log.Printf("%v\n", err.Error())
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&queryResult); err != nil {
		log.Printf("%v\n", err.Error())
		return
	}
	user.AttendedContestsCount = queryResult.Data.UserContestRanking.AttendedContestsCount
	user.Rating = queryResult.Data.UserContestRanking.Rating
	return
}

func buildQueryContestRankByNameAndPageURL(contestName string, pageNum int) string {
	return CN_LCContestRankQueryPrefix + contestName + LCContestRankQueryMidfix +
		strconv.Itoa(pageNum) + LCContestRankQuerySuffix
}
func (contest *Contest) QueryContestantNumByContestName(ctx context.Context) error {
	page1, err := contest.QueryContestRankByPage(ctx, 1)
	if err != nil {
		return err
	}
	if page1.UserNum == 0 {
		return errors.New("nil page")
	}
	contest.contestantNum = page1.UserNum
	return nil
}
func (contest *Contest) QueryContestRankByPage(ctx context.Context, pageNum int) (*RankPage, error) {
	result := &RankPage{}
	url := buildQueryContestRankByNameAndPageURL(contest.TitleSlug, pageNum)

	req, genReqErr := util.GenGetReq(url)
	if genReqErr != nil {
		log.Printf("[Error][QueryContestRankByPage]%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	queryErr := util.Retry(100, 10000*time.Millisecond, func() error {
		var err error
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		res, err = client.Do(req)
		if res == nil {
			log.Printf("[Error][I18NQueryContestRankByPage]query res nil\n")
			return errResNil
		}
		if res.StatusCode != 200 {
			log.Printf("[Error][I18NQueryContestRankByPage]query res status code not 200,code:%d\n", res.StatusCode)

			//log.Printf(util.ReadCloserToString(res.Body))
			return errStatusCodeNot200
		}
		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			log.Printf("[Error][I18NQueryContestRankByPage]%v\n", err.Error())
		}
		return err
	})
	defer util.CloseResponseBody(res)
	return result, queryErr
}

func I18NQueryLatelyContest(ctx context.Context) (*Contest, error) {
	client := util.GetHttpClient(ctx)
	req, genReqErr := util.GenPostReq(I18N_LCGraphQLURL, I18N_LCContestQueryGraphQLPostBody)
	if genReqErr != nil {
		log.Printf("[Error][I18NQueryLatelyContest]%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	result := &LCContestInfoGraphQLResult{}
	queryErr := util.Retry(100, 300*time.Millisecond, func() error {
		var err error
		res, err = client.Do(req)
		if res == nil {
			log.Printf("[Error][I18NQueryLatelyContest]query res nil\n")
			return errResNil
		}
		if res.StatusCode != 200 {
			log.Printf("[Error][I18NQueryLatelyContest]query res status code not 200,code:%d\n", res.StatusCode)
			return errStatusCodeNot200
		}
		err = json.NewDecoder(res.Body).Decode(result)
		if err != nil {
			log.Printf("[Error][I18NQueryLatelyContest]%v\n", err.Error())
		}
		return err
	})
	defer util.CloseResponseBody(res)
	if result == nil || result.Data.PastContests.Data == nil {
		log.Printf("[Error][I18NQueryLatelyContest]nil contest result\n")
		return nil, errors.New("nil contest result")
	}
	return &result.Data.PastContests.Data[0], queryErr
}

func CNQueryUpComingContest(ctx context.Context) ([]*Contest, error) {
	client := util.GetHttpClient(ctx)
	req, genReqErr := util.GenPostReq(CN_LCGraphQLURL_SHORT, CN_UpComingContestGraphQLQueryPostBody)
	if genReqErr != nil {
		log.Printf("[Error][CNQueryUpComingContest]%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	result := &LCCNQueryUpComingContest{}
	queryErr := util.Retry(100, 300*time.Millisecond, func() error {
		var err error
		res, err = client.Do(req)
		if res == nil {
			log.Printf("[Error][CNQueryUpComingContest]query res nil\n")
			return errResNil
		}
		if res.StatusCode != 200 {
			log.Printf("[Error][CNQueryUpComingContest]query res status code not 200,code:%d\n", res.StatusCode)
			log.Println(util.ReadCloserToString(res.Body))
			return errStatusCodeNot200
		}
		err = json.NewDecoder(res.Body).Decode(result)
		if err != nil {
			log.Printf("[Error][CNQueryUpComingContest]%v\n", err.Error())
		}
		return err
	})
	defer util.CloseResponseBody(res)
	return result.Data.ContestUpcomingContests, queryErr
}
