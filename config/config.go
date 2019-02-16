package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Dev    bool
	Secret string
	Http   struct {
		ListenAddress string
	}
	Mail struct {
		Host string
		From string
		To   string
	}
	Apps []struct {
		Name string
		Dir  string
	}
}

type AppConfig struct {
	Name string
	Dir  string
}

var config Config

func GetHttpListenerAddress() string { return config.Http.ListenAddress }

func IsDevEnv() bool { return config.Dev }

func GetSecret() string { return config.Secret }

func GetMailTo() string { return config.Mail.To }

func GetMailFrom() string { return config.Mail.From }

func GetMailHost() string { return config.Mail.Host }

func GetAppConfig(appName string) (appConfig AppConfig, ok bool) {
	for _, app := range config.Apps {
		if app.Name == appName {
			return AppConfig{
				app.Name,
				app.Dir,
			}, true
		}
	}
	return appConfig, false
}

func init() {
	configFile, err := os.Open("config.json")
	defer configFile.Close()

	if err != nil {
		panic("failed to load config.json")
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}
