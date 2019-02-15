package types

type Artefact struct {
	ReleaseID   int
	Name        string
	Tag         string
	ArtefactURL string
}

type Context struct {
	AppName string
	AppDir  string
	Mail
	Artefact
}

func NewContext(artefact Artefact, config AppConfig) Context {
	return Context{
		config.Name,
		config.Dir,
		config.Mail,
		artefact,
	}
}
