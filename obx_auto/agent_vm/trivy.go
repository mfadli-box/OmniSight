package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TrivyScanner struct {
	cacheDir string
}

func NewTrivyScanner(cacheDir string) *TrivyScanner {
	if cacheDir == "" {
		cacheDir = "/tmp/trivy"
	}
	return &TrivyScanner{cacheDir: cacheDir}
}

func (t *TrivyScanner) ScanImage(image string) (*VulnScanResult, error) {
	tmpDir := filepath.Join(os.TempDir(), "trivy_output")
	os.MkdirAll(tmpDir, 0755)
	outputFile := filepath.Join(tmpDir, "trivy_result.json")

	_, err := exec.Command("trivy", "image",
		"--cache-dir", t.cacheDir,
		"--format", "json",
		"--output", outputFile,
		"--severity", "CRITICAL,HIGH,MEDIUM,LOW",
		image,
	).CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("trivy scan: %w", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("read trivy output: %w", err)
	}

	return t.parseResult(image, string(data))
}

func (t *TrivyScanner) ScanFilesystem(path string) (*VulnScanResult, error) {
	tmpDir := filepath.Join(os.TempDir(), "trivy_output")
	os.MkdirAll(tmpDir, 0755)
	outputFile := filepath.Join(tmpDir, "trivy_fs_result.json")

	_, err := exec.Command("trivy", "fs",
		"--cache-dir", t.cacheDir,
		"--format", "json",
		"--output", outputFile,
		"--severity", "CRITICAL,HIGH,MEDIUM,LOW",
		path,
	).CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("trivy fs scan: %w", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("read trivy output: %w", err)
	}

	return t.parseResult(path, string(data))
}

func (t *TrivyScanner) parseResult(target string, jsonData string) (*VulnScanResult, error) {
	result := &VulnScanResult{
		Target: target,
		Vulns:  []Vuln{},
	}

	var trivyResults []TrivyResults
	if err := json.Unmarshal([]byte(jsonData), &trivyResults); err != nil {
		return nil, fmt.Errorf("parse trivy json: %w", err)
	}

	for _, r := range trivyResults {
		if r.Vulnerabilities != nil {
			for _, v := range *r.Vulnerabilities {
				vuln := Vuln{
					ID:           v.VulnerabilityID,
					PkgName:      v.PkgName,
					InstalledVer: v.InstalledVersion,
					FixedVer:     v.FixedVersion,
					Severity:     v.Severity,
					Title:        v.Title,
					Description:  v.Description,
				}
				result.Vulns = append(result.Vulns, vuln)

				switch strings.ToLower(v.Severity) {
				case "critical":
					result.Critical++
				case "high":
					result.High++
				case "medium":
					result.Medium++
				case "low":
					result.Low++
				}
			}
		}
	}

	return result, nil
}

type VulnScanResult struct {
	Target   string `json:"target"`
	Critical int    `json:"critical"`
	High     int    `json:"high"`
	Medium   int    `json:"medium"`
	Low      int    `json:"low"`
	Vulns    []Vuln `json:"vulns"`
}

type Vuln struct {
	ID           string `json:"id"`
	PkgName      string `json:"pkg_name"`
	InstalledVer string `json:"installed_version"`
	FixedVer     string `json:"fixed_version"`
	Severity     string `json:"severity"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type TrivyResults struct {
	Target          string                `json:"Target"`
	Vulnerabilities *[]TrivyVulnerability `json:"Vulnerabilities"`
}

type TrivyVulnerability struct {
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	InstalledVersion string `json:"InstalledVersion"`
	FixedVersion     string `json:"FixedVersion"`
	Severity         string `json:"Severity"`
	Title            string `json:"Title"`
	Description      string `json:"Description"`
}
