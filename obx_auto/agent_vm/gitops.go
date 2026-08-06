package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitOpsManager struct {
	workDir string
}

func NewGitOpsManager(workDir string) *GitOpsManager {
	return &GitOpsManager{workDir: workDir}
}

func (g *GitOpsManager) Sync(repoURL, branch, deployCmd string) (string, error) {
	if err := os.MkdirAll(g.workDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir workdir: %w", err)
	}

	gitDir := filepath.Join(g.workDir, ".git")
	var out string

	if _, err := os.Stat(gitDir); err == nil {
		pullOut, pullErr := g.gitPull(branch)
		out = pullOut
		if pullErr != nil {
			return out, pullErr
		}
	} else {
		cloneOut, cloneErr := g.gitClone(repoURL, branch)
		out = cloneOut
		if cloneErr != nil {
			return out, cloneErr
		}
	}

	if deployCmd != "" {
		deployOut, deployErr := g.runDeploy(deployCmd)
		out += "\n" + deployOut
		if deployErr != nil {
			return out, deployErr
		}
	}

	return out, nil
}

func (g *GitOpsManager) gitClone(repoURL, branch string) (string, error) {
	args := []string{"clone"}
	if branch != "" && branch != "main" && branch != "master" {
		args = append(args, "-b", branch)
	}
	args = append(args, repoURL, g.workDir)

	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("git clone: %w", err)
	}
	return string(out), nil
}

func (g *GitOpsManager) gitPull(branch string) (string, error) {
	out, err := exec.Command("git", "-C", g.workDir, "checkout", branch).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("git checkout: %w", err)
	}

	pullOut, pullErr := exec.Command("git", "-C", g.workDir, "pull", "origin", branch).CombinedOutput()
	if pullErr != nil {
		return string(pullOut), fmt.Errorf("git pull: %w", pullErr)
	}

	return string(out) + "\n" + string(pullOut), nil
}

func (g *GitOpsManager) runDeploy(cmd string) (string, error) {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return "", nil
	}

	program := fields[0]
	args := fields[1:]

	execCmd := exec.Command(program, args...)
	execCmd.Dir = g.workDir
	out, err := execCmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("deploy cmd: %w", err)
	}
	return string(out), nil
}

func (g *GitOpsManager) GetStatus() (string, error) {
	out, err := exec.Command("git", "-C", g.workDir, "status").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git status: %w", err)
	}
	return string(out), nil
}

func (g *GitOpsManager) GetCurrentCommit() (string, error) {
	out, err := exec.Command("git", "-C", g.workDir, "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
