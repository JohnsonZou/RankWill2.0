package dto

import (
	"RankWillServer/dao"
)

type ContestDto struct {
	ContestName   string `json:"contestName"`
	UpdateTime    int64  `json:"updateTime"`
	ContestantNum int    `json:"contestantNum"`
}

func ToContestDto(c []dao.Contest) []ContestDto {
	var res []ContestDto
	for _, v := range c {
		res = append(res, ContestDto{
			ContestName:   v.TitleSlug,
			UpdateTime:    v.StartTime,
			ContestantNum: v.ContestantNum,
		})
	}
	return res
}
