package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	APIURL            string
	RobotToken        string
	AgentHostID       string
	PollInterval      int
	HeartbeatInterval int
	LogLevel          string
	SyslogBind        string
	SyslogBufferDir   string
	DecryptKeyPath    string
}

func LoadConfig() *Config {
	godotenv.Load()

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3666"
	}

	pollInterval := 30
	if envPoll := os.Getenv("IVL_POOL"); envPoll != "" {
		if val := parseInt(envPoll); val > 0 {
			pollInterval = val
		}
	}

	heartbeatInterval := 60
	if envHB := os.Getenv("IVL_HEARTBEAT"); envHB != "" {
		if val := parseInt(envHB); val > 0 {
			heartbeatInterval = val
		}
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	syslogBind := os.Getenv("IPP_SYSLOG")
	if syslogBind == "" {
		syslogBind = "0.0.0.0:514"
	}

	syslogBufferDir := os.Getenv("DIR_SYSLOG")
	if syslogBufferDir == "" {
		syslogBufferDir = "./syslog_buffer"
	}

	return &Config{
		APIURL:            apiURL,
		RobotToken:        os.Getenv("KEY_ROBOT"),
		AgentHostID:       os.Getenv("IDA_HOST"),
		PollInterval:      pollInterval,
		HeartbeatInterval: heartbeatInterval,
		LogLevel:          logLevel,
		SyslogBind:        syslogBind,
		SyslogBufferDir:   syslogBufferDir,
		DecryptKeyPath:    os.Getenv("DIR_DECRYPT"),
	}
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
			Timeout: 30 * time.Second,
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

func (c *HTTPClient) Heartbeat() error {
	path := fmt.Sprintf("/rest/robot/IM02/ping?device_id=%s", c.robotToken)
	resp, _, err := c.do("GET", path, nil)
	if err != nil {
		return fmt.Errorf("heartbeat: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("heartbeat: status %d", resp.StatusCode)
	}
	log.Printf("[agent] Heartbeat sent")
	return nil
}

type PendingJob struct {
	ID        string `json:"id"`
	HostID    string `json:"host_id"`
	RefID     string `json:"ref_id"`
	JobType   string `json:"job_type"`
	Payload   string `json:"payload"`
	CreatedAt string `json:"created_at"`
}

type PendingResponse struct {
	Data    []PendingJob `json:"data"`
	Message string       `json:"message"`
}

func (c *HTTPClient) FetchPendingJobs(hostID string) ([]PendingJob, error) {
	path := fmt.Sprintf("/rest/robot/IM02/pending?host_id=%s", hostID)
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

type ReportRequest struct {
	DeviceID string `json:"device_id"`
	JobType  string `json:"job_type"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Result   string `json:"result"`
}

func (c *HTTPClient) Report(req ReportRequest) error {
	resp, body, err := c.do("POST", "/rest/robot/IM02/report", req)
	if err != nil {
		return fmt.Errorf("report: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("report: status %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

type StatsRequest struct {
	DeviceID   string           `json:"device_id"`
	CPUCount   int              `json:"cpu_count"`
	CPULoad    float64          `json:"cpu_load"`
	MemTotal   int64            `json:"memory_total"`
	MemUsed    int64            `json:"memory_used"`
	DiskTotal  int64            `json:"disk_total"`
	DiskUsed   int64            `json:"disk_used"`
	Uptime     string           `json:"uptime"`
	Interfaces []InterfaceStats `json:"interfaces"`
}

type InterfaceStats struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	RXBytes  int64  `json:"rx_bytes"`
	TXBytes  int64  `json:"tx_bytes"`
	Running  bool   `json:"running"`
	Disabled bool   `json:"disabled"`
}

func (c *HTTPClient) SendStats(req StatsRequest) error {
	resp, body, err := c.do("POST", "/rest/robot/IM02/stats", req)
	if err != nil {
		return fmt.Errorf("send stats: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("send stats: status %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

type LogEntry struct {
	Topics   string `json:"topics"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	SourceIP string `json:"source_ip"`
	Time     string `json:"time"`
}

type LogIngestRequest struct {
	DeviceID string     `json:"device_id"`
	Logs     []LogEntry `json:"logs"`
}

func (c *HTTPClient) IngestLogs(req LogIngestRequest) error {
	resp, body, err := c.do("POST", "/rest/robot/IM02/logs/ingest", req)
	if err != nil {
		return fmt.Errorf("ingest logs: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("ingest logs: status %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

type DeviceCredential struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip_address"`
	APIPort   int    `json:"api_port"`
	APITLS    bool   `json:"api_tls"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (c *HTTPClient) GetDeviceCredential(deviceID string) (*DeviceCredential, error) {
	path := fmt.Sprintf("/rest/robot/IM02/%s/credential", deviceID)
	resp, body, err := c.do("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("get device: status %d", resp.StatusCode)
	}

	var result struct {
		Data DeviceCredential `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal device: %w", err)
	}
	return &result.Data, nil
}

type DeviceInfo struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip_address"`
	APIPort   int    `json:"api_port"`
	APITLS    bool   `json:"api_tls"`
	Username  string `json:"username"`
	Name      string `json:"name"`
}

func (c *HTTPClient) ListDevices() ([]DeviceInfo, error) {
	path := "/rest/robot/IM02/devices"
	resp, body, err := c.do("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("list devices: status %d", resp.StatusCode)
	}

	var result struct {
		Data []DeviceInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal devices: %w", err)
	}
	return result.Data, nil
}
