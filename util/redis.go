package util

import (
	"strconv"
	"strings"
)

func BuildRedisContestKey(contestName string) string {
	return "ContestKey###" + contestName
}
func BuildRedisFetchedContestantSKey(contestID int, uname string) string {
	return strconv.Itoa(contestID) + "###" + uname
}
func BuildRedisContestantSVal(rating float64, attendedContestCount int) string {
	return strconv.FormatFloat(rating, 'f', -1, 64) + "#" + strconv.Itoa(attendedContestCount)
}

func BuildRedisPredictedContestantSKey(contestID int, uname string) string {
	return strconv.Itoa(contestID) + "######" + uname
}
func ParseContestSVal(key string) (rating float64, attendedContestCount int64, err error) {
	strArr := strings.Split(key, "#")
	rating, err = strconv.ParseFloat(strArr[0], 64)
	if err != nil {
		return
	}
	attendedContestCount, err = strconv.ParseInt(strArr[1], 10, 32)
	return
}
