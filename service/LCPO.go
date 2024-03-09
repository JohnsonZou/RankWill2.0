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
