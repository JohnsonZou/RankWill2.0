package util

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
)

var lock sync.Mutex

func GetHttpClient(ctx context.Context) *http.Client {
	cli := ctx.Value(httpClientKey)
	res, ok := cli.(*http.Client)
	if ok {
		return res
	}
	return nil
}

func GetChanelFromCtxByKey(ctx context.Context, key string) chan int {
	ch := ctx.Value(key)
	res, ok := ch.(chan int)
	if ok {
		return res
	}
	return nil
}
func SetTestMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, TestModeKey, "t")
}
func IsTestMode(ctx context.Context) bool {
	t := ctx.Value(TestModeKey)
	_, ok := t.(string)
	return ok
}
func SetHttpClient(ctx context.Context, cli *http.Client) context.Context {
	return context.WithValue(ctx, httpClientKey, cli)
}
func InitMainScheduler(ctx context.Context) context.Context {
	s := gocron.NewScheduler(time.Local)
	s.StartAsync()
	return context.WithValue(ctx, MainSchedulerKey, s)
}
func GetMainMQChanel(ctx context.Context) *amqp.Channel {
	ch := ctx.Value(MainMQChanelKey)
	res, ok := ch.(*amqp.Channel)
	if ok {
		return res
	}

	return nil
}
func GetMainScheduler(ctx context.Context) *gocron.Scheduler {
	cli := ctx.Value(MainSchedulerKey)
	res, ok := cli.(*gocron.Scheduler)
	if ok {
		return res
	}
	return nil
}

func InitTimer(ctx context.Context) context.Context {
	cntMap := make(map[string]int)
	return context.WithValue(context.WithValue(ctx, StartTimeKey, time.Now().UnixMilli()), QueryCounterKey, cntMap)
}

func AddCounterAndGetSpeed(ctx context.Context) (float64, int) {
	cnt := ctx.Value(QueryCounterKey).(map[string]int)

	lock.Lock()
	cnt[QueryCounterKey]++
	tmp := cnt[QueryCounterKey]
	lock.Unlock()

	deltaT := time.Now().UnixMilli() - ctx.Value(StartTimeKey).(int64)
	return 1000 * float64(tmp) / float64(deltaT), tmp
}
