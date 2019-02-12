package main

import (
	"encoding/json"
	"ghwebhooks/types"
	"os"
)

func LoadConfig() types.Config {
	var config types.Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()

	if err != nil {
		panic("failed to load config.json")
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
