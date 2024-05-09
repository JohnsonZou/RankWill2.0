package controller

import (
	dao2 "RankWillServer/dao"
	"RankWillServer/web_server/common"
	"RankWillServer/web_server/dto"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
)

type allContest []dao2.Contest

func (a allContest) Len() int {
	return len(a)
}
func (a allContest) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a allContest) Less(i, j int) bool {
	return a[i].StartTime > a[j].StartTime
}
func GetContest(c *gin.Context) {
	page, err := strconv.Atoi(c.PostForm("page"))
	if page <= 0 || err != nil {
		common.Fail(c, nil, "page err!")
		return
	}
	var tot []dao2.Contest
	db := dao2.GetDB()
	_ = db.Find(&tot)
	sort.Sort(allContest(tot))
	total := len(tot)

	const pagesize = 10

	if ((page - 1) * pagesize) > total {
		common.Success(c, gin.H{"totnum": total, "data": nil}, "there is no such page")
	} else {
		common.Success(c, gin.H{"totnum": total, "data": dto.ToContestDto(tot[(page-1)*pagesize : min(page*pagesize, total)])}, "query successfully")
	}

}
