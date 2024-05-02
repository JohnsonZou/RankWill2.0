package service

import "errors"

var (
	errStatusCodeNot200 = errors.New("status code not 200")
	errResNil           = errors.New("nil query result")
)

const (
	defaultUserRating           = 1500.0
	predictRatingMinimum        = 0.0
	predictRatingMaximum        = 4000.0
	predictRatingDeviationDelta = 0.5
)

const (
	maxGoroutineNum = 7

	stopTimeMilliSec = 10000
)

const (
	I18N_LCGraphQLURL                         = "https://leetcode.com/graphql/"
	I18N_UserRatingGraphQLQueryPostBodyPrefix = "{\"query\":\"\\n    query userContestRankingInfo($username: String!) {\\n  userContestRanking(username: $username) {\\n rating\\n   attendedContestsCount\\n }\\n \\n}\\n\",\"variables\":{\"username\":\""
	I18N_UserRatingGraphQLQueryPostBodySuffix = "\"}}"

	CN_LCGraphQLURL                         = "https://leetcode.cn/graphql/noj-go/"
	CN_LCGraphQLURL_SHORT                   = "https://leetcode.cn/graphql/"
	CN_UserRatingGraphQLQueryPostBodyPrefix = "{\"query\":\"query userContestRankingInfo($userSlug: String!) {  userContestRanking(userSlug: $userSlug) {     rating    attendedContestsCount  }  }    \",\"variables\":{\"userSlug\":\""
	CN_UserRatingGraphQLQueryPostBodySuffix = "\"}}"
	CN_UpComingContestGraphQLQueryPostBody  = "{\"operationName\":null,\"variables\":{},\"query\":\"{\\n  contestUpcomingContests {\\n    containsPremium\\n    title\\n    cardImg\\n    titleSlug\\n    description\\n    startTime\\n    duration\\n    originStartTime\\n    isVirtual\\n    isLightCardFontColor\\n    company {\\n      watermark\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\"}"
	I18N_LCContestRankQueryPrefix           = "https://leetcode.com/contest/api/ranking/"

	CN_LCContestRankQueryPrefix = "https://leetcode.cn/contest/api/ranking/"
	LCContestRankQueryMidfix    = "/?pagination="
	LCContestRankQuerySuffix    = "&region=global"

	I18N_LCContestQueryGraphQLPostBody = "\n{\"query\":\"\\n    query pastContests($pageNo: Int, $numPerPage: Int) {\\n  pastContests(pageNo: $pageNo, numPerPage: $numPerPage) {\\n    pageNum\\n    currentPage\\n    totalNum\\n    numPerPage\\n    data {\\n      title\\n      titleSlug\\n      startTime\\n      originStartTime\\n       }\\n  }\\n}\\n    \",\"variables\":{\"pageNo\":1}}"
)

const (
	testPageNum = 100
)
