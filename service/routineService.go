package service

import (
	"RankWillServer/util"
	"context"
	"time"
)

func RoutineService(ctx context.Context) {
	scheduler := util.GetMainScheduler(ctx)
	scheduler.Every(15 * time.Minute).Do(FindAndRegisterContest(ctx))

}

func consumeMQMsg(ctx context.Context) {
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
		c := Contest{
			TitleSlug: string(d.Body),
		}
		err := c.HandleContest(ctx)
		if err != nil {
			panic(err)
		}
	}
}
