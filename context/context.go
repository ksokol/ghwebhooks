package context

import (
	"ghwebhooks/config"
	"ghwebhooks/github"
	"ghwebhooks/types"
	"io"
)

type Artefact struct {
	ReleaseID   int
	Name        string
	Tag         string
	ArtefactURL string
}

type Context struct {
	AppName string
	AppDir  string
	Event   github.GithubEvent
	Artefact
}

func NewContext(body io.ReadCloser, status *types.Status) (context Context, err error) {
	if githubEvent, err := github.Parse(body, status); err != nil {
		return context, err
	} else {
		if appConfig, err := config.GetAppConfig(githubEvent.Repository.Name); err != nil {
			return context, err
		} else {
			return Context{
				appConfig.Name,
				appConfig.Dir,
				githubEvent,
				Artefact{
					githubEvent.Release.ID,
					githubEvent.Repository.Name,
					githubEvent.Release.TagName,
					githubEvent.Release.Assets[0].Url,
				},
			}, err
		}
	}
}
