package github

import (
	"bytes"
	"fmt"
	"ghwebhooks/config"
	"ghwebhooks/types"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Release struct {
	ID      int
	TagName string `json:"tag_name"`
	Draft   bool   `json:"draft"`
	Assets  []struct {
		Url string `json:"browser_download_url"`
	}
}

func (r *Release) isValid() bool {
	return r.ID != 0 && !r.Draft
}

type Repository struct {
	Name       string `json:"name"`
	ReleaseUrl string `json:"releases_url"`
	GitRefsUrl string `json:"git_refs_url"`
}

func (r *Repository) releasesUrl() string {
	return strings.Replace(r.ReleaseUrl, "{/id}", "", 1)
}

func (r *Repository) gitRefsUrl() string {
	return strings.Replace(r.GitRefsUrl, "{/sha}", "", 1)
}

func (r *Repository) releaseUrlFor(path interface{}) string {
	return fmt.Sprintf("%v/%v", r.releasesUrl(), path)
}

func (r *Repository) gitRefsUrlFor(ref interface{}) string {
	return fmt.Sprintf("%v/%v", r.gitRefsUrl(), ref)
}

type GithubEvent struct {
	Action     string
	Repository Repository
	Release    Release
}

func (g *GithubEvent) isValidRelease() bool {
	return g.Release.isValid() && g.Repository.Name != ""
}

func supportsApp(githubEvent *GithubEvent) bool {
	if _, err := config.GetAppConfig(githubEvent.Repository.Name); err == nil {
		return true
	}
	return false
}

func supportsAction(githubEvent *GithubEvent) bool {
	return githubEvent.Action == "published"
}

func Supports(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var githubEvent GithubEvent
		buf, _ := ioutil.ReadAll(req.Body)
		bodyClone := ioutil.NopCloser(bytes.NewBuffer(buf))
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		if err := parseBody(bodyClone, &githubEvent); err == nil {
			if supportsAction(&githubEvent) && supportsApp(&githubEvent) {
				handler(resp, req)
				return
			}
		}

		resp.WriteHeader(404)
	}
}

func Parse(body io.ReadCloser, status *types.Status) (GithubEvent, error) {
	var githubEvent GithubEvent

	if err := parseBody(body, &githubEvent); err != nil {
		return githubEvent, err
	}

	name := githubEvent.Repository.Name

	if !githubEvent.isValidRelease() {
		status.LogF("no valid release found for '%s'", name)

		if err := updateWithLatestRelease(&githubEvent, status); err != nil {
			return githubEvent, err
		}

		if !githubEvent.isValidRelease() {
			return githubEvent, fmt.Errorf("release not found for '%s'", name)
		}
	}

	return githubEvent, nil
}

func RemoveDraftReleases(githubEvent *GithubEvent, status *types.Status) {
	var releases []Release
	url := githubEvent.Repository.releasesUrl()

	status.LogF("fetching latest releases from '%s'", url)

	if err := Get(url, &releases); err != nil {
		status.Fail(err)
		return
	}

	for _, release := range releases {
		status.LogF("release: %s, draft: %t", release.TagName, release.Draft)

		if release.Draft {
			status.LogF("removing draft release '%s'", release.TagName)
			if err := Delete(githubEvent.Repository.releaseUrlFor(release.ID)); err != nil {
				status.LogF("draft release '%s' not removed due to '%s'", release.TagName, err.Error())
			} else {
				status.LogF("removed draft release '%s'", release.TagName)
			}
		}
	}
}

func RemovePreviousReleases(githubEvent *GithubEvent, status *types.Status) {
	var releases []Release
	url := githubEvent.Repository.releasesUrl()

	status.LogF("fetching latest releases from '%s'", url)

	if err := Get(url, &releases); err != nil {
		status.Fail(err)
		return
	}

	for _, release := range releases {
		if release.ID == githubEvent.Release.ID {
			continue
		}

		status.LogF("release: %s, draft: %t", release.TagName, release.Draft)

		if !release.Draft {
			status.LogF("removing release '%s'", release.TagName)
			if err := Delete(githubEvent.Repository.gitRefsUrlFor("tags/" + release.TagName)); err != nil {
				status.LogF("release '%s' not removed due to '%s'", release.TagName, err.Error())
			} else {
				status.LogF("removed release '%s'", release.TagName)
			}
		}
	}
}

func updateWithLatestRelease(githubEvent *GithubEvent, status *types.Status) error {
	latestReleaseURL := githubEvent.Repository.releaseUrlFor("latest")

	status.LogF("fetching latest release for '%s' from '%s'", githubEvent.Repository.Name, latestReleaseURL)
	release, err := fetchRelease(latestReleaseURL)
	githubEvent.Release = release

	return err
}

func fetchRelease(releaseURL string) (Release, error) {
	var release Release
	err := Get(releaseURL, &release)
	return release, err
}
