package service

import (
	"RankWillServer/model"
	"encoding/json"
	"log"
	"net/http"
)

func buildCNUserRatingGraphQLQueryPostBody(userName string) string {
	return CN_UserRatingGraphQLQueryPostBodyPrefix +
		userName + CN_UserRatingGraphQLQueryPostBodySuffix
}
func CNQueryusercurrentrating(client *http.Client, userName string) (userInfo model.LCUserInfo) {
	bodyStr := buildCNUserRatingGraphQLQueryPostBody(userName)
	req, err := GenPostReq(CN_LCGraphQLURL, bodyStr)

	if err != nil {
		log.Fatalf("%v\n", err.Error())
		return
	}
	res, queryErr := client.Do(req)
	defer res.Body.Close()
	client.CloseIdleConnections()
	if queryErr != nil {
		log.Fatalf("%v\n", queryErr.Error())
		return
	}
	var queryResult LCUserRatingGraphQLResult
	if resDecodeErr := json.NewDecoder(res.Body).Decode(&queryResult); resDecodeErr != nil {
		log.Fatalf("%v\n", resDecodeErr.Error())
		return
	}
	userInfo.UserName = userName
	userInfo.Rating = queryResult.Data.UserContestRanking.Rating
	userInfo.AttendedContestsCount = queryResult.Data.UserContestRanking.AttendedContestsCount
	return
}
