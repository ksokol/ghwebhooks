package main

import (
	"encoding/json"
	"ghwebhooks/config"
	"ghwebhooks/deploy"
	"ghwebhooks/security"
	"log"
	"net/http"
	"sync"
)

func main() {
	config := config.LoadConfig()

	var activeDeployments sync.Map

	http.HandleFunc("/", security.Secured(func(resp http.ResponseWriter, req *http.Request) {
		for _, app := range config.Apps {
			if app.Name == req.URL.Path[1:] {
				_, loaded := activeDeployments.LoadOrStore(app.Name, nil)

				if loaded == true {
					resp.WriteHeader(409)
					return
				}

				deployLog := deploy.Deploy(app.Dir, &config)
				json, err := json.MarshalIndent(deployLog, "", "  ")

				if err != nil {
					resp.WriteHeader(500)
					return
				}

				resp.Header().Set("Content-Type", "application/json")
				resp.Write(json)

				activeDeployments.Delete(app.Name)

				return
			}
		}

		resp.WriteHeader(404)
	}, &config))

	log.Fatal(http.ListenAndServe(config.Http.ListenAddress, nil))
}
