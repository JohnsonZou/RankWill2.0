package dto

import (
	"RankWillServer/dao"
)

type UserDto struct {
	Email string `json:"email"`
}

func ToUserDto(user dao.User) UserDto {
	return UserDto{
		Email: user.Email,
	}
}
