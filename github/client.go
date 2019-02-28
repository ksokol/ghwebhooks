package github

import (
	"encoding/json"
	"fmt"
	"ghwebhooks/config"
	"io"
	"net/http"
	"time"
)

const (
	shortDuration = 5
	longDuration  = 120
)

func Get(url string, v interface{}) error {
	return exchange("GET", url, v)
}

func Delete(url string) error {
	return exchange("DELETE", url, nil)
}

func Download(url string, dst io.Writer) error {
	var resp, err = client(longDuration).Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(dst, resp.Body)

	return err
}

func exchange(method string, url string, v interface{}) error {
	if req, err := newRequest(method, url); err != nil {
		return err
	} else {
		if resp, err := client(shortDuration).Do(req); err == nil {
			err = parseBody(resp.Body, v)
			defer resp.Body.Close()
		}
		return err
	}
}

func newRequest(method string, url string) (*http.Request, error) {
	if req, err := http.NewRequest(method, url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetAccessToken()))
		return req, err
	}
}

func client(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout * time.Second),
	}
}

func parseBody(payload io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(payload)
	return decoder.Decode(v)
}
