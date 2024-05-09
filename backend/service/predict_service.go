package service

import (
	"RankWillServer/backend/model"
	"RankWillServer/dao"
	myredis "RankWillServer/redis"
	"context"
	"github.com/go-redis/redis"
	"log"
	"math"
)

func calcPossibleExpectedRankForAllRatings(ratings []float64) []float64 {
	var res []float64
	for r := predictRatingMinimum; r <= predictRatingMaximum; r += predictRatingDeviationDelta {
		var curER float64 = 0
		for _, v := range ratings {
			curER += 1 / (1 + math.Pow(10, (r-v)/400))
		}
		res = append(res, curER)
	}
	return res
}

func lcFixDelta(user *model.UserRankInfo) float64 {
	return 1.0 / (1.0 + (1-math.Pow(5.0/7.0, float64(user.AttendedContestsCount+1)))*3.5)
}

func Predict(ctx context.Context, contest *model.Contest) error {
	var ratings []float64
	for _, p := range contest.RankPages {
		for i := range p.TotalRank {
			err := myredis.GetUserFetchRatingInfo(ctx, &p.TotalRank[i])
			if err == redis.Nil {
				p.TotalRank[i].Rating = defaultUserRating //cover
			}
			ratings = append(ratings, p.TotalRank[i].Rating)
		}
	}
	ERSlice := calcPossibleExpectedRankForAllRatings(ratings)
	for _, p := range contest.RankPages {
		for i, u := range p.TotalRank {

			pos := math.Ceil(u.Rating) / predictRatingDeviationDelta

			q_u := math.Sqrt(ERSlice[int(pos)] * float64(u.Rank))
			l := 0
			r := len(ERSlice) - 1

			for l < r {
				mid := (l + r) / 2
				if ERSlice[mid] <= q_u {
					r = mid
				} else {
					l = mid + 1
				}
			}

			p.TotalRank[i].PredictedRating = u.Rating + (float64(l)*predictRatingDeviationDelta-u.Rating)*lcFixDelta(&u)
			p.TotalRank[i].AttendedContestsCount++

			err := myredis.SetUserPredictedRatingInfo(ctx, &u)
			if err != nil {
				//log
				log.Println("[Error][Predict]", u.UserSlug)
			}
		}
	}
	dao.InsertIntoDB(contest)
	return nil
}
