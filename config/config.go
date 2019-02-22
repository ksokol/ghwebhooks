package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Dev    bool
	Github struct {
		Secret      string
		AccessToken string
	}
	Http struct {
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

func GetSecret() string { return config.Github.Secret }

func GetAccessToken() string { return config.Github.AccessToken }

func GetMailTo() string { return config.Mail.To }

func GetMailFrom() string { return config.Mail.From }

func GetMailHost() string { return config.Mail.Host }

func GetAppConfig(appName string) (appConfig AppConfig, err error) {
	for _, app := range config.Apps {
		if app.Name == appName {
			return AppConfig{
				app.Name,
				app.Dir,
			}, err
		}
	}
	return appConfig, fmt.Errorf("app config for '%s' not found", appName)
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
