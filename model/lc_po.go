package model

import (
	"RankWillServer/util"
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

type Contest struct {
	ContestName   string
	ContestantNum int
	PageNum       int
	RankPages     map[int]*RankPage
}
type RankPage struct {
	Time        float64                 `json:"time"`
	Submissions []map[string]Submission `json:"submissions"`
	Total_Rank  []UserRankInfo          `json:"total_rank"`
	UserNum     int                     `json:"user_num"`
}
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
type UserRankInfo struct {
	ContestId     int    `json:"contest_id"`
	Username      string `json:"username"`
	UserSlug      string `json:"user_slug"`
	CountryCode   string `json:"country_code"`
	CountryName   string `json:"country_name"`
	Rank          int    `json:"rank"`
	Score         int    `json:"score"`
	FinishTime    int64  `json:"finish_time"`
	GlobalRanking int    `json:"global_ranking"`
	DataRegion    string `json:"data_region"`
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
