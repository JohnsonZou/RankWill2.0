package service

import (
	"RankWillServer/backend/model"
	"RankWillServer/backend/util"
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
func buildUserRatingGraphQLQueryPostRequest(user *model.UserRankInfo) (*http.Request, error) {
	if user.DataRegion == "CN" {
		return util.GenPostReq(CN_LCGraphQLURL, buildCNUserRatingGraphQLQueryPostBody(user.Username))
	} else {
		return util.GenPostReq(I18N_LCGraphQLURL, buildI18NUserRatingGraphQLQueryPostBody(user.Username))
	}
}
func QueryUserCurrentRating(ctx context.Context, user *model.UserRankInfo) (err error) {
	t0 := time.Now().UnixMilli()

	req, genReqErr := buildUserRatingGraphQLQueryPostRequest(user)
	defer func() {
		t1 := time.Now().UnixMilli()

		if err == nil {
			s, c := util.AddCounterAndGetSpeed(ctx)
			log.Println("[QueryUserCurrentRating] success, cost ", t1-t0, " ms, speed: ", s, ", cur total q :", c,
				"cur user: ", user.Username, "rating: ", user.Rating, "re: ", user.DataRegion)
		}

	}()
	if genReqErr != nil {
		//to fix
		log.Fatalf("%v\n", genReqErr.Error())
		return genReqErr
	}
	var res *http.Response
	var queryResult model.LCUserRatingGraphQLResult

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
func QueryContestantNumByContestName(ctx context.Context, contest *model.Contest) error {
	page1, err := QueryContestRankByPage(ctx, contest, 1)
	if err != nil {
		return err
	}
	if page1.UserNum == 0 {
		return errors.New("nil page")
	}
	contest.ContestantNum = page1.UserNum
	return nil
}
func QueryContestRankByPage(ctx context.Context, contest *model.Contest, pageNum int) (*model.RankPage, error) {
	result := &model.RankPage{}
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

func CNQueryUpComingContest(ctx context.Context) ([]*model.Contest, error) {
	client := util.GetHttpClient(ctx)
	req, genReqErr := util.GenPostReq(CN_LCGraphQLURL_SHORT, CN_UpComingContestGraphQLQueryPostBody)
	if genReqErr != nil {
		log.Printf("[Error][CNQueryUpComingContest]%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	result := &model.LCCNQueryUpComingContest{}
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
