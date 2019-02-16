package types

import "ghwebhooks/config"

type Artefact struct {
	ReleaseID   int
	Name        string
	Tag         string
	ArtefactURL string
}

type Context struct {
	AppName string
	AppDir  string
	Artefact
}

func NewContext(artefact Artefact, config config.AppConfig) Context {
	return Context{
		config.Name,
		config.Dir,
		artefact,
	}
}
