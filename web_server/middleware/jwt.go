package middleware

import (
	"RankWillServer/dao"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func GetJWTKey() string {
	viper := viper.New()
	viper.SetConfigName("web_server_config")
	viper.SetConfigType("yaml")
	dir, _ := os.Getwd()
	viper.AddConfigPath(dir + "\\config\\")
	if err := viper.ReadInConfig(); err != nil {
		return ""
	}
	log.Println(viper.GetString("jwt_key"))
	return viper.GetString("jwt_key")
}

func ReleaseToken(user dao.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "rankwill",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(GetJWTKey()))
	if err != nil {
		return "", err
	}
	return tokenString, err
}
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(GetJWTKey()), nil
	})
	return token, claims, err
}
