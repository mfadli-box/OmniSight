package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type ContainerExecutor struct{}

func NewContainerExecutor() *ContainerExecutor {
	return &ContainerExecutor{}
}

type ContainerRunConfig struct {
	Name          string   `json:"name"`
	Image         string   `json:"image"`
	Tag           string   `json:"tag"`
	Ports         []string `json:"ports"`
	Volumes       []string `json:"volumes"`
	Env           []string `json:"env"`
	Network       string   `json:"network"`
	RestartPolicy string   `json:"restart_policy"`
	Command       string   `json:"command"`
}

func (e *ContainerExecutor) Run(cfg ContainerRunConfig) (string, error) {
	args := []string{"run", "-d"}
	if cfg.Name != "" {
		args = append(args, "--name", cfg.Name)
	}
	for _, p := range cfg.Ports {
		if strings.TrimSpace(p) != "" {
			args = append(args, "-p", p)
		}
	}
	for _, v := range cfg.Volumes {
		if strings.TrimSpace(v) != "" {
			args = append(args, "-v", v)
		}
	}
	for _, env := range cfg.Env {
		if strings.TrimSpace(env) != "" {
			args = append(args, "-e", env)
		}
	}
	if cfg.Network != "" {
		args = append(args, "--network", cfg.Network)
	}
	restart := cfg.RestartPolicy
	if restart == "" {
		restart = "unless-stopped"
	}
	args = append(args, "--restart", restart)
	image := cfg.Image
	if cfg.Tag != "" {
		image += ":" + cfg.Tag
	}
	args = append(args, image)
	if strings.TrimSpace(cfg.Command) != "" {
		args = append(args, strings.Fields(cfg.Command)...)
	}
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker run: %w", err)
	}
	return string(out), nil
}

func (e *ContainerExecutor) Start(container string) (string, error) {
	out, err := exec.Command("docker", "start", container).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker start: %w", err)
	}
	return string(out), nil
}

func (e *ContainerExecutor) Stop(container string) (string, error) {
	out, err := exec.Command("docker", "stop", container).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker stop: %w", err)
	}
	return string(out), nil
}

func (e *ContainerExecutor) Restart(container string) (string, error) {
	out, err := exec.Command("docker", "restart", container).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker restart: %w", err)
	}
	return string(out), nil
}

func (e *ContainerExecutor) Remove(container string, force bool) (string, error) {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, container)
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker rm: %w", err)
	}
	return string(out), nil
}

func (e *ContainerExecutor) Exec(container, command string) (string, int, error) {
	args := strings.Fields(command)
	cmdArgs := append([]string{"exec", container}, args...)
	out, err := exec.Command("docker", cmdArgs...).CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return string(out), exitCode, nil
}

func (e *ContainerExecutor) Logs(container string, tail int) (string, error) {
	args := []string{"logs"}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}
	args = append(args, container)
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker logs: %w", err)
	}
	return string(out), nil
}

func (e *ContainerExecutor) Inspect(container string) (string, error) {
	out, err := exec.Command("docker", "inspect", container).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker inspect: %w", err)
	}
	return string(out), nil
}

func (e *ContainerExecutor) Stats(container string) (string, error) {
	args := []string{"stats", "--no-stream", "--format", "{{json .}}"}
	if container != "" {
		args = append(args, container)
	}
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker stats: %w", err)
	}
	return string(out), nil
}

func (e *ContainerExecutor) List() ([]ContainerInfo, error) {
	out, err := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Status}}|{{.State}}").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker ps: %w", err)
	}

	var containers []ContainerInfo
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 5 {
			containers = append(containers, ContainerInfo{
				ID:     parts[0],
				Name:   parts[1],
				Image:  parts[2],
				Status: parts[3],
				State:  parts[4],
			})
		}
	}
	return containers, nil
}

type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	Status string
	State  string
}
