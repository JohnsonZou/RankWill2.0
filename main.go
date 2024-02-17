package main

import (
	"RankWillServer/util"
	"context"
	"log"
)

func init() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
func main() {
	ctx := context.Background()
	err := util.InitRedisClient(ctx)
	if err != nil {
		println(err.Error())
	}
}
