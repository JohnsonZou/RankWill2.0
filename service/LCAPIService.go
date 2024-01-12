package service

import (
	"RankWillServer/util"
	"context"
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
func (user userRankInfo) buildUserRatingGraphQLQueryPostRequest() (*http.Request, error) {
	if user.DataRegion == "CN" {
		return util.GenPostReq(CN_LCGraphQLURL, buildCNUserRatingGraphQLQueryPostBody(user.Username))
	} else {
		return util.GenPostReq(I18N_LCGraphQLURL, buildI18NUserRatingGraphQLQueryPostBody(user.Username))
	}
}
func (user *userRankInfo) QueryUserCurrentRating(ctx context.Context) (*LCUserInfo, error) {
	req, genReqErr := user.buildUserRatingGraphQLQueryPostRequest()
	client := util.GetHttpClient(ctx)
	if genReqErr != nil {
		//to fix
		log.Fatalf("%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	var queryResult LCUserRatingGraphQLResult
	if retryErr := util.Retry(100, 500*time.Millisecond, func() error {
		var err error
		res, err = client.Do(req)
		defer util.CloseResponseBody(res)
		if err != nil {
			log.Printf("%v\n", err.Error())
			return err
		}
		if res.StatusCode != 200 {
			return errStatusCodeNot200
		}
		if err = json.NewDecoder(res.Body).Decode(&queryResult); err != nil {
			log.Printf("%v\n", err.Error())
		}
		return err
	}); retryErr != nil {
		log.Fatalf("%v\n", retryErr)
	}
	return &LCUserInfo{
		UserName:              user.Username,
		UserSlug:              user.UserSlug,
		Rating:                queryResult.Data.UserContestRanking.Rating,
		AttendedContestsCount: queryResult.Data.UserContestRanking.AttendedContestsCount,
	}, nil
}

func buildQueryContestRankByNameAndPageURL(contestName string, pageNum int) string {
	return I18N_LCContestRankQueryPrefix + contestName + I18N_LCContestRankQueryMidfix +
		strconv.Itoa(pageNum) + I18N_LCContestRankQuerySuffix
}
func (contest *Contest) I18NQueryContestantNumByContestName(ctx context.Context) error {
	page1, err := contest.I18NQueryContestRankByPage(ctx, 1)
	if err != nil {
		return err
	}
	if page1.UserNum == 0 {
		return errors.New("nil page")
	}
	contest.contestantNum = page1.UserNum
	return nil
}
func (contest *Contest) I18NQueryContestRankByPage(ctx context.Context, pageNum int) (*RankPage, error) {
	client := util.GetHttpClient(ctx)
	result := &RankPage{}
	url := buildQueryContestRankByNameAndPageURL(contest.contestName, pageNum)
	req, genReqErr := util.GenGetReq(url)
	if genReqErr != nil {
		log.Fatalf("%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	queryErr := util.Retry(100, 300*time.Millisecond, func() error {
		var err error
		res, err = client.Do(req)
		if res.StatusCode != 200 {
			return errStatusCodeNot200
		}
		return err
	})
	if queryErr != nil {
		log.Fatalf("%v\n", queryErr.Error())
		return nil, queryErr
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("%v\n", err.Error())
		}
	}()
	if resDecodeErr := json.NewDecoder(res.Body).Decode(&result); resDecodeErr != nil {
		log.Fatalf("%v\n", resDecodeErr.Error())
		return nil, resDecodeErr
	}
	return result, nil
}
