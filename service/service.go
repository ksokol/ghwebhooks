package service

import (
	"fmt"
	"ghwebhooks/context"
	"ghwebhooks/github"
	"ghwebhooks/types"
	"os"
)

func Update(context *context.Context, status *types.Status) {
	user := Lookup(context, status)

	if status.Success != true {
		return
	}

	if downloadArtefact(context, status); status.Success != true {
		return
	}
	if stopService(context, status); status.Success != true {
		return
	}

	if replaceArtefact(user, context, status); status.Success != true {
		return
	}

	startService(context, status)
}

func stopService(context *context.Context, status *types.Status) {
	status.Log("stopping service")
	if ok, err := Stop(context.AppName); err != nil {
		status.Fail(err)
	} else {
		status.LogF("service stopped: %v", ok)
	}
}

func downloadArtefact(context *context.Context, status *types.Status) {
	if len(context.ArtefactURL) == 0 {
		status.FailMessage("artefact url not present")
		return
	}

	tmpFile := tmpFile(context)

	out, err := os.Create(tmpFile)
	defer out.Close()

	if err != nil {
		status.Fail(err)
		return
	}

	url := context.ArtefactURL
	status.LogF("downloading artefact %s", url)

	if err := github.Download(url, out); err != nil {
		status.Fail(err)
	}
}

func replaceArtefact(user UserLookup, context *context.Context, status *types.Status) {
	oldpath := tmpFile(context)
	newpath := fmt.Sprintf("%s/%s.jar", context.AppDir, context.AppName)

	status.LogF("replacing %s with %s", oldpath, newpath)

	if err := os.Rename(oldpath, newpath); err != nil {
		status.Fail(err)
		return
	}

	if err := os.Chown(newpath, user.uid, user.gid); err != nil {
		status.Fail(err)
	}
}

func tmpFile(context *context.Context) string {
	return fmt.Sprintf("%s/%s.jar", context.AppDir, "tmp")
}

func startService(context *context.Context, status *types.Status) {
	status.Log("starting service")
	if ok, err := Start(context.AppName); err != nil {
		status.Fail(err)
	} else {
		status.LogF("service started: %v", ok)
	}
}
