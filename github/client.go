package github

import (
	"encoding/json"
	"fmt"
	"ghwebhooks/config"
	"io"
	"net/http"
	"time"
)

var client http.Client

func init() {
	client = http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
}

func Get(url string, v interface{}) error {
	return exchange("GET", url, v)
}

func Delete(url string) error {
	return exchange("DELETE", url, nil)
}

func exchange(method string, url string, v interface{}) error {
	if req, err := newRequest(method, url); err != nil {
		return err
	} else {
		if resp, err := client.Do(req); err == nil {
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

func parseBody(payload io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(payload)
	return decoder.Decode(v)
}
