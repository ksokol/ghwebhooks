package main

import (
	"ghwebhooks/config"
	"ghwebhooks/github"
	"ghwebhooks/handler"
	"ghwebhooks/security"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", security.Secured(github.SupportsApp(handler.DeployHandler)))

	log.Fatal(http.ListenAndServe(config.GetHttpListenerAddress(), nil))
}
