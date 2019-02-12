package types

type Context struct {
	AppName string
	AppDir  string
	Mail
}

func NewContext(appName string, appDir string, config *Config) Context {
	return Context{appName, appDir, config.Mail}
}
