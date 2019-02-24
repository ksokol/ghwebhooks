package deploy

import (
	"ghwebhooks/context"
	"ghwebhooks/deploy/mail"
	"ghwebhooks/github"
	"ghwebhooks/types"
	"os"
	"os/exec"
)

func Deploy(context *context.Context, status *types.Status) {
	os.Chdir(context.AppDir)
	out, err := exec.Command("python", "cron.py", context.ArtefactURL).Output()

	if err != nil {
		status.Fail(err)
		return
	}

	status.Log(string(out[:]))
	mail.Sendmail(context, status)

	if status.Success != true {
		return
	}

	github.RemovePreviousReleases(&context.Event, status)
	github.RemoveDraftReleases(&context.Event, status)
}
