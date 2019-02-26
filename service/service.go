package service

import (
	"ghwebhooks/context"
	"ghwebhooks/types"
	"os"
	"os/exec"
)

func Update(context *context.Context, status *types.Status) {
	os.Chdir(context.AppDir)
	out, err := exec.Command("python", "cron.py", context.ArtefactURL).Output()

	if err != nil {
		status.Fail(err)
		return
	}

	status.Log(string(out[:]))
}
