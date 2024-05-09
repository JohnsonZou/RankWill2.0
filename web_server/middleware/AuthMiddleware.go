package middleware

import (
	dao2 "RankWillServer/dao"
	"RankWillServer/web_server/common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		log.Println(tokenString)
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			log.Println(tokenString)
			common.Response(c, http.StatusUnauthorized, 401, nil, "Unauthorized token")
			c.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			log.Println(err)
			common.Response(c, http.StatusUnauthorized, 401, nil, "Unauthorized token")
			c.Abort()
			return
		}
		userId := claims.UserId
		DB := dao2.GetDB()
		var user dao2.User
		DB.First(&user, userId)
		if user.ID == 0 {
			log.Println(user.ID)
			common.Response(c, http.StatusUnauthorized, 401, nil, "Unauthorized token")
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
