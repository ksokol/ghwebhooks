package main

import (
	"encoding/json"
	"ghwebhooks/deploy"
	"ghwebhooks/security"
	"ghwebhooks/types"
	"log"
	"net/http"
	"sync"
)

func main() {
	config := LoadConfig()

	var activeDeployments sync.Map

	http.HandleFunc("/", security.Secured(func(resp http.ResponseWriter, req *http.Request) {
		for _, app := range config.Apps {
			if app.Name == req.URL.Path[1:] {
				context := types.NewContext(app.Name, app.Dir, &config)
				status := types.NewStatus()

				_, loaded := activeDeployments.LoadOrStore(app.Name, nil)

				if loaded == true {
					resp.WriteHeader(409)
					return
				}

				deploy.Deploy(&context, &status)
				json, err := json.MarshalIndent(status, "", "  ")

				if err != nil {
					resp.WriteHeader(500)
					return
				}

				var statusCode int
				if statusCode = 200; !status.Success {
					statusCode = 500
				}

				resp.WriteHeader(statusCode)
				resp.Header().Set("Content-Type", "application/json")
				resp.Write(json)

				activeDeployments.Delete(context.AppName)

				return
			}
		}

		resp.WriteHeader(404)
	}, &config))

	log.Fatal(http.ListenAndServe(config.Http.ListenAddress, nil))
}
