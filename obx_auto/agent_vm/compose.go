package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ComposeExecutor struct {
	workDir string
}

func NewComposeExecutor() *ComposeExecutor {
	return &ComposeExecutor{}
}

func (e *ComposeExecutor) Deploy(composeFile, projectName string, envVars []string) (string, error) {
	if err := os.WriteFile(filepath.Join(e.workDir, "docker-compose.yml"), []byte(composeFile), 0644); err != nil {
		return "", fmt.Errorf("write compose file: %w", err)
	}
	if err := e.writeEnvFile(envVars); err != nil {
		return "", err
	}

	pullOut, err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "-p", projectName, "pull").CombinedOutput()
	if err != nil {
		return string(pullOut), fmt.Errorf("docker compose pull: %w", err)
	}

	upOut, err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "-p", projectName, "up", "-d").CombinedOutput()
	if err != nil {
		return string(upOut), fmt.Errorf("docker compose up: %w", err)
	}

	return string(pullOut) + "\n" + string(upOut), nil
}

func (e *ComposeExecutor) Down(composeFile, projectName string, envVars []string) (string, error) {
	if err := os.WriteFile(filepath.Join(e.workDir, "docker-compose.yml"), []byte(composeFile), 0644); err != nil {
		return "", fmt.Errorf("write compose file: %w", err)
	}
	if err := e.writeEnvFile(envVars); err != nil {
		return "", err
	}
	out, err := exec.Command("docker", "compose", "-f", filepath.Join(e.workDir, "docker-compose.yml"), "-p", projectName, "down").CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker compose down: %w", err)
	}
	return string(out), nil
}

func (e *ComposeExecutor) writeEnvFile(envVars []string) error {
	if len(envVars) == 0 {
		return nil
	}
	if err := os.WriteFile(filepath.Join(e.workDir, ".env"), []byte(strings.Join(envVars, "\n")), 0644); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}
	return nil
}

func (e *ComposeExecutor) Pull(composeFile, projectName string) (string, error) {
	out, err := exec.Command("docker", "compose", "-f", filepath.Join(e.workDir, "docker-compose.yml"), "-p", projectName, "pull").CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker compose pull: %w", err)
	}
	return string(out), nil
}

func (e *ComposeExecutor) Logs(composeFile, projectName, service string, tail int) (string, error) {
	args := []string{"compose", "-f", filepath.Join(e.workDir, "docker-compose.yml"), "-p", projectName, "logs"}
	if service != "" {
		args = append(args, service)
	}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker compose logs: %w", err)
	}
	return string(out), nil
}

func (e *ComposeExecutor) SetWorkDir(dir string) {
	e.workDir = dir
}

func WriteComposeFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir compose dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write compose file: %w", err)
	}
	return nil
}

func ParseComposeProjectName(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "COMPOSE_PROJECT_NAME") || strings.Contains(line, "PREFIX") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
