package service

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
type LCUserInfo struct {
	Rating                float64 `json:"rating"`
	UserName              string  `json:"userName"`
	UserSlug              string  `json:"userSlug"`
	AttendedContestsCount int     `json:"attendedContestsCount"`
}
