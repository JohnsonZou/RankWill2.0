package util

import (
	"bytes"
	"io"
	"log"
	"time"
)

func Retry(maxTime int, gapDuration time.Duration, f func() error) error {
	err := f()
	for err != nil && maxTime > 1 {
		maxTime--
		time.Sleep(gapDuration)
		err = f()

		log.Println(maxTime, err)
	}
	return err
}

func ReadCloserToString(rc io.ReadCloser) string {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(rc)
	if err != nil {
		return ""
	}
	// Convert the bytes.Buffer to a string and return it
	return buf.String()
}
