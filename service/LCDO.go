package service

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

type submission struct {
	Id           int64  `json:"id"`
	Date         int64  `json:"date"`
	QuestionId   int64  `json:"question_id"`
	SubmissionId int64  `json:"submission_id"`
	Status       int16  `json:"status"`
	ContestId    int16  `json:"contest_id"`
	DataRegion   string `json:"data_region"`
	FailCount    int16  `json:"fail_count"`
	Lang         string `json:"lang"`
}
type userRank struct {
	ContestId     int16  `json:"contest_id"`
	Username      string `json:"username"`
	UserSlug      string `json:"user_slug"`
	CountryCode   string `json:"country_code"`
	CountryName   string `json:"country_name"`
	Rank          int32  `json:"rank"`
	Score         int16  `json:"score"`
	FinishTime    int64  `json:"finish_time"`
	GlobalRanking int32  `json:"global_ranking"`
	DataRegion    string `json:"data_region"`
}
type RankPage struct {
	Time        float64                 `json:"time"`
	Submissions []map[string]submission `json:"submissions"`
	Total_Rank  []userRank              `json:"total_rank"`
	UserNum     int32                   `json:"user_num"`
}
