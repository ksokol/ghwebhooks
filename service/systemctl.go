package service

import (
	"os/exec"
)

func Start(service string) (bool, error) {
	return command("restart", service)
}

func Stop(service string) (bool, error) {
	return command("stop", service)
}

func command(command string,  service string) (bool, error) {
	if err := systemctl(command, service); err != nil {
		return false, err
	}
	return true, nil
}

func systemctl(arg ...string) error {
	return exec.Command("systemctl", arg...).Run()
}
