package main

import (
	"ghwebhooks/config"
	"ghwebhooks/deploy"
	"log"
	"net/http"
)

func main() {
	config := config.LoadConfig()

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		for _, app := range config.Apps {
			if app.Name == req.URL.Path[1:] {
				deploy.Deploy(app.Dir, &config)
				return
			}
		}

		resp.WriteHeader(404)
	})

	log.Fatal(http.ListenAndServe(config.Http.ListenAddress, nil))
}
