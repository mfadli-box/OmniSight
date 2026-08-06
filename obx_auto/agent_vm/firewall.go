package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type FirewallManager struct{}

func NewFirewallManager() *FirewallManager {
	return &FirewallManager{}
}

func (f *FirewallManager) Deploy(rulesContent string) (string, error) {
	if runtime.GOOS == "windows" {
		return f.deployWindows(rulesContent)
	}
	return f.deployLinux(rulesContent)
}

func (f *FirewallManager) deployLinux(rulesContent string) (string, error) {
	tmpPath := "/tmp/nftables_ruleset.conf"
	if err := os.WriteFile(tmpPath, []byte(rulesContent), 0644); err != nil {
		return "", fmt.Errorf("write nftables config: %w", err)
	}

	out, err := exec.Command("nft", "-f", tmpPath).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("nft apply: %w", err)
	}

	os.Remove(tmpPath)
	return string(out), nil
}

func (f *FirewallManager) deployWindows(_ string) (string, error) {
	return "", fmt.Errorf("windows firewall deploy not yet implemented")
}

func (f *FirewallManager) Test(rulesContent string) (bool, string, error) {
	if runtime.GOOS == "windows" {
		return f.testWindows(rulesContent)
	}
	return f.testLinux(rulesContent)
}

func (f *FirewallManager) testLinux(rulesContent string) (bool, string, error) {
	tmpPath := "/tmp/nftables_test.conf"
	if err := os.WriteFile(tmpPath, []byte(rulesContent), 0644); err != nil {
		return false, "", fmt.Errorf("write test config: %w", err)
	}
	defer os.Remove(tmpPath)

	out, err := exec.Command("nft", "-c", "-f", tmpPath).CombinedOutput()
	if err != nil {
		return false, string(out), nil
	}
	return true, string(out), nil
}

func (f *FirewallManager) testWindows(_ string) (bool, string, error) {
	return false, "", fmt.Errorf("windows firewall test not yet implemented")
}

func (f *FirewallManager) Backup() (string, error) {
	if runtime.GOOS == "windows" {
		return f.backupWindows()
	}
	return f.backupLinux()
}

func (f *FirewallManager) backupLinux() (string, error) {
	out, err := exec.Command("nft", "list", "ruleset").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("nft list: %w", err)
	}
	return string(out), nil
}

func (f *FirewallManager) backupWindows() (string, error) {
	return "", fmt.Errorf("windows firewall backup not yet implemented")
}

func (f *FirewallManager) Restore(rulesContent string) (string, error) {
	if runtime.GOOS == "windows" {
		return f.restoreWindows(rulesContent)
	}
	return f.restoreLinux(rulesContent)
}

func (f *FirewallManager) restoreLinux(rulesContent string) (string, error) {
	tmpPath := "/tmp/nftables_restore.conf"
	if err := os.WriteFile(tmpPath, []byte(rulesContent), 0644); err != nil {
		return "", fmt.Errorf("write restore config: %w", err)
	}
	defer os.Remove(tmpPath)

	out, err := exec.Command("nft", "-f", tmpPath).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("nft restore: %w", err)
	}
	return string(out), nil
}

func (f *FirewallManager) restoreWindows(_ string) (string, error) {
	return "", fmt.Errorf("windows firewall restore not yet implemented")
}
