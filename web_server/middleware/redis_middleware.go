package middleware

import (
	"RankWillServer/redis"
	"context"

	"github.com/gin-gonic/gin"
)

func RedisMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, err := redis.InitRedisClient(context.Background())
		if err != nil {
			panic(err)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
