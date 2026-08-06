package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := LoadConfig()

	log.Printf("[agent] MikroTik Agent starting...")
	log.Printf("[agent] API URL: %s", cfg.APIURL)
	log.Printf("[agent] Agent Host ID: %s", cfg.AgentHostID)
	log.Printf("[agent] Poll Interval: %ds", cfg.PollInterval)
	log.Printf("[agent] Heartbeat Interval: %ds", cfg.HeartbeatInterval)
	log.Printf("[agent] Syslog Bind: %s", cfg.SyslogBind)

	client := NewHTTPClient(cfg.APIURL, cfg.RobotToken)

	syslogServer := NewSyslogServer(cfg.SyslogBind, cfg.SyslogBufferDir, func(entries []LogEntry) error {
		if len(entries) == 0 {
			return nil
		}
		req := LogIngestRequest{
			DeviceID: cfg.AgentHostID,
			Logs:     entries,
		}
		return client.IngestLogs(req)
	})

	if err := syslogServer.Start(); err != nil {
		log.Printf("[agent] Failed to start syslog server: %v", err)
	} else {
		log.Printf("[agent] Syslog server running on %s", cfg.SyslogBind)
	}

	heartbeatTicker := time.NewTicker(time.Duration(cfg.HeartbeatInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	statsTicker := time.NewTicker(5 * time.Minute)

	stopCh := make(chan struct{})
	go heartbeatLoop(heartbeatTicker, client, cfg.AgentHostID)
	go pollLoop(pollTicker, client, cfg.AgentHostID)
	go statsLoop(statsTicker, client, cfg.AgentHostID)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Printf("[agent] Shutting down...")
	close(stopCh)
	heartbeatTicker.Stop()
	pollTicker.Stop()
	statsTicker.Stop()
	syslogServer.Stop()
	log.Printf("[agent] Stopped")
}

func heartbeatLoop(ticker *time.Ticker, client *HTTPClient, _ string) {
	for range ticker.C {
		if err := client.Heartbeat(); err != nil {
			log.Printf("[agent] Heartbeat error: %v", err)
		}
	}
}

func pollLoop(ticker *time.Ticker, client *HTTPClient, hostID string) {
	for range ticker.C {
		jobs, err := client.FetchPendingJobs(hostID)
		if err != nil {
			log.Printf("[agent] Fetch pending jobs error: %v", err)
			continue
		}

		for _, job := range jobs {
			go handleJob(job, client)
		}
	}
}

func statsLoop(ticker *time.Ticker, client *HTTPClient, hostID string) {
	for range ticker.C {
		stats, err := collectStats(hostID, client)
		if err != nil {
			log.Printf("[agent] Collect stats error: %v", err)
			continue
		}
		if err := client.SendStats(*stats); err != nil {
			log.Printf("[agent] Send stats error: %v", err)
		}
	}
}

func handleJob(job PendingJob, client *HTTPClient) {
	log.Printf("[agent] Processing job: %s (type: %s)", job.ID, job.JobType)

	var success bool
	var resultMsg string

	switch job.JobType {
	case "mikrotik_test":
		success, resultMsg = handleMikrotikTest(job, client)
	case "mikrotik_firewall_deploy":
		success, resultMsg = handleMikrotikFirewallDeploy(job, client)
	case "mikrotik_firewall_sync":
		success, resultMsg = handleMikrotikFirewallSync(job, client)
	case "mikrotik_backup":
		success, resultMsg = handleMikrotikBackup(job, client)
	case "mikrotik_restore":
		success, resultMsg = handleMikrotikRestore(job, client)
	case "mikrotik_syslog_setup":
		success, resultMsg = handleMikrotikSyslogSetup(job, client)
	default:
		success = false
		resultMsg = "Unknown job type: " + job.JobType
	}

	report := ReportRequest{
		DeviceID: job.HostID,
		JobType:  job.JobType,
		Success:  success,
		Message:  resultMsg,
		Result:   resultMsg,
	}

	if err := client.Report(report); err != nil {
		log.Printf("[agent] Report error for job %s: %v", job.ID, err)
	}
}

func handleMikrotikTest(job PendingJob, client *HTTPClient) (bool, string) {
	var payload struct {
		DeviceID string `json:"device_id"`
	}
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return false, "Invalid payload: " + err.Error()
	}

	cred, err := client.GetDeviceCredential(payload.DeviceID)
	if err != nil {
		return false, "Failed to get device credential: " + err.Error()
	}

	ros := NewRouterOSClient(cred.IPAddress, cred.APIPort, cred.APITLS, cred.Username, cred.Password)
	result := ros.TestConnection()

	if result.Success {
		return true, "Connection successful"
	}
	return false, "Connection failed: " + result.Error
}

func handleMikrotikFirewallDeploy(job PendingJob, client *HTTPClient) (bool, string) {
	var payload struct {
		DeviceID    string `json:"device_id"`
		Chain       string `json:"chain"`
		Action      string `json:"action"`
		SrcAddress  string `json:"src_address"`
		DstAddress  string `json:"dst_address"`
		Protocol    string `json:"protocol"`
		SrcPort     string `json:"src_port"`
		DstPort     string `json:"dst_port"`
		Comment     string `json:"comment"`
		RuleID      string `json:"rule_id"`
		RuleEnabled bool   `json:"rule_enabled"`
	}
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return false, "Invalid payload: " + err.Error()
	}

	cred, err := client.GetDeviceCredential(payload.DeviceID)
	if err != nil {
		return false, "Failed to get device credential: " + err.Error()
	}

	ros := NewRouterOSClient(cred.IPAddress, cred.APIPort, cred.APITLS, cred.Username, cred.Password)

	if payload.RuleID != "" {
		if payload.RuleEnabled {
			if err := ros.EnableFirewallRule(payload.RuleID); err != nil {
				return false, "Failed to enable rule: " + err.Error()
			}
			return true, "Rule " + payload.RuleID + " enabled"
		}
		if err := ros.DisableFirewallRule(payload.RuleID); err != nil {
			return false, "Failed to disable rule: " + err.Error()
		}
		return true, "Rule " + payload.RuleID + " disabled"
	}

	chain := payload.Chain
	if chain == "" {
		chain = "forward"
	}
	action := payload.Action
	if action == "" {
		action = "accept"
	}
	comment := payload.Comment
	if comment == "" {
		comment = "Added by OmniSight Agent"
	}

	if err := ros.AddFirewallRule(chain, action, comment, false); err != nil {
		return false, "Failed to add rule: " + err.Error()
	}
	return true, "Firewall rule added to chain " + chain
}

func handleMikrotikFirewallSync(job PendingJob, client *HTTPClient) (bool, string) {
	var payload struct {
		DeviceID string `json:"device_id"`
	}
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return false, "Invalid payload: " + err.Error()
	}

	cred, err := client.GetDeviceCredential(payload.DeviceID)
	if err != nil {
		return false, "Failed to get device credential: " + err.Error()
	}

	ros := NewRouterOSClient(cred.IPAddress, cred.APIPort, cred.APITLS, cred.Username, cred.Password)
	rules, err := ros.GetFirewallRules()
	if err != nil {
		return false, "Failed to get firewall rules: " + err.Error()
	}

	result := map[string]interface{}{
		"device_id": payload.DeviceID,
		"rules":     rules,
		"count":     len(rules),
	}
	data, _ := json.Marshal(result)

	return true, string(data)
}

func handleMikrotikBackup(job PendingJob, client *HTTPClient) (bool, string) {
	var payload struct {
		DeviceID   string `json:"device_id"`
		BackupName string `json:"backup_name"`
		Password   string `json:"password"`
		ExportFile bool   `json:"export_file"`
	}
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return false, "Invalid payload: " + err.Error()
	}

	cred, err := client.GetDeviceCredential(payload.DeviceID)
	if err != nil {
		return false, "Failed to get device credential: " + err.Error()
	}

	ros := NewRouterOSClient(cred.IPAddress, cred.APIPort, cred.APITLS, cred.Username, cred.Password)

	backupName := payload.BackupName
	if backupName == "" {
		backupName = "omnisight_backup"
	}

	if payload.ExportFile {
		content, err := ros.ExportBackup(backupName)
		if err != nil {
			return false, "Failed to export config: " + err.Error()
		}
		return true, "Export saved (size: " + string(rune(len(content))) + " bytes)"
	}

	if err := ros.CreateBackup(backupName, payload.Password); err != nil {
		return false, "Failed to create backup: " + err.Error()
	}
	return true, "Backup created: " + backupName
}

func handleMikrotikRestore(job PendingJob, client *HTTPClient) (bool, string) {
	var payload struct {
		DeviceID   string `json:"device_id"`
		BackupName string `json:"backup_name"`
		LoadBackup bool   `json:"load_backup"`
	}
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return false, "Invalid payload: " + err.Error()
	}

	cred, err := client.GetDeviceCredential(payload.DeviceID)
	if err != nil {
		return false, "Failed to get device credential: " + err.Error()
	}

	ros := NewRouterOSClient(cred.IPAddress, cred.APIPort, cred.APITLS, cred.Username, cred.Password)

	backupName := payload.BackupName
	if backupName == "" {
		return false, "Backup name is required"
	}

	if payload.LoadBackup {
		if err := ros.LoadBackup(backupName); err != nil {
			return false, "Failed to load backup: " + err.Error()
		}
		return true, "Backup loaded: " + backupName
	}

	return false, "Specify load_backup=true to load the backup"
}

func handleMikrotikSyslogSetup(job PendingJob, client *HTTPClient) (bool, string) {
	var payload struct {
		DeviceID   string `json:"device_id"`
		RemoteAddr string `json:"remote_addr"`
		RemotePort int    `json:"remote_port"`
		SrcAddress string `json:"src_address"`
		Action     string `json:"action"`
	}
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return false, "Invalid payload: " + err.Error()
	}

	cred, err := client.GetDeviceCredential(payload.DeviceID)
	if err != nil {
		return false, "Failed to get device credential: " + err.Error()
	}

	ros := NewRouterOSClient(cred.IPAddress, cred.APIPort, cred.APITLS, cred.Username, cred.Password)

	remoteAddr := payload.RemoteAddr
	if remoteAddr == "" {
		return false, "Remote address is required"
	}
	remotePort := payload.RemotePort
	if remotePort == 0 {
		remotePort = 514
	}
	action := payload.Action
	if action == "" {
		action = "add"
	}

	switch action {
	case "add":
		if err := ros.AddSyslogTarget(remoteAddr, remotePort); err != nil {
			return false, "Failed to add syslog target: " + err.Error()
		}
		return true, fmt.Sprintf("Syslog target added: %s:%d", remoteAddr, remotePort)
	case "remove":
		if err := ros.RemoveSyslogTarget(remoteAddr); err != nil {
			return false, "Failed to remove syslog target: " + err.Error()
		}
		return true, fmt.Sprintf("Syslog target removed: %s", remoteAddr)
	case "list":
		targets, err := ros.ListSyslogTargets()
		if err != nil {
			return false, "Failed to list syslog targets: " + err.Error()
		}
		return true, fmt.Sprintf("Syslog targets: %v", targets)
	default:
		return false, "Unknown action: " + action + ". Use: add, remove, or list"
	}
}

func collectStats(_ string, client *HTTPClient) (*StatsRequest, error) {
	devices, err := client.ListDevices()
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices assigned to this agent")
	}

	for _, dev := range devices {
		cred := DeviceCredential{
			ID:        dev.ID,
			IPAddress: dev.IPAddress,
			APIPort:   dev.APIPort,
			APITLS:    dev.APITLS,
			Username:  dev.Username,
			Password:  "",
		}
		creds, err := client.GetDeviceCredential(dev.ID)
		if err != nil {
			log.Printf("[agent] Failed to get credential for device %s: %v", dev.ID, err)
			continue
		}
		cred.Password = creds.Password

		ros := NewRouterOSClient(cred.IPAddress, cred.APIPort, cred.APITLS, cred.Username, cred.Password)
		res, err := ros.GetSystemResources()
		if err != nil {
			log.Printf("[agent] Failed to get stats from device %s: %v", dev.ID, err)
			continue
		}

		ifaces, err := ros.GetInterfaces()
		if err != nil {
			log.Printf("[agent] Failed to get interfaces from device %s: %v", dev.ID, err)
			continue
		}

		var ifaceStats []InterfaceStats
		for _, iface := range ifaces {
			ifaceStats = append(ifaceStats, InterfaceStats{
				Name:     iface.Name,
				Type:     iface.Type,
				RXBytes:  iface.RXBytes,
				TXBytes:  iface.TXBytes,
				Running:  iface.Running,
				Disabled: iface.Disabled,
			})
		}

		return &StatsRequest{
			DeviceID:   dev.ID,
			CPUCount:   res.CPUCount,
			CPULoad:    res.CPULoad,
			MemTotal:   res.MemoryTotal,
			MemUsed:    res.MemoryUsed,
			DiskTotal:  res.DiskTotal,
			DiskUsed:   res.DiskUsed,
			Uptime:     res.Uptime,
			Interfaces: ifaceStats,
		}, nil
	}

	return nil, fmt.Errorf("failed to collect stats from any device")
}
