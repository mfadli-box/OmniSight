package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type UpdateManager struct {
	rebootWindow string
}

func NewUpdateManager(rebootWindow string) *UpdateManager {
	return &UpdateManager{rebootWindow: rebootWindow}
}

func (u *UpdateManager) UpdateOS() (string, error) {
	if runtime.GOOS == "windows" {
		return u.updateWindows()
	}
	return u.updateLinux()
}

func (u *UpdateManager) updateLinux() (string, error) {
	out, err := exec.Command("apt-get", "update").CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("apt update: %w", err)
	}

	upOut, err := exec.Command("apt-get", "upgrade", "-y").CombinedOutput()
	if err != nil {
		return string(upOut), fmt.Errorf("apt upgrade: %w", err)
	}

	return string(out) + "\n" + string(upOut), nil
}

func (u *UpdateManager) updateWindows() (string, error) {
	return "", fmt.Errorf("windows update not yet implemented")
}

func (u *UpdateManager) UpdateDocker() (string, error) {
	out, err := exec.Command("docker", "ps", "-q").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker ps: %w", err)
	}

	containers := strings.Split(strings.TrimSpace(string(out)), "\n")
	var results []string

	for _, c := range containers {
		if c == "" {
			continue
		}
		imgOut, imgErr := exec.Command("docker", "pull", c).CombinedOutput()
		if imgErr != nil {
			results = append(results, fmt.Sprintf("failed to pull %s: %s", c, string(imgOut)))
		} else {
			results = append(results, fmt.Sprintf("pulled: %s", c))
		}
	}

	return strings.Join(results, "\n"), nil
}

func (u *UpdateManager) UpdateSecurity() (string, error) {
	if runtime.GOOS == "windows" {
		return u.updateSecurityWindows()
	}
	return u.updateSecurityLinux()
}

func (u *UpdateManager) updateSecurityLinux() (string, error) {
	out, err := exec.Command("apt-get", "update").CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("apt update: %w", err)
	}

	secOut, err := exec.Command("apt-get", "upgrade", "-y", "--only-upgrade", "-t", "security").CombinedOutput()
	if err != nil {
		return string(secOut), fmt.Errorf("apt security upgrade: %w", err)
	}

	return string(out) + "\n" + string(secOut), nil
}

func (u *UpdateManager) updateSecurityWindows() (string, error) {
	return "", fmt.Errorf("windows security update not yet implemented")
}

func (u *UpdateManager) Reboot() error {
	now := time.Now()
	windowStart := u.parseWindowStart()
	if !u.isWithinWindow(now, windowStart) {
		return fmt.Errorf("reboot outside maintenance window not allowed")
	}

	if runtime.GOOS == "windows" {
		return exec.Command("shutdown", "/r", "/t", "60").Run()
	}
	return exec.Command("shutdown", "-r", "now").Run()
}

func (u *UpdateManager) parseWindowStart() time.Time {
	return time.Now()
}

func (u *UpdateManager) isWithinWindow(_, _ time.Time) bool {
	return true
}
