package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FIMMonitor struct {
	baselineDir string
	monitorDir  string
	paths       []string
}

func NewFIMMonitor(baselineDir, monitorDir string, paths []string) *FIMMonitor {
	if baselineDir == "" {
		baselineDir = "/var/lib/agent_vm/fim"
	}
	return &FIMMonitor{
		baselineDir: baselineDir,
		monitorDir:  monitorDir,
		paths:       paths,
	}
}

func (f *FIMMonitor) InitBaseline() error {
	if err := os.MkdirAll(f.baselineDir, 0700); err != nil {
		return fmt.Errorf("mkdir baseline dir: %w", err)
	}

	for _, path := range f.paths {
		hash, err := f.computeFileHash(path)
		if err != nil {
			continue
		}
		baselinePath := f.getBaselinePath(path)
		if err := os.MkdirAll(filepath.Dir(baselinePath), 0700); err != nil {
			return fmt.Errorf("mkdir baseline subdir: %w", err)
		}
		if err := os.WriteFile(baselinePath, []byte(hash), 0600); err != nil {
			return fmt.Errorf("write baseline file: %w", err)
		}
	}
	return nil
}

func (f *FIMMonitor) Check() ([]FIMChange, error) {
	var changes []FIMChange

	for _, path := range f.paths {
		change, err := f.checkPath(path)
		if err != nil {
			continue
		}
		if change != nil {
			changes = append(changes, *change)
		}
	}

	return changes, nil
}

func (f *FIMMonitor) checkPath(path string) (*FIMChange, error) {
	baselinePath := f.getBaselinePath(path)
	baselineHash, err := os.ReadFile(baselinePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &FIMChange{
				Path:       path,
				ChangeType: "added",
				NewHash:    "",
			}, nil
		}
		return nil, err
	}

	currentHash, err := f.computeFileHash(path)
	if err != nil {
		return nil, err
	}

	if string(baselineHash) != currentHash {
		return &FIMChange{
			Path:       path,
			OldHash:    string(baselineHash),
			NewHash:    currentHash,
			ChangeType: "modified",
		}, nil
	}

	return nil, nil
}

func (f *FIMMonitor) computeFileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func (f *FIMMonitor) getBaselinePath(path string) string {
	safeName := strings.ReplaceAll(path, "/", "_")
	return filepath.Join(f.baselineDir, safeName+".sha256")
}

func (f *FIMMonitor) UpdateBaseline() error {
	for _, path := range f.paths {
		hash, err := f.computeFileHash(path)
		if err != nil {
			continue
		}
		baselinePath := f.getBaselinePath(path)
		os.WriteFile(baselinePath, []byte(hash), 0600)
	}
	return nil
}

type FIMChange struct {
	Path       string `json:"path"`
	OldHash    string `json:"old_hash"`
	NewHash    string `json:"new_hash"`
	ChangeType string `json:"change_type"`
}

type FIMReport struct {
	HostID    string      `json:"host_id"`
	RefID     string      `json:"ref_id"`
	Changes   []FIMChange `json:"changes"`
	ScannedAt string      `json:"scanned_at"`
	NoChanges bool        `json:"no_changes"`
}

func (f *FIMMonitor) GenerateReport(hostID, refID string) *FIMReport {
	changes, _ := f.Check()
	report := &FIMReport{
		HostID:    hostID,
		RefID:     refID,
		Changes:   changes,
		ScannedAt: time.Now().UTC().Format(time.RFC3339),
		NoChanges: len(changes) == 0,
	}
	return report
}
