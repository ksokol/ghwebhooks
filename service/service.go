package service

import (
	"fmt"
	"ghwebhooks/context"
	"ghwebhooks/types"
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

type userLookup struct {
	uid, gid int
}

func Update(context *context.Context, status *types.Status) {
	user := lookup(context, status)

	if status.Success != true {
		return
	}

	if download(context, status); status.Success != true {
		return
	}

	if replaceArtefact(user, context, status); status.Success != true {
		return
	}

	status.Log("starting service")
	if err := start(context.AppName); err != nil {
		status.Fail(err)
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

func replaceArtefact(user userLookup, context *context.Context, status *types.Status) {
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

func lookup(context *context.Context, status *types.Status) (lookup userLookup) {
	username := context.AppName
	if user, err := user.Lookup(username); err != nil {
		status.Fail(err)
	} else {
		uid, err := strconv.Atoi(user.Uid)
		gid, err := strconv.Atoi(user.Gid)

		if err != nil {
			status.Fail(err)
		} else {
			status.LogF("found uid: %v, gid: %v for user: %s", uid, gid, username)
			lookup = userLookup{
				uid,
				gid,
			}
		}
	}
	return lookup
}

func start(service string) error {
	return systemctl("start", service)
}

func systemctl(arg ...string) error {
	return exec.Command("systemctl", arg...).Run()
}
