package dto

import (
	"RankWillServer/dao"
)

type FollowDto struct {
	Lcusername string `json:"lcusername"`
}

func ToFollowDto(f []dao.Following) []FollowDto {
	var res []FollowDto
	for _, k := range f {
		res = append(res, FollowDto{
			Lcusername: k.Lcusername,
		})
	}
	return res
}
