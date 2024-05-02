package service

import (
	"context"
	"log"
	"math"

	"github.com/go-redis/redis"
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
	ERSlice := calcPossibleExpectedRankForAllRatings(ratings)
	for _, p := range contest.rankPages {
		for _, u := range p.TotalRank {

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

			log.Println("[predicted] ", u.Rating, float64(l)*predictRatingDeviationDelta, ERSlice[int(pos)], u.Rank)
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
