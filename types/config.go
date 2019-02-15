package types

import "fmt"

type Http struct {
	ListenAddress string
}

type Mail struct {
	Host string
	From string
	To   string
}

type App struct {
	Name string
	Dir  string
}

type Config struct {
	Dev    bool
	Secret string
	Http   Http
	Mail   Mail
	Apps   []App
}

type AppConfig struct {
	App
	Mail
}

func (c *Config) For(appName string) AppConfig {
	for _, app := range c.Apps {
		if app.Name == appName {
			return AppConfig{
				app,
				c.Mail,
			}
		}
	}

	panic(fmt.Sprintf("config for app '%s' not found", appName))
}
