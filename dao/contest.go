package dao

import (
	"RankWillServer/backend/model"
	"log"
	"time"

	"gorm.io/gorm"
)

type Contest struct {
	gorm.Model
	TitleSlug     string `gorm:"type:varchar(30);not null"`
	StartTime     int64  `gorm:"type:bigint"`
	ContestantNum int    `gorm:"type:int"`
}

func isContestantExisted(db *gorm.DB, con Contestant) bool {
	var c Contestant
	res := db.Where("contestname=?", con.Contestname)
	res = res.Where("username=?", con.Username)
	res.First(&c)
	return c.ID != 0
}

func isContestExisted(contestName string) bool {
	db := GetDB()
	var a Contest
	db.Where("title_slug=?", contestName).First(&a)
	return a.ID != 0
}
func InsertContestIntoDB(contest *model.Contest) {
	log.Println("Start to inser into DB")
	db := GetDB()
	contestName := contest.TitleSlug
	for _, p := range contest.RankPages {
		for _, u := range p.TotalRank {
			v := Contestant{
				Contestname:          contestName,
				Username:             u.UserSlug,
				Rank:                 u.Rank,
				FinishTime:           u.FinishTime,
				DataRegion:           u.DataRegion,
				AttendedContestCount: u.AttendedContestsCount,
				Score:                u.Score,
				Rating:               u.Rating,
				PredictedRating:      u.PredictedRating,
			}

			if !isContestantExisted(db, v) {
				db.Create(&v)
			} else {
				db.Where("contestname=?", contestName).Where("username=?", v.Username).Updates(&v)
			}
		}
	}
	a := Contest{
		StartTime:     time.Now().Unix(),
		TitleSlug:     contestName,
		ContestantNum: contest.ContestantNum,
	}
	if isContestExisted(contestName) {
		db.Where("title_slug=?", contestName).Updates(&a)
	} else {
		db.Create(&a)
	}
}
