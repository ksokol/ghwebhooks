package service

import (
	"fmt"
	"ghwebhooks/context"
	"ghwebhooks/types"
	"os"
	"os/exec"
)

func Update(context *context.Context, status *types.Status) {
	user := Lookup(context, status)

	if status.Success != true {
		return
	}

	if stopService(context, status); status.Success != true {
		return
	}

	if download(context, status); status.Success != true {
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

func download(context *context.Context, status *types.Status) {
	os.Chdir(context.AppDir)
	out, err := exec.Command("python", "cron.py", context.ArtefactURL).Output()

	if err != nil {
		status.Fail(err)
		return
	}

	status.Log(string(out[:]))
}

func replaceArtefact(user UserLookup, context *context.Context, status *types.Status) {
	oldpath := fmt.Sprintf("%s/%s.jar", context.AppDir, "tmp")
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

func startService(context *context.Context, status *types.Status) {
	status.Log("starting service")
	if ok, err := Start(context.AppName); err != nil {
		status.Fail(err)
	} else {
		status.LogF("service started: %v", ok)
	}
}
