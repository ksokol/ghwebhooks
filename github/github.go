package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ghwebhooks/config"
	"ghwebhooks/types"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Release struct {
	ID      int
	TagName string `json:"tag_name"`
	Assets  []struct {
		Url string `json:"browser_download_url"`
	}
}

type GithubEvent struct {
	Repository struct {
		Name       string `json:"name"`
		ReleaseUrl string `json:"releases_url"`
	}
	Release Release
}

func SupportsApp(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var githubEvent GithubEvent
		statusCode := 400
		buf, _ := ioutil.ReadAll(req.Body)
		bodyClone := ioutil.NopCloser(bytes.NewBuffer(buf))
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		if err := parsePayload(bodyClone, &githubEvent); err == nil {
			if _, err := config.GetAppConfig(githubEvent.Repository.Name); err == nil {
				handler(resp, req)
				return
			}

			statusCode = 404
		}

		resp.WriteHeader(statusCode)
	}
}

func Parse(body io.ReadCloser, status *types.Status) (GithubEvent, error) {
	var githubEvent GithubEvent

	if err := parsePayload(body, &githubEvent); err != nil {
		return githubEvent, err
	}

	name := githubEvent.Repository.Name

	if isInvalid(&githubEvent) {
		status.LogF("did not find a valid release for '%s'", name)

		if err := updateWithLatestRelease(&githubEvent, status); err != nil {
			return githubEvent, err
		}

		if isInvalid(&githubEvent) {
			return githubEvent, fmt.Errorf("release not found for '%s'", name)
		}
	}

	return githubEvent, nil
}

func isInvalid(githubEvent *GithubEvent) bool {
	return githubEvent.Release.ID == 0 || githubEvent.Repository.Name == ""
}

func updateWithLatestRelease(githubEvent *GithubEvent, status *types.Status) error {
	latestReleaseURL := strings.Replace(githubEvent.Repository.ReleaseUrl, "{/id}", "/latest", 1)

	status.LogF("fetching latest release for '%s' from '%s'", githubEvent.Repository.Name, latestReleaseURL)
	release, err := fetchRelease(latestReleaseURL)
	githubEvent.Release = release

	return err
}

func fetchRelease(releaseURL string) (Release, error) {
	var release Release
	var err error
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	if resp, err := client.Get(releaseURL); err == nil {
		err = parsePayload(resp.Body, &release)
	}

	return release, err
}

func parsePayload(payload io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(payload)
	return decoder.Decode(v)
}
