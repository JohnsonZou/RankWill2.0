package service

import (
	"RankWillServer/backend/model"
	"RankWillServer/backend/util"
	"context"
	"encoding/json"
	"sync"
	"time"
)

func RoutineService(ctx context.Context) {
	scheduler := util.GetMainScheduler(ctx)
	scheduler.Every(15 * time.Minute).Do(FindAndRegisterContest(ctx))
}

var lock sync.Mutex

func ConsumeMQMsg(ctx context.Context) {
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
		lock.Lock()

		msg := model.MQMessage{}
		_ = json.Unmarshal(d.Body, &msg)

		c := model.Contest{
			TitleSlug: msg.ContestName,
		}

		_ = FetchContest(ctx, &c)

		if msg.PredictNeeded {
			_ = Predict(ctx, &c)
		}
		lock.Unlock()
	}
}
