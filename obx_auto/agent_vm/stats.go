package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type StatsCollector struct {
	dockerSocket string
	interval     time.Duration
	stopCh       chan struct{}
	wg           sync.WaitGroup
	onStats      func(HostStatsPayload)
}

func NewStatsCollector(dockerSocket string, interval time.Duration, onStats func(HostStatsPayload)) *StatsCollector {
	return &StatsCollector{
		dockerSocket: dockerSocket,
		interval:     interval,
		stopCh:       make(chan struct{}),
		onStats:      onStats,
	}
}

func (s *StatsCollector) Start() {
	s.wg.Add(1)
	go s.collectLoop()
}

func (s *StatsCollector) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

func (s *StatsCollector) collectLoop() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			stats, err := s.Collect()
			if err != nil {
				continue
			}
			if s.onStats != nil {
				s.onStats(*stats)
			}
		}
	}
}

func (s *StatsCollector) Collect() (*HostStatsPayload, error) {
	stats := &HostStatsPayload{}

	hostname, _ := os.Hostname()
	stats.Hostname = hostname
	stats.Uptime = getUptime()
	stats.CPUCount = runtime.NumCPU()
	stats.CPUPercent = getCPUUsage()

	memInfo := getMemoryInfo()
	stats.MemTotal = memInfo.Total
	stats.MemUsed = memInfo.Used

	diskInfo := getDiskInfo()
	stats.DiskTotal = diskInfo.Total
	stats.DiskUsed = diskInfo.Used

	containers, _ := s.getContainerStats()
	stats.Containers = containers

	stats.DockerRunning = s.checkDocker()
	if stats.DockerRunning {
		stats.DockerVersion = s.getDockerVersion()
	}

	stats.NginxRunning = checkNginx()

	return stats, nil
}

func (s *StatsCollector) getContainerStats() ([]ContainerStats, error) {
	out, err := exec.Command("docker", "stats", "--no-stream", "--format", "{{json .}}").CombinedOutput()
	if err != nil {
		return nil, err
	}

	var containers []ContainerStats
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var ids []string
	for _, line := range lines {
		if line == "" {
			continue
		}
		var c ContainerStats
		if err := json.Unmarshal([]byte(line), &c); err == nil {
			containers = append(containers, c)
			ids = append(ids, c.ID)
		}
	}

	if len(ids) > 0 {
		details := s.getContainerDetails(ids)
		for i := range containers {
			if d, ok := details[containers[i].ID]; ok {
				containers[i].Status = d.Status
				containers[i].RestartCount = d.RestartCount
				containers[i].UptimeSeconds = d.UptimeSeconds
			}
		}
	}
	return containers, nil
}

type containerDetail struct {
	Status        string
	RestartCount  int
	UptimeSeconds int64
}

func (s *StatsCollector) getContainerDetails(ids []string) map[string]containerDetail {
	details := map[string]containerDetail{}
	if len(ids) == 0 {
		return details
	}
	args := append([]string{"inspect", "--format", "{{.Id}}|{{.State.Status}}|{{.RestartCount}}|{{.State.StartedAt}}"}, ids...)
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return details
	}
	now := time.Now()
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}
		id := strings.TrimSpace(parts[0])
		status := strings.TrimSpace(parts[1])
		restartCount, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		startedAt := strings.TrimSpace(parts[3])
		var uptime int64
		if t, err := time.Parse(time.RFC3339, startedAt); err == nil {
			uptime = int64(now.Sub(t).Seconds())
			if uptime < 0 {
				uptime = 0
			}
		}
		details[id] = containerDetail{Status: status, RestartCount: restartCount, UptimeSeconds: uptime}
	}
	return details
}

func (s *StatsCollector) checkDocker() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func (s *StatsCollector) getDockerVersion() string {
	out, err := exec.Command("docker", "version", "--format", "{{.Server.Version}}").CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

type memInfo struct {
	Total int64
	Used  int64
}

func getMemoryInfo() memInfo {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return memInfo{}
	}
	var total, avail uint64
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			if f := strings.Fields(line); len(f) >= 2 {
				total, _ = strconv.ParseUint(f[1], 10, 64)
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			if f := strings.Fields(line); len(f) >= 2 {
				avail, _ = strconv.ParseUint(f[1], 10, 64)
			}
		}
	}
	if total == 0 {
		return memInfo{}
	}
	return memInfo{Total: int64(total * 1024), Used: int64((total - avail) * 1024)}
}

type diskInfo struct {
	Total int64
	Used  int64
}

func getDiskInfo() diskInfo {
	target := hostDiskRoot()
	out, err := exec.Command("df", "-P", target).CombinedOutput()
	if err != nil {
		return diskInfo{}
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 6 || fields[len(fields)-1] != target {
			continue
		}
		total, _ := strconv.ParseInt(fields[1], 10, 64)
		used, _ := strconv.ParseInt(fields[2], 10, 64)
		return diskInfo{Total: total * 1024, Used: used * 1024}
	}
	return diskInfo{}
}

func hostDiskRoot() string {
	if fi, err := os.Stat("/hostfs"); err == nil && fi.IsDir() {
		return "/hostfs"
	}
	return "/"
}

func getCPUUsage() float64 {
	firstIdle, firstTotal := readCPUStats()
	if firstTotal == 0 {
		return 0
	}
	time.Sleep(200 * time.Millisecond)
	secondIdle, secondTotal := readCPUStats()
	if secondTotal == 0 {
		return 0
	}
	dIdle := secondIdle - firstIdle
	dTotal := secondTotal - firstTotal
	if dTotal <= 0 {
		return 0
	}
	usage := 100.0 * (1 - float64(dIdle)/float64(dTotal))
	if usage < 0 {
		usage = 0
	}
	if usage > 100 {
		usage = 100
	}
	return usage
}

func readCPUStats() (idle, total uint64) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 8 {
			return 0, 0
		}
		var values [7]uint64
		for i := 0; i < 7; i++ {
			values[i], _ = strconv.ParseUint(fields[i+1], 10, 64)
		}
		idle = values[3] + values[4]
		for _, v := range values {
			total += v
		}
		return idle, total
	}
	return 0, 0
}

func getUptime() string {
	out, err := exec.Command("uptime", "-s").CombinedOutput()
	if err != nil {
		upout, _ := exec.Command("uptime").CombinedOutput()
		return strings.TrimSpace(string(upout))
	}
	return strings.TrimSpace(string(out))
}

func checkNginx() bool {
	_, err := exec.LookPath("nginx")
	return err == nil
}

type LinuxStats struct{}

func (s *LinuxStats) Collect() (*HostStatsPayload, error) {
	coll := NewStatsCollector("/var/run/docker.sock", 15*time.Second, nil)
	return coll.Collect()
}

type WindowsStats struct{}

func (s *WindowsStats) Collect() (*HostStatsPayload, error) {
	return &HostStatsPayload{}, nil
}

func NewStatsCollectorPlatform() interface{} {
	if runtime.GOOS == "windows" {
		return &WindowsStats{}
	}
	return &LinuxStats{}
}
