package model

import "sync"

type Contest struct {
	Title           string `json:"title"`
	TitleSlug       string `json:"titleSlug"`
	StartTime       int64  `json:"startTime"`
	OriginStartTime int64  `json:"originStartTime"`
	Duration        int64  `json:"duration"`
	ContestantNum   int
	PageNum         int
	RankPages       map[int]*RankPage
	Lock            sync.Mutex
}
type RankPage struct {
	Time        float64                 `json:"time"`
	Submissions []map[string]submission `json:"submissions"`
	TotalRank   []UserRankInfo          `json:"total_rank"`
	UserNum     int                     `json:"user_num"`
}
type submission struct {
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
	ContestId             int     `json:"contest_id"`
	Username              string  `json:"username"`
	UserSlug              string  `json:"user_slug"`
	CountryCode           string  `json:"country_code"`
	CountryName           string  `json:"country_name"`
	Rank                  int     `json:"rank"`
	Score                 int     `json:"score"`
	FinishTime            int64   `json:"finish_time"`
	GlobalRanking         int     `json:"global_ranking"`
	DataRegion            string  `json:"data_region"`
	AttendedContestsCount int     `json:"attendedContestsCount"`
	Rating                float64 `json:"rating"`
	PredictedRating       float64 `json:"predict_rating"`
	ContestName           string  `json:"contest_name"`
}

type MQMessage struct {
	ContestName   string `json:"contestName"`
	PredictNeeded bool   `json:"predictNeeded"`
}
