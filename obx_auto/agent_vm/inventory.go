package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var inventoryRetry = RetryConfig{MaxRetries: 2, BaseDelay: 500 * time.Millisecond, MaxDelay: 3 * time.Second}

func execWithRetry(name string, args ...string) ([]byte, error) {
	var out []byte
	err := WithRetry(inventoryRetry, name, func() error {
		var err error
		out, err = exec.Command(name, args...).Output()
		return err
	})
	return out, err
}

type InventoryImage struct {
	Name   string `json:"name"`
	Tag    string `json:"tag"`
	Digest string `json:"digest"`
	Size   string `json:"size"`
}

type InventoryNetwork struct {
	Name       string `json:"name"`
	Driver     string `json:"driver"`
	Subnet     string `json:"subnet"`
	Gateway    string `json:"gateway"`
	Internal   bool   `json:"internal"`
	Attachable bool   `json:"attachable"`
}

type InventoryContainer struct {
	Name        string `json:"name"`
	ContainerID string `json:"container_id"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	State       string `json:"state"`
}

type InventoryCompose struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Config string `json:"config_file"`
}

type InventoryData struct {
	HostID     string               `json:"host_id"`
	Images     []InventoryImage     `json:"images"`
	Networks   []InventoryNetwork   `json:"networks"`
	Containers []InventoryContainer `json:"containers"`
	Compose    []InventoryCompose   `json:"compose"`
}

type InventoryCollector struct{}

func NewInventoryCollector() *InventoryCollector {
	return &InventoryCollector{}
}

func (c *InventoryCollector) Collect(hostID string) (*InventoryData, error) {
	data := &InventoryData{HostID: hostID}

	images, err := c.collectImages()
	if err != nil {
		return nil, fmt.Errorf("collect images: %w", err)
	}
	data.Images = images

	networks, err := c.collectNetworks()
	if err != nil {
		return nil, fmt.Errorf("collect networks: %w", err)
	}
	data.Networks = networks

	containers, err := c.collectContainers()
	if err != nil {
		return nil, fmt.Errorf("collect containers: %w", err)
	}
	data.Containers = containers

	compose, err := c.collectCompose()
	if err != nil {
		return nil, fmt.Errorf("collect compose: %w", err)
	}
	data.Compose = compose

	return data, nil
}

func (c *InventoryCollector) collectImages() ([]InventoryImage, error) {
	out, err := execWithRetry("docker", "images", "--format", "{{json .}}")
	if err != nil {
		return nil, err
	}

	var images []InventoryImage
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var raw struct {
			Repository string `json:"Repository"`
			Tag        string `json:"Tag"`
			Digest     string `json:"ID"`
			Size       string `json:"Size"`
			CreatedAt  string `json:"CreatedAt"`
		}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		if raw.Repository == "<none>" {
			continue
		}

		tag := raw.Tag
		if tag == "<none>" {
			tag = "latest"
		}

		images = append(images, InventoryImage{
			Name:   raw.Repository,
			Tag:    tag,
			Digest: raw.Digest,
			Size:   raw.Size,
		})
	}

	return images, nil
}

func (c *InventoryCollector) collectNetworks() ([]InventoryNetwork, error) {
	out, err := execWithRetry("docker", "network", "ls", "--format", "{{json .}}")
	if err != nil {
		return nil, err
	}

	type rawNetwork struct {
		Name   string `json:"Name"`
		Driver string `json:"Driver"`
		ID     string `json:"ID"`
	}

	var networks []InventoryNetwork
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var raw rawNetwork
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		if raw.Name == "bridge" || raw.Name == "host" || raw.Name == "none" {
			continue
		}

		netInfo, err := c.inspectNetwork(raw.ID)
		if err != nil {
			networks = append(networks, InventoryNetwork{
				Name:   raw.Name,
				Driver: raw.Driver,
			})
			continue
		}

		networks = append(networks, *netInfo)
	}

	return networks, nil
}

func (c *InventoryCollector) inspectNetwork(id string) (*InventoryNetwork, error) {
	out, err := execWithRetry("docker", "network", "inspect", id)
	if err != nil {
		return nil, err
	}

	var inspect []struct {
		Name       string `json:"Name"`
		Driver     string `json:"Driver"`
		Internal   bool   `json:"Internal"`
		Attachable bool   `json:"Attachable"`
		IPAM       struct {
			Config []struct {
				Subnet  string `json:"Subnet"`
				Gateway string `json:"Gateway"`
			} `json:"Config"`
		} `json:"IPAM"`
	}
	if err := json.Unmarshal(out, &inspect); err != nil {
		return nil, err
	}
	if len(inspect) == 0 {
		return nil, fmt.Errorf("no network found")
	}

	net := inspect[0]
	info := &InventoryNetwork{
		Name:       net.Name,
		Driver:     net.Driver,
		Internal:   net.Internal,
		Attachable: net.Attachable,
	}

	if len(net.IPAM.Config) > 0 {
		info.Subnet = net.IPAM.Config[0].Subnet
		info.Gateway = net.IPAM.Config[0].Gateway
	}

	return info, nil
}

func (c *InventoryCollector) collectContainers() ([]InventoryContainer, error) {
	out, err := execWithRetry("docker", "ps", "-a", "--format", "{{json .}}")
	if err != nil {
		return nil, err
	}

	var containers []InventoryContainer
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var raw struct {
			ID         string `json:"ID"`
			Name       string `json:"Names"`
			Image      string `json:"Image"`
			Status     string `json:"Status"`
			State      string `json:"State"`
			RunningFor string `json:"RunningFor"`
		}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		containers = append(containers, InventoryContainer{
			Name:        raw.Name,
			ContainerID: raw.ID,
			Image:       raw.Image,
			Status:      raw.Status,
			State:       raw.State,
		})
	}

	return containers, nil
}

func (c *InventoryCollector) collectCompose() ([]InventoryCompose, error) {
	out, err := execWithRetry("docker", "compose", "ls", "--format", "{{json .}}")
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "executable file not found") {
			return nil, nil
		}
		return nil, err
	}

	var composeList []InventoryCompose
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var raw struct {
			Name            string   `json:"Name"`
			Status          string   `json:"Status"`
			ComposeFiles    []string `json:"ConfigFiles"`
			RunningServices int      `json:"Running"`
			TotalServices   int      `json:"Total"`
		}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		configFile := ""
		if len(raw.ComposeFiles) > 0 {
			configFile = raw.ComposeFiles[0]
		}

		composeList = append(composeList, InventoryCompose{
			Name:   raw.Name,
			Status: raw.Status,
			Config: configFile,
		})
	}

	return composeList, nil
}
