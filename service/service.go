package service

import (
	"ghwebhooks/context"
	"ghwebhooks/types"
	"os"
	"os/exec"
)

func systemctl(arg ...string) error {
	return exec.Command("systemctl", arg...).Run()
}

func start(service string) error {
	return systemctl("start", service)
}

func Update(context *context.Context, status *types.Status) {
	os.Chdir(context.AppDir)
	out, err := exec.Command("python", "cron.py", context.ArtefactURL).Output()

	if err != nil {
		status.Fail(err)
		return
	}

	status.Log(string(out[:]))

	status.Log("starting service")
	if err := start(context.AppName); err != nil {
		status.Fail(err)
	}
}
