package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Config struct {
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

func loadConfig() (Config, error) {
	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()

	if err == nil {
		jsonParser := json.NewDecoder(configFile)
		jsonParser.Decode(&config)
	}

	return config, err
}

func deploy(dir string, config *Config) {
	os.Chdir(dir)
	cmd := exec.Command("python", "cron.py", config.Mail.From, config.Mail.To, config.Mail.Host)
	cmd.Run()
}

func main() {
	config, err := loadConfig()

	if err != nil {
		panic("failed to load config.json")
	}

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		for _, app := range config.Apps {
			if app.Name == req.URL.Path[1:] {
				deploy(app.Dir, &config)
				return
			}
		}

		resp.WriteHeader(404)
	})

	log.Fatal(http.ListenAndServe(config.Http.ListenAddress, nil))
}
