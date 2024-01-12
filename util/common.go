package util

import "time"

func Retry(maxTime int, gapDuration time.Duration, f func() error) error {
	err := f()
	for err != nil && maxTime > 1 {
		maxTime--
		time.Sleep(gapDuration)
		err = f()
	}
	return err
}
