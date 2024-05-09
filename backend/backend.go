package backend

import (
	"RankWillServer/backend/mq"
	"RankWillServer/backend/service"
	"RankWillServer/backend/util"
	myredis "RankWillServer/redis"
	"context"
)

func InitContext(ctx context.Context) context.Context {
	ContestChanel := make(chan int, 2000)
	ctx = context.WithValue(ctx, util.ContestChanelKey, ContestChanel)
	ctx, err := myredis.InitRedisClient(ctx)
	if err != nil {
		panic(err)
	}
	ctx = util.InitMainScheduler(ctx)

	cli := util.GemNewHTTPClient()
	ctx = util.SetHttpClient(ctx, cli)
	ctx, err = mq.InitMQChanel(ctx)

	if err != nil {
		panic(err)
	}
	ctx = util.InitTimer(ctx)
	return ctx
}

func Serve() {
	ctx := InitContext(context.Background())
	go service.RoutineService(ctx)
	go service.ConsumeMQMsg(ctx)
}
