package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type LegoSSL struct {
	legoDir  string
	email    string
	domains  []string
	provider string
}

func NewLegoSSL(legoDir, email string, domains []string) *LegoSSL {
	return &LegoSSL{
		legoDir:  legoDir,
		email:    email,
		domains:  domains,
		provider: "cloudflare",
	}
}

func (l *LegoSSL) SetProvider(provider string) {
	l.provider = provider
}

func (l *LegoSSL) Provision() (string, string, error) {
	if err := os.MkdirAll(l.legoDir, 0700); err != nil {
		return "", "", fmt.Errorf("mkdir lego dir: %w", err)
	}

	args := []string{
		"--path", l.legoDir,
		"--email", l.email,
		"--domains", l.domains[0],
		"--dns", l.provider,
		"--accept-tos",
		"run",
	}

	if len(l.domains) > 1 {
		for _, d := range l.domains[1:] {
			args = append(args, "--domains", d)
		}
	}

	cmd := exec.Command("lego", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", string(out), fmt.Errorf("lego run: %w", err)
	}

	certPath := filepath.Join(l.legoDir, "certificates", l.domains[0]+".crt")
	keyPath := filepath.Join(l.legoDir, "certificates", l.domains[0]+".key")

	if _, err := os.Stat(certPath); err != nil {
		return "", string(out), fmt.Errorf("certificate not found: %w", err)
	}

	return certPath, keyPath, nil
}

func (l *LegoSSL) Renew(domain string) (string, error) {
	args := []string{
		"--path", l.legoDir,
		"--domains", domain,
		"--dns", l.provider,
		"--accept-tos",
		"renew",
	}

	out, err := exec.Command("lego", args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("lego renew: %w", err)
	}
	return string(out), nil
}

func (l *LegoSSL) List() ([]CertInfo, error) {
	certDir := filepath.Join(l.legoDir, "certificates")
	entries, err := os.ReadDir(certDir)
	if err != nil {
		return nil, fmt.Errorf("read cert dir: %w", err)
	}

	var certs []CertInfo
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".crt" {
			continue
		}
		certs = append(certs, CertInfo{
			Domain: e.Name(),
			Path:   filepath.Join(certDir, e.Name()),
		})
	}
	return certs, nil
}

type CertInfo struct {
	Domain string
	Path   string
}
