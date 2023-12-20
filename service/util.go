package service

import (
	"net/http"
	"strings"
	"time"
)

func GenPostReq(url string, body string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)
	return req, err
}

func Retry(maxTime int, gapDuration time.Duration, f func() error) error {
	err := f()
	for err != nil && maxTime > 1 {
		maxTime--
		time.Sleep(gapDuration)
		err = f()
	}
	return err
}

func GenGetReq(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cookie", cookie)
	return req, err
}
func GenNewClient() http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()

	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	t.DisableKeepAlives = true
	return http.Client{
		Transport: t,
		Timeout:   3 * time.Second,
	}
}
