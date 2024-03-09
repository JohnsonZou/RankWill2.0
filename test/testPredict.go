package test

import (
	"RankWillServer/service"
	"RankWillServer/util"
	"context"
)

func MainTest(ctx context.Context) {
	ctx = util.SetTestMode(ctx)
	c := service.Contest{
		TitleSlug: "biweekly-contest-125",
	}
	err := c.HandleContestTest(ctx)

	if err != nil {
		panic(err)
	}
	err = c.Predict(ctx)
	if err != nil {
		panic(err)
	}
}
