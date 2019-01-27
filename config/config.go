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

func LoadConfig() Config {
	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()

	if err != nil {
		panic("failed to load config.json")
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
