package handler

import (
	"encoding/json"
	"ghwebhooks/context"
	"ghwebhooks/deploy"
	"ghwebhooks/types"
	"net/http"
	"sync"
)

var activeDeployments sync.Map

func writeRespone(resp http.ResponseWriter, status *types.Status, statusCode int) {
	json, err := json.MarshalIndent(status, "", "  ")

	if err != nil || !status.Success {
		statusCode = 500
	}

	resp.WriteHeader(statusCode)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(json)
}

func DeployHandler(resp http.ResponseWriter, req *http.Request) {
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
}
