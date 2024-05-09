package controller

import (
	dao2 "RankWillServer/dao"
	"RankWillServer/web_server/common"
	"RankWillServer/web_server/dto"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
)

type allContestant []dao2.Contestant

func (a allContestant) Len() int {
	return len(a)
}
func (a allContestant) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a allContestant) Less(i, j int) bool {
	return a[i].Rank < a[j].Rank
}
func Getpage(c *gin.Context) {
	db := dao2.GetDB()
	contestname := c.PostForm("contestname")
	page := c.PostForm("page")
	if contestname == "" || page == "" {
		common.Fail(c, gin.H{}, "invalid form data")
		return
	}
	var con []dao2.Contestant
	p, err := strconv.Atoi(page)
	if err != nil {
		panic(err)
	}
	var curContest dao2.Contest
	db.Where("title_slug=?", contestname).First(&curContest)
	if curContest.ID == 0 {
		common.Fail(c, nil, "no such contest")
		return
	}
	db.Where("`rank`>?", (p-1)*25).Where("`rank`<=?", p*25).Where("contestname=?", contestname).Find(&con)
	if con == nil {
		common.Fail(c, nil, "page empty")
		return
	}
	sort.Sort(allContestant(con))
	common.Success(c, gin.H{"result": dto.ToQueryPageDto(con), "contestantnum": curContest.ContestantNum}, "Successfully query page")
}
func Getbyname(c *gin.Context) {
	db := dao2.GetDB()
	contestname := c.PostForm("contestname")
	contestantname := c.PostForm("contestantname")
	if contestname == "" || contestantname == "" {
		common.Fail(c, gin.H{}, "invalid form data")
		return
	}
	var con dao2.Contestant
	db.Where("username=?", contestantname).Where("contestname=?", contestname).First(&con)
	if con.ID == 0 {
		common.Fail(c, gin.H{}, "no such user")
		return
	}
	common.Success(c, gin.H{"result": dto.ToQueryByNameDto(con)}, "Successfully query by name")
}
