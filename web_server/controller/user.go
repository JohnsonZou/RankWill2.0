package controller

import (
	dao2 "RankWillServer/dao"
	"RankWillServer/web_server/common"
	"RankWillServer/web_server/dto"
	"RankWillServer/web_server/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"regexp"
)

func isEmailExisted(db *gorm.DB, email string) bool {
	var user dao2.User
	db.Where("email=?", email).First(&user)
	return user.ID != 0
}
func isFollowExisted(db *gorm.DB, uname string, lcusername string) bool {
	var fl dao2.Following
	db.Where("email=?", uname).Where("lcusername=?", lcusername).First(&fl)
	return fl.ID != 0
}
func getUserByEmail(db *gorm.DB, email string) dao2.User {
	var user dao2.User
	db.Where("email=?", email).First(&user)
	return user
}
func validEmail(email string) (bool, error) {
	regex := "^([a-z0-9A-Z]+[-|\\.]?)+[a-z0-9A-Z]@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?\\.)+[a-zA-Z]{2,}$"
	return regexp.MatchString(regex, email)
}
func Register(c *gin.Context) {
	db := dao2.GetDB()
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	log.Println(username, email, password)

	if res, matchErr := validEmail(email); res == false || matchErr != nil {
		if matchErr != nil {
			common.Response(c, http.StatusInternalServerError, 500, nil, "Email matching fail")
		}
		common.Response(c, http.StatusUnprocessableEntity, 422, nil, "Email invalid")
		return
	}
	if isEmailExisted(db, email) {
		common.Response(c, http.StatusUnprocessableEntity, 422, nil, "Register failed,email existed.")
		return
	}
	newUser := dao2.User{
		Email:    email,
		Password: password,
	}
	db.Create(&newUser)
	common.Success(c, nil, "Successfully register")
}
func Login(c *gin.Context) {
	db := dao2.GetDB()
	email := c.PostForm("email")
	password := c.PostForm("password")
	log.Println(email, password)
	loginUser := getUserByEmail(db, email)
	if loginUser.ID == 0 {
		common.Response(c, http.StatusUnprocessableEntity, 422, nil, "Login failed,email not exist")
		return
	}
	if loginUser.Password != password {
		common.Fail(c, nil, "Wrong password")
		return
	}
	token, tokenGenErr := middleware.ReleaseToken(loginUser)
	if tokenGenErr != nil {
		common.Response(c, http.StatusInternalServerError, 500, nil, "token generation failed")
		log.Println("token generation failed", tokenGenErr.Error())
		return
	}
	common.Success(c, gin.H{"token": token}, "Successfully login")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	common.Success(c, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(dao2.User))}}, "UserInfo request successfully")
}
