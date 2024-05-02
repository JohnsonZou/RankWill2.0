package model

import (
	"RankWillServer/backend/util"
	"net/http"
)

const (
	I18N_LCGraphQLURL                         = "https://leetcode.com/graphql/"
	I18N_UserRatingGraphQLQueryPostBodyPrefix = "{\"query\":\"\\n    query userContestRankingInfo($username: String!) {\\n  userContestRanking(username: $username) {\\n rating\\n   attendedContestsCount\\n }\\n \\n}\\n\",\"variables\":{\"username\":\""
	I18N_UserRatingGraphQLQueryPostBodySuffix = "\"}}"

	CN_LCGraphQLURL                         = "https://leetcode.cn/graphql/noj-go/"
	CN_UserRatingGraphQLQueryPostBodyPrefix = "{\"query\":\"query userContestRankingInfo($userSlug: String!) {  userContestRanking(userSlug: $userSlug) {     rating    attendedContestsCount  }  }    \",\"variables\":{\"userSlug\":\""
	CN_UserRatingGraphQLQueryPostBodySuffix = "\"}}"
)

type Submission struct {
	Id           int64  `json:"id"`
	Date         int64  `json:"date"`
	QuestionId   int64  `json:"question_id"`
	SubmissionId int64  `json:"submission_id"`
	Status       int    `json:"status"`
	ContestId    int    `json:"contest_id"`
	DataRegion   string `json:"data_region"`
	FailCount    int    `json:"fail_count"`
	Lang         string `json:"lang"`
}

func buildCNUserRatingGraphQLQueryPostBody(userName string) string {
	return CN_UserRatingGraphQLQueryPostBodyPrefix +
		userName + CN_UserRatingGraphQLQueryPostBodySuffix
}
func buildI18NUserRatingGraphQLQueryPostBody(userName string) string {
	return I18N_UserRatingGraphQLQueryPostBodyPrefix +
		userName + I18N_UserRatingGraphQLQueryPostBodySuffix
}

func (user *UserRankInfo) BuildUserRatingGraphQLQueryPostRequest() (*http.Request, error) {
	if user.DataRegion == "CN" {
		return util.GenPostReq(CN_LCGraphQLURL, buildCNUserRatingGraphQLQueryPostBody(user.Username))
	} else {
		return util.GenPostReq(I18N_LCGraphQLURL, buildI18NUserRatingGraphQLQueryPostBody(user.Username))
	}
}

type LCUserRatingGraphQLResultInfo struct {
	Rating                float64 `json:"rating"`
	AttendedContestsCount int     `json:"attendedContestsCount"`
}
type LCUserRatingGraphQLResultData struct {
	UserContestRanking LCUserRatingGraphQLResultInfo `json:"userContestRanking"`
}
type LCUserRatingGraphQLResult struct {
	Data LCUserRatingGraphQLResultData `json:"data"`
}
type LCContestInfoGraphQLResultData struct {
	PastContests LCContestInfoGraphQLResultPastContest `json:"pastContests"`
}

type LCContestInfoGraphQLResultPastContest struct {
	Data []Contest `json:"data"`
}
type LCContestInfoGraphQLResult struct {
	Data LCContestInfoGraphQLResultData `json:"data"`
}

type LCCNQueryUpComingContestData struct {
	ContestUpcomingContests []*Contest `json:"contestUpcomingContests"`
}
type LCCNQueryUpComingContest struct {
	Data LCCNQueryUpComingContestData `json:"data"`
}
