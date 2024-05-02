package test

import (
	"RankWillServer/backend/model"
	"RankWillServer/backend/service"
	"RankWillServer/backend/util"
	"context"
)

func MainTest(ctx context.Context) {
	ctx = util.SetTestMode(ctx)
	c := model.Contest{
		TitleSlug: "biweekly-contest-125",
	}
	err := service.FetchContest(ctx, &c)

	if err != nil {
		panic(err)
	}

}
