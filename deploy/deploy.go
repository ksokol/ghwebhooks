package deploy

import (
	"ghwebhooks/config"
	"os"
	"os/exec"
	"time"
)

type DeployLogMessage struct {
	Timestamp time.Time
	Message   string
}

type DeployLog struct {
	Success  bool
	Messages []DeployLogMessage
}

func Deploy(dir string, config *config.Config) DeployLog {
	os.Chdir(dir)
	out, err := exec.Command("python", "cron.py", config.Mail.From, config.Mail.To, config.Mail.Host).Output()

	if err != nil {
		return DeployLog{
			Success: false,
			Messages: []DeployLogMessage{DeployLogMessage{
				Timestamp: time.Now(),
				Message:   err.Error()}}}
	}

	return DeployLog{
		Success: true,
		Messages: []DeployLogMessage{{
			Timestamp: time.Now(),
			Message:   string(out[:])}}}
}
