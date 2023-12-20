package service

import (
	"RankWillServer/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func buildI18NUserRatingGraphQLQueryPostBody(userName string) string {
	return I18N_UserRatingGraphQLQueryPostBodyPrefix +
		userName + I18N_UserRatingGraphQLQueryPostBodySuffix
}
func I18NQueryUserCurrentRating(client *http.Client, userName string) (userInfo *model.LCUserInfo, err error) {
	bodyStr := buildI18NUserRatingGraphQLQueryPostBody(userName)
	req, genReqErr := GenPostReq(I18N_LCGraphQLURL, bodyStr)

	if err != nil {
		log.Fatalf("%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	res, queryErr := client.Do(req)
	defer res.Body.Close()
	if queryErr != nil {
		log.Fatalf("%v\n", queryErr.Error())
		return nil, err
	}
	var queryResult LCUserRatingGraphQLResult
	if resDecodeErr := json.NewDecoder(res.Body).Decode(&queryResult); resDecodeErr != nil {
		log.Fatalf("%v\n", resDecodeErr.Error())
		return nil, resDecodeErr
	}
	userInfo.UserName = userName
	userInfo.Rating = queryResult.Data.UserContestRanking.Rating
	userInfo.AttendedContestsCount = queryResult.Data.UserContestRanking.AttendedContestsCount
	return userInfo, nil
}

func buildQueryContestRankByNameAndPageURL(contestName string, pageNum int16) string {
	return "https://leetcode.com/contest/api/ranking/" + contestName + "/?pagination=" +
		strconv.Itoa(int(pageNum)) + "&region=global"
}
func QueryContestRankByNameAndPage(client *http.Client, contestName string, pageNum int16) (*RankPage, error) {
	result := &RankPage{}
	url := buildQueryContestRankByNameAndPageURL(contestName, pageNum)
	req, genReqErr := GenGetReq(url)
	if genReqErr != nil {
		log.Fatalf("%v\n", genReqErr.Error())
		return nil, genReqErr
	}
	var res *http.Response
	queryErr := Retry(100, 300*time.Millisecond, func() error {
		var err error
		res, err = client.Do(req)
		return err
	})
	if queryErr != nil {
		log.Fatalf("%v\n", queryErr.Error())
		return nil, queryErr
	}
	defer res.Body.Close()
	if resDecodeErr := json.NewDecoder(res.Body).Decode(&result); resDecodeErr != nil {
		log.Fatalf("%v\n", resDecodeErr.Error())
		return nil, resDecodeErr
	}
	return result, nil
}
