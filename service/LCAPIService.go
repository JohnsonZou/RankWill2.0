package service

import (
	"RankWillServer/model"
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
		return GenPostReq(CN_LCGraphQLURL, buildCNUserRatingGraphQLQueryPostBody(user.Username))
	} else {
		return GenPostReq(I18N_LCGraphQLURL, buildI18NUserRatingGraphQLQueryPostBody(user.Username))
	}
}
func (user *userRankInfo) QueryUserCurrentRating(client *http.Client) (userInfo *model.LCUserInfo, err error) {
	bodyStr := buildI18NUserRatingGraphQLQueryPostBody(user.Username)
	req, genReqErr := GenPostReq(I18N_LCGraphQLURL, bodyStr)

	if genReqErr != nil {
		log.Fatalf("%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	var queryResult model.LCUserRatingGraphQLResult
	if retryErr := Retry(100, 500*time.Millisecond, func() error {
		var err error
		res, err = client.Do(req)
		defer closeResponseBody(res)
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
	userInfo.UserName = user.Username
	userInfo.Rating = queryResult.Data.UserContestRanking.Rating
	userInfo.AttendedContestsCount = queryResult.Data.UserContestRanking.AttendedContestsCount
	return userInfo, nil
}

func buildQueryContestRankByNameAndPageURL(contestName string, pageNum int16) string {
	return "https://leetcode.com/contest/api/ranking/" + contestName + "/?pagination=" +
		strconv.Itoa(int(pageNum)) + "&region=global"
}
func (contest *Contest) I18NQueryContestantNumByContestName(client *http.Client) error {
	page1, err := contest.I18NQueryContestRankByPage(client, 1)
	if err != nil {
		return err
	}
	if page1.UserNum == 0 {
		return errors.New("nil page")
	}
	contest.contestantNum = page1.UserNum
	return nil
}
func (contest *Contest) I18NQueryContestRankByPage(client *http.Client, pageNum int16) (*RankPage, error) {
	result := &RankPage{}
	url := buildQueryContestRankByNameAndPageURL(contest.contestName, pageNum)
	req, genReqErr := GenGetReq(url)
	if genReqErr != nil {
		log.Fatalf("%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	queryErr := Retry(100, 300*time.Millisecond, func() error {
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
