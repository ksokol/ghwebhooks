package deploy

import (
	"ghwebhooks/deploy/mail"
	"ghwebhooks/types"
	"os"
	"os/exec"
)

func Deploy(context *types.Context, status *types.Status) {
	os.Chdir(context.AppDir)
	out, err := exec.Command("python", "cron.py").Output()

	if err != nil {
		status.Fail(err)
		return
	}

	status.Log(string(out[:]))
	status.Log("sending email")

	if err := mail.Sendmail(context); err != nil {
		status.Fail(err)
	} else {
		status.Log("email send")
	}
}
