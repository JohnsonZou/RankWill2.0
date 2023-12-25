package model

type LCUserRatingGraphQLResultInfo struct {
	Rating                float64 `json:"rating"`
	AttendedContestsCount int16   `json:"attendedContestsCount"`
}
type LCUserRatingGraphQLResultData struct {
	UserContestRanking LCUserRatingGraphQLResultInfo `json:"userContestRanking"`
}
type LCUserRatingGraphQLResult struct {
	Data LCUserRatingGraphQLResultData `json:"data"`
}
