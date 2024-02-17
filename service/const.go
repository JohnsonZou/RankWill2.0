package service

import "errors"

var (
	errStatusCodeNot200 = errors.New("status code not 200")
)

const (
	maxGoroutineNum = 15
)

const (
	I18N_LCGraphQLURL                         = "https://leetcode.com/graphql/"
	I18N_UserRatingGraphQLQueryPostBodyPrefix = "{\"query\":\"\\n    query userContestRankingInfo($username: String!) {\\n  userContestRanking(username: $username) {\\n rating\\n   attendedContestsCount\\n }\\n \\n}\\n\",\"variables\":{\"username\":\""
	I18N_UserRatingGraphQLQueryPostBodySuffix = "\"}}"

	CN_LCGraphQLURL                         = "https://leetcode.cn/graphql/noj-go/"
	CN_UserRatingGraphQLQueryPostBodyPrefix = "{\"query\":\"query userContestRankingInfo($userSlug: String!) {  userContestRanking(userSlug: $userSlug) {     rating    attendedContestsCount  }  }    \",\"variables\":{\"userSlug\":\""
	CN_UserRatingGraphQLQueryPostBodySuffix = "\"}}"

	I18N_LCContestRankQueryPrefix = "https://leetcode.com/contest/api/ranking/"
	I18N_LCContestRankQueryMidfix = "/?pagination="
	I18N_LCContestRankQuerySuffix = "&region=global"
)
