package web_server

import (
	"RankWillServer/web_server/controller"
	"RankWillServer/web_server/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "gorm.io/driver/mysql"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.POST("/api/querypage", controller.Getpage)
	r.POST("/api/querybyname", controller.Getbyname)
	r.POST("/api/auth/follow", middleware.AuthMiddleware(), controller.Follow)
	r.POST("/api/auth/unfollow", middleware.AuthMiddleware(), controller.Unfollow)
	r.POST("/api/auth/getfollowing", middleware.AuthMiddleware(), controller.GetFollowing)
	r.POST("/api/getcontestbypage", controller.GetContest)
	r.POST("/api/auth/getfollowlist", middleware.AuthMiddleware(), controller.GetFollowList)
	return r
}

func GinRun() {
	viper := viper.New()
	viper.SetConfigName("web_server_config")
	viper.SetConfigType("yaml")
	dir, _ := os.Getwd()
	viper.AddConfigPath(dir + "\\config\\")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	port := viper.GetString("port")

	r := gin.Default()
	r.Use(middleware.RedisMiddleware())
	r.Use(middleware.CORS())
	r = CollectRoute(r)

	panic(r.Run(":" + port))
}
