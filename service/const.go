package service

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
)

const (
	cookie = "csrftoken=1TCqj8frSz5I2tPQhFSbLhyx9vWXFhPAxdrap9ezE5GTHC3lEWslq4ZUr9bmdxKO"
)
