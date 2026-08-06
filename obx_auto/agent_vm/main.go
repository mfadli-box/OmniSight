package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	cfg := LoadConfig()

	log.Printf("[agent] VM Agent starting...")
	log.Printf("[agent] API URL: %s", cfg.APIURL)
	log.Printf("[agent] Agent Host ID: %s", cfg.AgentHostID)
	log.Printf("[agent] Poll Interval: %ds", cfg.PollInterval)
	log.Printf("[agent] Heartbeat Interval: %ds", cfg.HeartbeatInterval)
	log.Printf("[agent] Docker Socket: %s", cfg.DockerSocket)

	client := NewHTTPClient(cfg.APIURL, cfg.RobotToken)

	composeExec := NewComposeExecutor()
	containerExec := NewContainerExecutor()
	imageExec := NewImageExecutor()
	nginxMgr := NewNginxManager(cfg.NginxConfigDir)
	legoSSL := NewLegoSSL(cfg.LegoDir, "", nil)
	firewallMgr := NewFirewallManager()
	updateMgr := NewUpdateManager("")
	gitopsMgr := NewGitOpsManager("/var/lib/agent_vm/gitops")
	trivyScanner := NewTrivyScanner(cfg.TrivyCacheDir)
	statsCollector := NewStatsCollector(cfg.DockerSocket, 15*time.Second, func(stats HostStatsPayload) {
		stats.HostID = cfg.AgentHostID
		if err := client.SendHostStats(stats); err != nil {
			log.Printf("[agent] Send stats error: %v", err)
		}
	})

	statsCollector.Start()

	vmPollTicker := time.NewTicker(time.Duration(cfg.VMPollInterval) * time.Second)
	inventorySyncTicker := time.NewTicker(60 * time.Second)

	heartbeatTicker := time.NewTicker(time.Duration(cfg.HeartbeatInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)

	stopCh := make(chan struct{})
	go heartbeatLoop(heartbeatTicker, client, cfg.AgentHostID)
	go pollLoop(pollTicker, client, cfg.AgentHostID, composeExec, containerExec, imageExec, nginxMgr, legoSSL, firewallMgr, updateMgr, gitopsMgr, trivyScanner)
	go vmPollLoop(vmPollTicker, client, cfg.AgentHostID)
	go inventorySyncLoop(inventorySyncTicker, client, cfg.AgentHostID)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Printf("[agent] Shutting down...")
	close(stopCh)
	heartbeatTicker.Stop()
	pollTicker.Stop()
	vmPollTicker.Stop()
	inventorySyncTicker.Stop()
	statsCollector.Stop()
	log.Printf("[agent] Stopped")
}

func heartbeatLoop(ticker *time.Ticker, client *HTTPClient, hostID string) {
	for range ticker.C {
		if err := client.Heartbeat(hostID); err != nil {
			log.Printf("[agent] Heartbeat error: %v", err)
		}
		dockerVersion, containers, images := getDockerInfo()
		if containers >= 0 {
			if err := client.DockerHeartbeat(hostID, dockerVersion, containers, images); err != nil {
				log.Printf("[agent] Docker Heartbeat error: %v", err)
			}
		}
	}
}

func getDockerInfo() (version string, containers, images int) {
	out, err := exec.Command("docker", "version", "--format", "{{.Server.Version}}").CombinedOutput()
	if err == nil {
		version = strings.TrimSpace(string(out))
	}
	out, err = exec.Command("docker", "ps", "-q").CombinedOutput()
	if err == nil {
		containers = len(strings.Fields(strings.TrimSpace(string(out))))
	}
	out, err = exec.Command("docker", "images", "-q").CombinedOutput()
	if err == nil {
		images = len(strings.Fields(strings.TrimSpace(string(out))))
	}
	return version, containers, images
}

func pollLoop(
	ticker *time.Ticker,
	client *HTTPClient,
	hostID string,
	composeExec *ComposeExecutor,
	containerExec *ContainerExecutor,
	imageExec *ImageExecutor,
	nginxMgr *NginxManager,
	legoSSL *LegoSSL,
	firewallMgr *FirewallManager,
	updateMgr *UpdateManager,
	gitopsMgr *GitOpsManager,
	trivyScanner *TrivyScanner,
) {
	for range ticker.C {
		jobs, err := client.FetchPendingJobs(hostID)
		if err != nil {
			log.Printf("[agent] Fetch pending jobs error: %v", err)
			continue
		}

		for _, job := range jobs {
			go handleJob(job, client, composeExec, containerExec, imageExec, nginxMgr, legoSSL, firewallMgr, updateMgr, gitopsMgr, trivyScanner)
		}
	}
}

func handleJob(
	job PendingJob,
	client *HTTPClient,
	composeExec *ComposeExecutor,
	containerExec *ContainerExecutor,
	imageExec *ImageExecutor,
	nginxMgr *NginxManager,
	legoSSL *LegoSSL,
	firewallMgr *FirewallManager,
	updateMgr *UpdateManager,
	gitopsMgr *GitOpsManager,
	trivyScanner *TrivyScanner,
) {
	log.Printf("[agent] Processing job: %s (type: %s, action: %s)", job.ID, job.JobType, job.Action)

	startTime := time.Now()
	var success bool
	var resultLog string
	var errorMsg string
	var exitCode int
	var containerID, imageDigest string
	var imageSize int64

	switch job.JobType {
	case "compose_deploy":
		success, resultLog, errorMsg = handleComposeDeploy(job, composeExec)
	case "compose_down":
		success, resultLog, errorMsg = handleComposeDown(job, composeExec)
	case "container_action":
		success, resultLog, errorMsg = handleContainerAction(job, containerExec)
	case "container_deploy":
		success, resultLog, errorMsg, containerID = handleContainerDeploy(job, containerExec)
	case "container_exec":
		success, exitCode, resultLog, errorMsg = handleContainerExec(job, containerExec)
	case "image_pull":
		success, resultLog, errorMsg, imageDigest, imageSize = handleImagePull(job, imageExec)
	case "image_remove":
		success, resultLog, errorMsg = handleImageRemove(job, imageExec)
	case "image_prune":
		success, resultLog, errorMsg = handleImagePrune(job, imageExec)
	case "nginx_config":
		success, resultLog, errorMsg = handleNginxConfig(job, nginxMgr)
	case "nginx_validate":
		success, resultLog, errorMsg = handleNginxValidate(job, nginxMgr)
	case "nginx_ssl":
		success, resultLog, errorMsg = handleNginxSSL(job, nginxMgr, legoSSL)
	case "firewall_deploy":
		success, resultLog, errorMsg = handleFirewallDeploy(job, firewallMgr)
	case "firewall_test":
		success, resultLog, errorMsg = handleFirewallTest(job, firewallMgr)
	case "firewall_backup":
		success, resultLog, errorMsg = handleFirewallBackup(job, firewallMgr)
	case "update_os":
		success, resultLog, errorMsg = handleUpdateOS(job, updateMgr)
	case "update_docker":
		success, resultLog, errorMsg = handleUpdateDocker(job, updateMgr)
	case "update_security":
		success, resultLog, errorMsg = handleUpdateSecurity(job, updateMgr)
	case "host_reboot":
		success, errorMsg = handleHostReboot(job, updateMgr)
	case "gitops_sync":
		success, resultLog, errorMsg = handleGitOpsSync(job, gitopsMgr)
	case "vuln_scan":
		success, resultLog, errorMsg = handleVulnScan(job, trivyScanner, client)
	case "fim_check":
		success, resultLog, errorMsg = handleFIMCheck(job, client)
	default:
		success = false
		errorMsg = "Unknown job type: " + job.JobType
	}

	completedAt := time.Now().UTC().Format(time.RFC3339)
	duration := time.Since(startTime).Seconds()

	update := JobStatusUpdate{
		Status:       map[bool]string{true: "success", false: "failed"}[success],
		Log:          fmt.Sprintf("[%.1fs] %s", duration, resultLog),
		ErrorMessage: errorMsg,
		ExitCode:     exitCode,
		CompletedAt:  completedAt,
		ContainerID:  containerID,
		Digest:       imageDigest,
		SizeBytes:    imageSize,
	}

	if err := client.ReportJob(job.ID, update); err != nil {
		log.Printf("[agent] Report job error for %s: %v", job.ID, err)
	}
}

type composePayload struct {
	ComposeFile string   `json:"compose_file"`
	ProjectName string   `json:"project_name"`
	EnvVars     []string `json:"env_vars"`
}

func handleComposeDeploy(job PendingJob, exec *ComposeExecutor) (bool, string, string) {
	var payload composePayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	out, err := exec.Deploy(payload.ComposeFile, payload.ProjectName, payload.EnvVars)
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleComposeDown(job PendingJob, exec *ComposeExecutor) (bool, string, string) {
	var payload composePayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	out, err := exec.Down(payload.ComposeFile, payload.ProjectName, payload.EnvVars)
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

type containerPayload struct {
	Container string `json:"container"`
	Command   string `json:"command"`
}

func handleContainerDeploy(job PendingJob, exec *ContainerExecutor) (bool, string, string, string) {
	var payload ContainerRunConfig
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error(), ""
	}
	out, err := exec.Run(payload)
	if err != nil {
		return false, out, err.Error(), ""
	}
	return true, out, "", strings.TrimSpace(out)
}

func handleContainerAction(job PendingJob, exec *ContainerExecutor) (bool, string, string) {
	var payload containerPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	action := job.Action
	var out string
	var err error
	switch action {
	case "start":
		out, err = exec.Start(payload.Container)
	case "stop":
		out, err = exec.Stop(payload.Container)
	case "restart":
		out, err = exec.Restart(payload.Container)
	case "rm":
		out, err = exec.Remove(payload.Container, true)
	default:
		return false, "", "Unknown action: " + action
	}
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleContainerExec(job PendingJob, exec *ContainerExecutor) (bool, int, string, string) {
	var payload containerPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, 1, "", "Invalid payload: " + err.Error()
	}
	out, exitCode, err := exec.Exec(payload.Container, payload.Command)
	if err != nil {
		return false, exitCode, out, err.Error()
	}
	return true, exitCode, out, ""
}

type imagePayload struct {
	Image string `json:"image"`
}

func handleImagePull(job PendingJob, exec *ImageExecutor) (bool, string, string, string, int64) {
	var payload imagePayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error(), "", 0
	}
	out, err := exec.Pull(payload.Image)
	if err != nil {
		return false, out, err.Error(), "", 0
	}
	digest, size, err := exec.Inspect(payload.Image)
	if err != nil {
		return true, out, "", "", 0
	}
	return true, out, "", digest, size
}

func handleImagePrune(_ PendingJob, exec *ImageExecutor) (bool, string, string) {
	out, err := exec.Prune(true)
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleImageRemove(job PendingJob, exec *ImageExecutor) (bool, string, string) {
	var payload imagePayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	out, err := exec.Remove(payload.Image, true)
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

type nginxPayload struct {
	SiteName string `json:"site_name"`
	Content  string `json:"content"`
	Action   string `json:"action"`
}

func handleNginxConfig(job PendingJob, mgr *NginxManager) (bool, string, string) {
	var payload nginxPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	path, err := mgr.DeployConfig(payload.SiteName, payload.Content)
	if err != nil {
		return false, path, err.Error()
	}
	if payload.Action == "enable" {
		if err := mgr.EnableSite(payload.SiteName); err != nil {
			return false, path, err.Error()
		}
	}
	out, err := mgr.Reload()
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleNginxValidate(job PendingJob, mgr *NginxManager) (bool, string, string) {
	var payload nginxPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	ok, out, err := mgr.TestConfigContent(payload.Content)
	if err != nil {
		return false, out, err.Error()
	}
	if !ok {
		return false, out, "config test failed"
	}
	return true, out, ""
}

func handleNginxSSL(job PendingJob, _ *NginxManager, legoSSL *LegoSSL) (bool, string, string) {
	var payload struct {
		Domains []string `json:"domains"`
		Email   string   `json:"email"`
		Action  string   `json:"action"`
	}
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}

	if len(payload.Domains) == 0 {
		return false, "", "At least one domain is required"
	}
	if payload.Email == "" {
		return false, "", "Email is required for Let's Encrypt"
	}

	legoSSL.email = payload.Email
	legoSSL.domains = payload.Domains

	action := payload.Action
	if action == "" {
		action = "provision"
	}

	switch action {
	case "provision":
		certPath, keyPath, err := legoSSL.Provision()
		if err != nil {
			return false, "", "Lego provision failed: " + err.Error()
		}
		return true, fmt.Sprintf("Cert: %s, Key: %s", certPath, keyPath), ""
	case "renew":
		domain := payload.Domains[0]
		output, err := legoSSL.Renew(domain)
		if err != nil {
			return false, "", "Lego renew failed: " + err.Error()
		}
		return true, output, ""
	case "list":
		certs, err := legoSSL.List()
		if err != nil {
			return false, "", "Lego list failed: " + err.Error()
		}
		result := fmt.Sprintf("Found %d certificates", len(certs))
		return true, result, ""
	default:
		return false, "", "Unknown action: " + action + ". Use: provision, renew, or list"
	}
}

type firewallPayload struct {
	Content string `json:"content"`
}

func handleFirewallDeploy(job PendingJob, mgr *FirewallManager) (bool, string, string) {
	var payload firewallPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	out, err := mgr.Deploy(payload.Content)
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleFirewallTest(job PendingJob, mgr *FirewallManager) (bool, string, string) {
	var payload firewallPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	ok, out, err := mgr.Test(payload.Content)
	if err != nil {
		return false, out, err.Error()
	}
	if !ok {
		return false, out, "firewall test failed"
	}
	return true, out, ""
}

func handleFirewallBackup(_ PendingJob, mgr *FirewallManager) (bool, string, string) {
	out, err := mgr.Backup()
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleUpdateOS(_ PendingJob, mgr *UpdateManager) (bool, string, string) {
	out, err := mgr.UpdateOS()
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleUpdateDocker(_ PendingJob, mgr *UpdateManager) (bool, string, string) {
	out, err := mgr.UpdateDocker()
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleUpdateSecurity(_ PendingJob, mgr *UpdateManager) (bool, string, string) {
	out, err := mgr.UpdateSecurity()
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

func handleHostReboot(_ PendingJob, mgr *UpdateManager) (bool, string) {
	if err := mgr.Reboot(); err != nil {
		return false, err.Error()
	}
	return true, ""
}

type gitopsPayload struct {
	RepoURL   string `json:"repo_url"`
	Branch    string `json:"branch"`
	DeployCmd string `json:"deploy_cmd"`
}

func handleGitOpsSync(job PendingJob, mgr *GitOpsManager) (bool, string, string) {
	var payload gitopsPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}
	out, err := mgr.Sync(payload.RepoURL, payload.Branch, payload.DeployCmd)
	if err != nil {
		return false, out, err.Error()
	}
	return true, out, ""
}

type vulnPayload struct {
	HostID   string `json:"host_id"`
	RefID    string `json:"ref_id"`
	ScanType string `json:"scan_type"`
	Target   string `json:"target"`
}

func handleVulnScan(job PendingJob, scanner *TrivyScanner, client *HTTPClient) (bool, string, string) {
	var payload vulnPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}

	var result *VulnScanResult
	var err error

	switch payload.ScanType {
	case "image":
		result, err = scanner.ScanImage(payload.Target)
	case "filesystem":
		result, err = scanner.ScanFilesystem(payload.Target)
	default:
		return false, "", "Unknown scan type: " + payload.ScanType
	}

	if err != nil {
		return false, "", err.Error()
	}

	req := VulnScanRequest{
		HostID:   payload.HostID,
		RefID:    payload.RefID,
		ScanType: payload.ScanType,
		Report:   fmt.Sprintf("Critical: %d, High: %d, Medium: %d, Low: %d", result.Critical, result.High, result.Medium, result.Low),
		Critical: result.Critical,
		High:     result.High,
		Medium:   result.Medium,
		Low:      result.Low,
	}

	if err := client.SubmitVulnScan(req); err != nil {
		return false, "", err.Error()
	}

	return true, req.Report, ""
}

type fimPayload struct {
	HostID string   `json:"host_id"`
	RefID  string   `json:"ref_id"`
	Paths  []string `json:"paths"`
}

func handleFIMCheck(job PendingJob, client *HTTPClient) (bool, string, string) {
	var payload fimPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "", "Invalid payload: " + err.Error()
	}

	monitor := NewFIMMonitor("", "", payload.Paths)
	changes, err := monitor.Check()
	if err != nil {
		return false, "", err.Error()
	}

	var fimChanges []FIMChange
	for _, c := range changes {
		fimChanges = append(fimChanges, FIMChange{
			Path:       c.Path,
			OldHash:    c.OldHash,
			NewHash:    c.NewHash,
			ChangeType: c.ChangeType,
		})
	}

	req := FIMReportRequest{
		HostID:    payload.HostID,
		RefID:     payload.RefID,
		Changes:   fimChanges,
		ScannedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := client.SubmitFIMReport(req); err != nil {
		return false, "", err.Error()
	}

	if len(changes) == 0 {
		return true, "No changes detected", ""
	}
	return true, fmt.Sprintf("Detected %d changes", len(changes)), ""
}

func vmPollLoop(ticker *time.Ticker, client *HTTPClient, agentHostID string) {
	for range ticker.C {
		hosts, err := client.FetchVMHosts()
		if err != nil {
			log.Printf("[agent] Fetch VM hosts error: %v", err)
			continue
		}

		for _, host := range hosts {
			if !host.IsActive {
				continue
			}
			jobs, err := client.FetchVMPendingJobs(host.ID)
			if err != nil {
				log.Printf("[agent] Fetch VM pending jobs for host %s error: %v", host.ID, err)
				continue
			}

			for _, job := range jobs {
				go handleVMJob(job, client, host)
			}
		}
	}
}

func inventorySyncLoop(ticker *time.Ticker, client *HTTPClient, hostID string) {
	collector := NewInventoryCollector()
	for range ticker.C {
		data, err := collector.Collect(hostID)
		if err != nil {
			log.Printf("[agent] Inventory collect error: %v", err)
			continue
		}
		synced, err := client.SyncInventory(InventorySyncRequest{
			HostID:     data.HostID,
			Images:     data.Images,
			Networks:   data.Networks,
			Containers: data.Containers,
			Compose:    data.Compose,
		})
		if err != nil {
			log.Printf("[agent] Inventory sync error: %v", err)
		} else {
			log.Printf("[agent] Inventory synced: %d items", synced)
		}
	}
}

type vmSyncPayload struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

func handleVMJob(job VMPendingJob, client *HTTPClient, host VMHostItem) {
	log.Printf("[agent] Processing VM job: %s (type: %s)", job.ID, job.JobType)

	startTime := time.Now()
	var success bool
	var message string
	var result string

	switch job.JobType {
	case "vm_test":
		success, message = handleVMTest(job, client, host)
	case "vm_sync":
		success, message, result = handleVMSync(job, client, host)
	case "vm_stats_collect":
		success, message, result = handleVMStatsCollect(job, client, host)
	default:
		success = false
		message = "Unknown job type: " + job.JobType
	}

	duration := time.Since(startTime).Seconds()

	req := VMReportRequest{
		JobID:   job.ID,
		HostID:  host.ID,
		JobType: job.JobType,
		Success: success,
		Message: message,
		Result:  fmt.Sprintf("[%.1fs] %s", duration, result),
	}

	if err := client.ReportVMJob(req); err != nil {
		log.Printf("[agent] Report VM job error for %s: %v", job.ID, err)
	}
}

func handleVMTest(job VMPendingJob, client *HTTPClient, host VMHostItem) (bool, string) {
	var payload vmSyncPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "Invalid payload: " + err.Error()
	}

	address := fmt.Sprintf("%s:%d", host.Host, host.Port)
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return false, "Connection failed: " + err.Error()
	}
	conn.Close()
	return true, fmt.Sprintf("Successfully connected to %s", address)
}

func handleVMSync(job VMPendingJob, client *HTTPClient, host VMHostItem) (bool, string, string) {
	var payload vmSyncPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "Invalid payload: " + err.Error(), ""
	}

	vms, err := queryHypervisorVMs(host.Type, host.Host, host.Port, host.Username, host.ID)
	if err != nil {
		return false, "Failed to query VMs: " + err.Error(), ""
	}

	synced, err := client.SyncVMs(VMSyncRequest{
		HostID: host.ID,
		VMs:    vms,
	})
	if err != nil {
		return false, "Failed to sync VMs: " + err.Error(), ""
	}

	return true, fmt.Sprintf("Synced %d VMs", synced), fmt.Sprintf("synced=%d", synced)
}

func handleVMStatsCollect(job VMPendingJob, client *HTTPClient, host VMHostItem) (bool, string, string) {
	var payload vmSyncPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return false, "Invalid payload: " + err.Error(), ""
	}

	stats, err := collectVMStats(host.Type, host.Host, host.Port, host.Username, host.ID)
	if err != nil {
		return false, "Failed to collect VM stats: " + err.Error(), ""
	}

	if err := client.SubmitVMStats(VMStatsRequest{
		HostID: host.ID,
		VMs:    stats,
	}); err != nil {
		return false, "Failed to submit VM stats: " + err.Error(), ""
	}

	return true, fmt.Sprintf("Collected stats for %d VMs", len(stats)), fmt.Sprintf("vms=%d", len(stats))
}

func queryHypervisorVMs(hvType, host string, port int, username, hostID string) ([]VMSyncEntry, error) {
	switch hvType {
	case "proxmox":
		return queryProxmoxVMs(host, port, username, hostID)
	case "vmware":
		return queryVMwareVMs(host, port, username, hostID)
	case "hyperv":
		return queryHyperVVMs(host, port, username, hostID)
	default:
		return nil, fmt.Errorf("unsupported hypervisor type: %s", hvType)
	}
}

func collectVMStats(hvType, host string, port int, username, hostID string) ([]VMStatsEntry, error) {
	switch hvType {
	case "proxmox":
		return collectProxmoxStats(host, port, username, hostID)
	case "vmware":
		return collectVMwareStats(host, port, username, hostID)
	case "hyperv":
		return collectHyperVStats(host, port, username, hostID)
	default:
		return nil, fmt.Errorf("unsupported hypervisor type: %s", hvType)
	}
}

func queryProxmoxVMs(host string, port int, username, hostID string) ([]VMSyncEntry, error) {
	return nil, fmt.Errorf("Proxmox API integration not yet implemented")
}

func queryVMwareVMs(host string, port int, username, hostID string) ([]VMSyncEntry, error) {
	return nil, fmt.Errorf("VMware API integration not yet implemented")
}

func queryHyperVVMs(host string, port int, username, hostID string) ([]VMSyncEntry, error) {
	return nil, fmt.Errorf("Hyper-V API integration not yet implemented")
}

func collectProxmoxStats(host string, port int, username, hostID string) ([]VMStatsEntry, error) {
	return nil, fmt.Errorf("Proxmox stats integration not yet implemented")
}

func collectVMwareStats(host string, port int, username, hostID string) ([]VMStatsEntry, error) {
	return nil, fmt.Errorf("VMware stats integration not yet implemented")
}

func collectHyperVStats(host string, port int, username, hostID string) ([]VMStatsEntry, error) {
	return nil, fmt.Errorf("Hyper-V stats integration not yet implemented")
}
