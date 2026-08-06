package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type NginxManager struct {
	configDir string
	testBin   string
	reloadBin string
}

func NewNginxManager(configDir string) *NginxManager {
	return &NginxManager{
		configDir: configDir,
		testBin:   "/usr/sbin/nginx",
		reloadBin: "/usr/sbin/nginx",
	}
}

func (n *NginxManager) DeployConfig(siteName, content string) (string, error) {
	sitePath := filepath.Join(n.configDir, "sites-available", siteName+".conf")
	if err := os.MkdirAll(filepath.Dir(sitePath), 0755); err != nil {
		return "", fmt.Errorf("mkdir site dir: %w", err)
	}
	if err := os.WriteFile(sitePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write nginx config: %w", err)
	}
	return sitePath, nil
}

func (n *NginxManager) EnableSite(siteName string) error {
	available := filepath.Join(n.configDir, "sites-available", siteName+".conf")
	enabled := filepath.Join(n.configDir, "sites-enabled", siteName+".conf")
	linkExists, _ := pathExists(enabled)
	if !linkExists {
		if err := os.Symlink(available, enabled); err != nil {
			return fmt.Errorf("symlink enable site: %w", err)
		}
	}
	return nil
}

func (n *NginxManager) DisableSite(siteName string) error {
	enabled := filepath.Join(n.configDir, "sites-enabled", siteName+".conf")
	if err := os.Remove(enabled); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove site symlink: %w", err)
	}
	return nil
}

func (n *NginxManager) TestConfig(configPath string) (bool, string, error) {
	args := []string{"-t", "-c", configPath}
	out, err := exec.Command(n.testBin, args...).CombinedOutput()
	if err != nil {
		return false, string(out), fmt.Errorf("nginx test: %w", err)
	}
	return true, string(out), nil
}

func (n *NginxManager) TestConfigContent(content string) (bool, string, error) {
	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("nginx_test_%d.conf", os.Getpid()))
	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return false, "", fmt.Errorf("write temp config: %w", err)
	}
	defer os.Remove(tmpPath)

	return n.TestConfig(tmpPath)
}

func (n *NginxManager) Reload() (string, error) {
	out, err := exec.Command(n.reloadBin, "-s", "reload").CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("nginx reload: %w", err)
	}
	return string(out), nil
}

func (n *NginxManager) CheckRunning() bool {
	out, err := exec.Command(n.testBin, "-v").CombinedOutput()
	if err != nil {
		return false
	}
	_ = out
	return true
}

func (n *NginxManager) GetVersion() (string, error) {
	out, err := exec.Command(n.testBin, "-v").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("nginx version: %w", err)
	}
	return string(out), nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
