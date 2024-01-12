package util

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func GenPostReq(url string, body string) (*http.Request, error) {
	var req *http.Request
	retryErr := Retry(100, 10*time.Millisecond, func() error {
		var err error
		req, err = http.NewRequest("POST", url, strings.NewReader(body))
		return err
	})
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)
	return req, retryErr
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
func CloseResponseBody(res *http.Response) {
	if err := Retry(5, 100*time.Millisecond, func() error {
		return res.Body.Close()
	}); err != nil {
		log.Printf("%v\n", err.Error())
	}
}
