package service

import "os/exec"

func Start(service string) error {
	if installed, err := isInstalled(service); err != nil || !installed {
		return err
	}
	return systemctl("start", service)
}

func isInstalled(service string) (bool, error) {
	out, err := exec.Command("systemctl", "is-active", service).Output()

	if err != nil {
		return false, err
	}

	state := string(out[:])

	return state == "active" || state == "inactive", err
}

func systemctl(arg ...string) error {
	return exec.Command("systemctl", arg...).Run()
}
