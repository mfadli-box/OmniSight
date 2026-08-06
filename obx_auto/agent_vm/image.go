package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type ImageExecutor struct{}

func NewImageExecutor() *ImageExecutor {
	return &ImageExecutor{}
}

func (e *ImageExecutor) Pull(image string) (string, error) {
	out, err := exec.Command("docker", "pull", image).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker pull: %w", err)
	}
	return string(out), nil
}

func (e *ImageExecutor) Inspect(image string) (digest string, size int64, err error) {
	out, err := exec.Command("docker", "image", "inspect", image, "--format", "{{json .}}").CombinedOutput()
	if err != nil {
		return "", 0, fmt.Errorf("docker image inspect: %w", err)
	}
	var info struct {
		RepoDigests []string `json:"RepoDigests"`
		Size        int64    `json:"Size"`
	}
	if err := json.Unmarshal(out, &info); err != nil {
		return "", 0, fmt.Errorf("docker image inspect parse: %w", err)
	}
	if len(info.RepoDigests) > 0 {
		digest = info.RepoDigests[0]
	}
	return digest, info.Size, nil
}

func (e *ImageExecutor) List() ([]ImageInfo, error) {
	out, err := exec.Command("docker", "images", "--format", "{{.Repository}}|{{.Tag}}|{{.ID}}|{{.Size}}").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker images: %w", err)
	}

	var images []ImageInfo
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 4 {
			images = append(images, ImageInfo{
				Repository: parts[0],
				Tag:        parts[1],
				ID:         parts[2],
				Size:       parts[3],
			})
		}
	}
	return images, nil
}

func (e *ImageExecutor) Prune(all bool) (string, error) {
	args := []string{"image", "prune", "-f"}
	if all {
		args = append(args, "-a")
	}
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker image prune: %w", err)
	}
	return string(out), nil
}

func (e *ImageExecutor) Remove(image string, force bool) (string, error) {
	args := []string{"image", "rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, image)
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("docker image rm: %w", err)
	}
	return string(out), nil
}

type ImageInfo struct {
	Repository string
	Tag        string
	ID         string
	Size       string
}
