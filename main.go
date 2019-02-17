package main

import (
	"encoding/json"
	"ghwebhooks/config"
	"ghwebhooks/context"
	"ghwebhooks/deploy"
	"ghwebhooks/github"
	"ghwebhooks/security"
	"ghwebhooks/types"
	"log"
	"net/http"
	"sync"
)

func main() {
	var activeDeployments sync.Map

	http.HandleFunc("/", security.Secured(github.SupportsApp(func(resp http.ResponseWriter, req *http.Request) {
		status := types.NewStatus()
		statusCode := 200

		if context, err := context.NewContext(req.Body, &status); err != nil {
			status.Fail(err)
			statusCode = 400
		} else {
			if _, loaded := activeDeployments.LoadOrStore(context.AppName, nil); loaded == true {
				statusCode = 409
			} else {
				deploy.Deploy(&context, &status)
				activeDeployments.Delete(context.AppName)
			}
		}

		writeRespone(resp, &status, statusCode)
	})))

	log.Fatal(http.ListenAndServe(config.GetHttpListenerAddress(), nil))
}

func writeRespone(resp http.ResponseWriter, status *types.Status, statusCode int) {
	json, err := json.MarshalIndent(status, "", "  ")

	if err != nil || !status.Success {
		statusCode = 500
	}

	resp.WriteHeader(statusCode)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(json)
}
