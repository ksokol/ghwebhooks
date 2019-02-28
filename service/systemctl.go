package service

import (
	"os/exec"
	"strings"
)

const (
	active   = "active"
	inactive = "inactive"
	unknown  = "unknown"
)

func Start(service string) (bool, error) { return command("start", inactive, service) }

func Stop(service string) (bool, error) { return command("stop", active, service) }

func command(command string, expectedState string, service string) (bool, error) {
	if currentState, err := getServiceState(service); err != nil || currentState != expectedState {
		return false, err
	}
	return true, systemctl(command, service)
}

func getServiceState(service string) (string, error) {
	out, err := exec.Command("systemctl", "is-active", service).Output()
	state := strings.TrimSuffix(string(out[:]), "\n")

	/*
	 systemctl returns code 3 in case of an inactive or unknown service.
	 Therefore we need to check stdout for the unit state before.
	*/
	if state == active || state == inactive || state == unknown {
		return state, nil
	}

	return "", err
}

func systemctl(arg ...string) error {
	return exec.Command("systemctl", arg...).Run()
}
