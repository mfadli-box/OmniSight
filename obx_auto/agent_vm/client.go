package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Config struct {
	APIURL            string
	RobotToken        string
	AgentHostID       string
	PollInterval      int
	HeartbeatInterval int
	VMPollInterval    int
	LogLevel          string
	DockerSocket      string
	NginxConfigDir    string
	LegoDir           string
	TrivyCacheDir     string
}

func LoadConfig() *Config {
	_ = os.Getenv("")

	apiURL := getEnv("API_URL", "http://localhost:3666")
	pollInterval := parseIntEnv("IVL_POOL", 30)
	heartbeatInterval := parseIntEnv("IVL_HEARTBEAT", 60)
	logLevel := getEnv("LOG_LEVEL", "info")
	dockerSocket := getEnv("SOC_DOCKER", "/var/run/docker.sock")
	nginxConfigDir := getEnv("DIR_NGINX", "/etc/nginx")
	legoDir := getEnv("DIR_LEGO", "/etc/lego")
	trivyCacheDir := getEnv("DIR_TRIVY", "/tmp/trivy")

	return &Config{
		APIURL:            apiURL,
		RobotToken:        getEnv("KEY_ROBOT", ""),
		AgentHostID:       getEnv("IDA_HOST", ""),
		PollInterval:      pollInterval,
		HeartbeatInterval: heartbeatInterval,
		VMPollInterval:    parseIntEnv("IVL_POOL", 60),
		LogLevel:          logLevel,
		DockerSocket:      dockerSocket,
		NginxConfigDir:    nginxConfigDir,
		LegoDir:           legoDir,
		TrivyCacheDir:     trivyCacheDir,
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func parseIntEnv(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if val := parseInt(v); val > 0 {
			return val
		}
	}
	return defaultVal
}

func parseInt(s string) int {
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}

type HTTPClient struct {
	baseURL    string
	robotToken string
	client     *http.Client
}

func NewHTTPClient(baseURL, token string) *HTTPClient {
	return &HTTPClient{
		baseURL:    baseURL,
		robotToken: token,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *HTTPClient) do(method, path string, body any) (*http.Response, []byte, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("json marshal: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.robotToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.robotToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do request: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("read body: %w", err)
	}
	resp.Body.Close()

	return resp, respBody, nil
}

func (c *HTTPClient) Heartbeat(hostID string) error {
	path := fmt.Sprintf("/rest/robot/VM01/ping?host_id=%s", hostID)
	resp, _, err := c.do("GET", path, nil)
	if err != nil {
		return fmt.Errorf("heartbeat: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("heartbeat: status %d", resp.StatusCode)
	}
	return nil
}

func (c *HTTPClient) DockerHeartbeat(hostID, version string, containers, images int) error {
	path := fmt.Sprintf("/rest/robot/DK01/ping?host_id=%s&version=%s&containers=%d&images=%d",
		hostID, version, containers, images)
	resp, _, err := c.do("GET", path, nil)
	if err != nil {
		return fmt.Errorf("docker heartbeat: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("docker heartbeat: status %d", resp.StatusCode)
	}
	return nil
}

type PendingJob struct {
	ID        string          `json:"id"`
	HostID    string          `json:"host_id"`
	RefID     string          `json:"ref_id"`
	JobType   string          `json:"job_type"`
	Action    string          `json:"action"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt string          `json:"created_at"`
}

type PendingResponse struct {
	Data    []PendingJob `json:"data"`
	Message string       `json:"message"`
}

func (c *HTTPClient) FetchPendingJobs(hostID string) ([]PendingJob, error) {
	path := fmt.Sprintf("/rest/robot/DK06/pending?host_id=%s", hostID)
	resp, body, err := c.do("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch pending: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("fetch pending: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result PendingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal pending: %w", err)
	}

	return result.Data, nil
}

type JobStatusUpdate struct {
	Status       string `json:"status"`
	Log          string `json:"log"`
	ErrorMessage string `json:"error_message"`
	ExitCode     int    `json:"exit_code"`
	CompletedAt  string `json:"completed_at"`
	ContainerID  string `json:"container_id"`
	Digest       string `json:"digest"`
	SizeBytes    int64  `json:"size_bytes"`
}

func (c *HTTPClient) ReportJob(jobID string, update JobStatusUpdate) error {
	path := "/rest/robot/DK06/report"
	updateWithID := struct {
		JobStatusUpdate
		ID string `json:"id"`
	}{
		JobStatusUpdate: update,
		ID:              jobID,
	}
	resp, body, err := c.do("POST", path, updateWithID)
	if err != nil {
		return fmt.Errorf("report job: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("report job: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type HostStatsPayload struct {
	HostID        string           `json:"host_id"`
	Hostname      string           `json:"hostname"`
	Uptime        string           `json:"uptime"`
	CPUCount      int              `json:"cpu_count"`
	CPUPercent    float64          `json:"cpu_percent"`
	MemTotal      int64            `json:"memory_total"`
	MemUsed       int64            `json:"memory_used"`
	DiskTotal     int64            `json:"disk_total"`
	DiskUsed      int64            `json:"disk_used"`
	Containers    []ContainerStats `json:"containers"`
	NetworkStats  []NetworkStats   `json:"network_stats"`
	DockerVersion string           `json:"docker_version"`
	DockerRunning bool             `json:"docker_running"`
	NginxRunning  bool             `json:"nginx_running"`
}

type ContainerStats struct {
	ID            string `json:"ID"`
	Name          string `json:"Name"`
	Image         string `json:"Image"`
	Status        string `json:"Status"`
	CPUPerc       string `json:"CPUPerc"`
	MemUsage      string `json:"MemUsage"`
	MemPerc       string `json:"MemPerc"`
	NetIO         string `json:"NetIO"`
	BlockIO       string `json:"BlockIO"`
	PIDs          string `json:"PIDs"`
	RestartCount  int    `json:"RestartCount"`
	UptimeSeconds int64  `json:"UptimeSeconds"`
}

type NetworkStats struct {
	Name      string `json:"name"`
	RxBytes   int64  `json:"rx_bytes"`
	TxBytes   int64  `json:"tx_bytes"`
	RxPackets int64  `json:"rx_packets"`
	TxPackets int64  `json:"tx_packets"`
}

func (c *HTTPClient) SendHostStats(stats HostStatsPayload) error {
	resp, body, err := c.do("POST", "/rest/robot/DK06/stats", stats)
	if err != nil {
		return fmt.Errorf("send stats: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("send stats: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type ExecLogRequest struct {
	HostID    string `json:"host_id"`
	Container string `json:"container"`
	Command   string `json:"command"`
	ExitCode  int    `json:"exit_code"`
	LogFile   string `json:"log_file"`
}

func (c *HTTPClient) ReportExecLog(req ExecLogRequest) error {
	resp, body, err := c.do("POST", "/rest/pages/DK06/exec-log", req)
	if err != nil {
		return fmt.Errorf("report exec log: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("report exec log: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type VulnScanRequest struct {
	HostID   string `json:"host_id"`
	RefID    string `json:"ref_id"`
	ScanType string `json:"scan_type"`
	Report   string `json:"report"`
	Critical int    `json:"critical"`
	High     int    `json:"high"`
	Medium   int    `json:"medium"`
	Low      int    `json:"low"`
}

func (c *HTTPClient) SubmitVulnScan(req VulnScanRequest) error {
	resp, body, err := c.do("POST", "/rest/pages/SC03/vuln-submit", req)
	if err != nil {
		return fmt.Errorf("submit vuln scan: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("submit vuln scan: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type FIMReportRequest struct {
	HostID    string      `json:"host_id"`
	RefID     string      `json:"ref_id"`
	Changes   []FIMChange `json:"changes"`
	ScannedAt string      `json:"scanned_at"`
}

func (c *HTTPClient) SubmitFIMReport(req FIMReportRequest) error {
	resp, body, err := c.do("POST", "/rest/pages/SC04/fim-submit", req)
	if err != nil {
		return fmt.Errorf("submit FIM report: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("submit FIM report: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type AlertNotifyRequest struct {
	HostID   string             `json:"host_id"`
	AlertID  string             `json:"alert_id"`
	RuleName string             `json:"rule_name"`
	Severity string             `json:"severity"`
	Message  string             `json:"message"`
	Metrics  map[string]float64 `json:"metrics"`
}

func (c *HTTPClient) SendAlertNotification(req AlertNotifyRequest) error {
	resp, body, err := c.do("POST", "/rest/pages/DK06/alert-notify", req)
	if err != nil {
		return fmt.Errorf("send alert notification: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("send alert notification: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type VMHostItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Status   string `json:"status"`
	VMCount  int    `json:"vm_count"`
	LastSeen string `json:"last_seen"`
	IsActive bool   `json:"is_active"`
}

type VMHostListResponse struct {
	Data    []VMHostItem `json:"data"`
	Message string       `json:"message"`
}

func (c *HTTPClient) FetchVMHosts() ([]VMHostItem, error) {
	resp, body, err := c.do("GET", "/rest/robot/VM01/hosts", nil)
	if err != nil {
		return nil, fmt.Errorf("fetch VM hosts: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("fetch VM hosts: status %d, body: %s", resp.StatusCode, string(body))
	}
	var result VMHostListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal VM hosts: %w", err)
	}
	return result.Data, nil
}

type VMPendingJob struct {
	ID        string          `json:"id"`
	HostID    string          `json:"host_id"`
	RefID     string          `json:"ref_id"`
	JobType   string          `json:"job_type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt string          `json:"created_at"`
}

type VMPendingResponse struct {
	Data    []VMPendingJob `json:"data"`
	Message string         `json:"message"`
}

func (c *HTTPClient) FetchVMPendingJobs(hostID string) ([]VMPendingJob, error) {
	path := fmt.Sprintf("/rest/robot/VM01/pending?host_id=%s", hostID)
	resp, body, err := c.do("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch VM pending: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("fetch VM pending: status %d, body: %s", resp.StatusCode, string(body))
	}
	var result VMPendingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal VM pending: %w", err)
	}
	return result.Data, nil
}

type VMReportRequest struct {
	JobID   string `json:"job_id" binding:"required"`
	HostID  string `json:"host_id" binding:"required"`
	JobType string `json:"job_type" binding:"required"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func (c *HTTPClient) ReportVMJob(req VMReportRequest) error {
	resp, body, err := c.do("POST", "/rest/robot/VM01/report", req)
	if err != nil {
		return fmt.Errorf("report VM job: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("report VM job: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type VMSyncEntry struct {
	VMID      string `json:"vm_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	OSType    string `json:"os_type"`
	CPUCount  int    `json:"cpu_count"`
	MemoryMB  int    `json:"memory_mb"`
	DiskGB    int    `json:"disk_gb"`
	IPAddress string `json:"ip_address"`
}

type VMSyncRequest struct {
	HostID string        `json:"host_id" binding:"required"`
	VMs    []VMSyncEntry `json:"vms" binding:"required"`
}

type VMSyncResponse struct {
	Message string `json:"message"`
	Synced  int    `json:"synced"`
}

func (c *HTTPClient) SyncVMs(req VMSyncRequest) (int, error) {
	resp, body, err := c.do("POST", "/rest/robot/VM01/sync", req)
	if err != nil {
		return 0, fmt.Errorf("sync VMs: %w", err)
	}
	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("sync VMs: status %d, body: %s", resp.StatusCode, string(body))
	}
	var result VMSyncResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("unmarshal sync response: %w", err)
	}
	return result.Synced, nil
}

type VMStatsEntry struct {
	VMID        string  `json:"vm_id"`
	Status      string  `json:"status"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemUsedMB   float64 `json:"mem_used_mb"`
	MemPercent  float64 `json:"mem_percent"`
	DiskUsedGB  float64 `json:"disk_used_gb"`
	DiskPercent float64 `json:"disk_percent"`
	NetRXBps    int64   `json:"net_rx_bps"`
	NetTXBps    int64   `json:"net_tx_bps"`
	UptimeSec   int64   `json:"uptime_sec"`
}

type VMStatsRequest struct {
	HostID string         `json:"host_id" binding:"required"`
	VMs    []VMStatsEntry `json:"vms" binding:"required"`
}

func (c *HTTPClient) SubmitVMStats(req VMStatsRequest) error {
	resp, body, err := c.do("POST", "/rest/robot/VM01/stats", req)
	if err != nil {
		return fmt.Errorf("submit VM stats: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("submit VM stats: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type VMPingResponse struct {
	Message string `json:"message"`
}

func (c *HTTPClient) PingVMHost(hostID, endpoint string, port int) error {
	path := fmt.Sprintf("/rest/pages/VM01/%s/ping?host=%s&port=%d", hostID, endpoint, port)
	resp, body, err := c.do("GET", path, nil)
	if err != nil {
		return fmt.Errorf("ping VM host: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("ping VM host: status %d, body: %s", resp.StatusCode, string(body))
	}
	_ = body
	return nil
}

type InventorySyncRequest struct {
	HostID     string               `json:"host_id" binding:"required"`
	Images     []InventoryImage     `json:"images"`
	Networks   []InventoryNetwork   `json:"networks"`
	Containers []InventoryContainer `json:"containers"`
	Compose    []InventoryCompose   `json:"compose"`
}

type InventorySyncResponse struct {
	Message string `json:"message"`
	Synced  int    `json:"synced"`
}

func (c *HTTPClient) SyncInventory(req InventorySyncRequest) (int, error) {
	var synced int
	err := WithRetry(DefaultRetry(), "SyncInventory", func() error {
		resp, body, err := c.do("POST", "/rest/robot/DK06/sync", req)
		if err != nil {
			return fmt.Errorf("sync inventory: %w", err)
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("sync inventory: status %d, body: %s", resp.StatusCode, string(body))
		}
		var result InventorySyncResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("unmarshal sync response: %w", err)
		}
		synced = result.Synced
		return nil
	})
	return synced, err
}
