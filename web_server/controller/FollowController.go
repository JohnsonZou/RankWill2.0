package controller

import (
	dao2 "RankWillServer/dao"
	"RankWillServer/web_server/common"
	"RankWillServer/web_server/dto"
	"github.com/gin-gonic/gin"
)

func Follow(c *gin.Context) {
	user, _ := c.Get("user")
	email := user.(dao2.User).Email
	lcusername := c.PostForm("username")
	if lcusername == "" {
		common.Fail(c, nil, "Empty leetcode username")
		return
	}
	db := dao2.GetDB()
	if isFollowExisted(db, email, lcusername) {
		common.Fail(c, nil, "Duplicated follow")
		return
	}
	f := dao2.Following{
		Email:      email,
		Lcusername: lcusername,
	}
	db.Create(&f)
	common.Success(c, nil, "Successfully follow")
}
func Unfollow(c *gin.Context) {
	user, _ := c.Get("user")
	email := user.(dao2.User).Email
	lcusername := c.PostForm("username")
	if lcusername == "" {
		common.Fail(c, nil, "Empty leetcode username")
		return
	}
	db := dao2.GetDB()
	if isFollowExisted(db, email, lcusername) {
		db.Where("email=?", email).Where("lcusername=?", lcusername).Delete(&dao2.Following{})
		common.Success(c, nil, "Successfully unfollow")
		return
	}
	common.Fail(c, nil, "Leetcode user not exist")
}
func GetFollowList(c *gin.Context) {
	user, _ := c.Get("user")
	email := user.(dao2.User).Email
	db := dao2.GetDB()
	var fol []dao2.Following
	db.Where("email=?", email).Find(&fol)
	common.Success(c, gin.H{"result": dto.ToFollowDto(fol)}, "Successfully get followlist")

}
func GetFollowing(c *gin.Context) {
	user, _ := c.Get("user")
	email := user.(dao2.User).Email
	contestname := c.PostForm("contestname")
	if contestname == "" {
		common.Fail(c, nil, "Empty contest name")
		return
	}
	db := dao2.GetDB()
	var res []dao2.Contestant
	var fol []dao2.Following
	db.Where("email=?", email).Find(&fol)
	for _, v := range fol {
		curname := v.Lcusername
		db.Where("contestname=?", contestname).Where("username=?", curname).Find(&res)
	}

	common.Success(c, gin.H{"result": dto.ToQueryPageDto(res)}, "Successfully get following")
}
