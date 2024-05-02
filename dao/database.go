package dao

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	viper := viper.New()
	viper.SetConfigName("mysql_config")
	viper.SetConfigType("yaml")
	dir, _ := os.Getwd()
	viper.AddConfigPath(dir + "\\config\\")
	if err := viper.ReadInConfig(); err != nil {
		return nil
	}
	host := viper.GetString("host")
	port := viper.GetString("port")
	database := viper.GetString("database")
	username := viper.GetString("username")
	password := viper.GetString("password")
	charset := viper.GetString("charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(mysql.Open(args))
	if err != nil {
		panic("fail to connect database,err: " + err.Error())
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Contest{})
	db.AutoMigrate(&Contestant{})
	db.AutoMigrate(&Following{})
	DB = db
	return db
}
func GetDB() *gorm.DB {
	return DB
}
