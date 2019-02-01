package main

import (
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

				deploy.Deploy(app.Dir, &config)
				activeDeployments.Delete(app.Name)

				return
			}
		}

		resp.WriteHeader(404)
	}, &config))

	log.Fatal(http.ListenAndServe(config.Http.ListenAddress, nil))
}
