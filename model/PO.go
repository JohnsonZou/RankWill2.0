package model

type LCUserInfo struct {
	Rating                float64 `json:"rating"`
	UserName              string  `json:"userName"`
	UserSlug              string  `json:"userSlug"`
	AttendedContestsCount int16   `json:"attendedContestsCount"`
}
