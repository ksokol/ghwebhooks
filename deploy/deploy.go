package deploy

import (
	"ghwebhooks/config"
	"os"
	"os/exec"
)

func Deploy(dir string, config *config.Config) {
	os.Chdir(dir)
	cmd := exec.Command("python", "cron.py", config.Mail.From, config.Mail.To, config.Mail.Host)
	cmd.Run()
}
