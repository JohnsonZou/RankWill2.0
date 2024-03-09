package main

import (
	"RankWillServer/service"
	"RankWillServer/test"
	"RankWillServer/util"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
func initContext(ctx context.Context) context.Context {
	ContestChanel := make(chan int, 2000)
	ctx = context.WithValue(ctx, util.ContestChanelKey, ContestChanel)
	ctx, err := util.InitRedisClient(ctx)
	if err != nil {
		panic(err)
	}
	ctx = util.InitMainScheduler(ctx)

	cli := util.GenNewClient()
	ctx = util.SetHttpClient(ctx, cli)
	ctx, err = util.InitMQChanel(ctx)

	if err != nil {
		panic(err)
	}
	ctx = util.InitTimer(ctx)
	return ctx
}
func main() {

	ctx := initContext(context.Background())

	go service.RoutineService(ctx)
	test.MainTest(ctx)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	fmt.Println(<-sig)
}

func t(ctx context.Context) {
	ch := util.GetMainMQChanel(ctx)
	msgs, _ := ch.Consume(
		"delayedQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	for d := range msgs {
		log.Printf("[MQ]Received a message: %s", d.Body)
		//contestTile := d.Body

	}
}
