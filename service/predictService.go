package service

import (
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

func (user *userRankInfo) lcFixDelta() float64 {
	return 1.0 / (1.0 + (1-math.Pow(5.0/7.0, float64(user.AttendedContestsCount+1)))*3.5)
}

func (contest *Contest) Predict(ctx context.Context) error {
	var ratings []float64
	for _, p := range contest.rankPages {
		for i, _ := range p.TotalRank {
			err := p.TotalRank[i].getUserFetchRatingInfo(ctx)
			if err == redis.Nil {
				p.TotalRank[i].Rating = defaultUserRating //cover
			}
			ratings = append(ratings, p.TotalRank[i].Rating)
		}
	}
	cnt := 0
	ERSlice := calcPossibleExpectedRankForAllRatings(ratings)
	for _, p := range contest.rankPages {
		for _, u := range p.TotalRank {
			cnt++
			l := 0
			r := len(ERSlice) - 1

			for l < r {
				mid := (l + r) / 2
				if ERSlice[mid] > u.Rating {
					r = mid
				} else {
					l = mid + 1
				}
			}
			u.Rating = u.Rating + u.lcFixDelta()*(float64(l)*predictRatingDeviationDelta-u.Rating)
			u.AttendedContestsCount++
			err := u.setUserPredictedRatingInfo(ctx)
			if err != nil {
				//log
				log.Println("[Error][Predict]", u.UserSlug)
			}
		}
	}
	return nil
}
