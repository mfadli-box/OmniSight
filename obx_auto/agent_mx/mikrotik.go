package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

type RouterOSClient struct {
	host     string
	port     int
	tls      bool
	username string
	password string
}

func NewRouterOSClient(host string, port int, tls bool, username, password string) *RouterOSClient {
	return &RouterOSClient{
		host:     host,
		port:     port,
		tls:      tls,
		username: username,
		password: password,
	}
}

type LoginResult struct {
	Success bool
	Error   string
}

func (c *RouterOSClient) TestConnection() LoginResult {
	cmd := "/ip/address/print"
	_, err := c.runCmd(cmd)
	if err != nil {
		return LoginResult{Success: false, Error: err.Error()}
	}
	return LoginResult{Success: true}
}

type FirewallRule struct {
	ID           string
	Chain        string
	Action       string
	SrcAddress   string
	DstAddress   string
	Protocol     string
	SrcPort      string
	DstPort      string
	InInterface  string
	OutInterface string
	Comment      string
	Disabled     bool
}

func (c *RouterOSClient) GetFirewallRules() ([]FirewallRule, error) {
	lines, err := c.runCmd("/ip/firewall/filter/print")
	if err != nil {
		return nil, fmt.Errorf("get firewall rules: %w", err)
	}

	var rules []FirewallRule
	for _, line := range lines {
		if !strings.HasPrefix(line, "=") {
			continue
		}
		rule := c.parseFirewallLine(line)
		if rule.ID != "" {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func (c *RouterOSClient) parseFirewallLine(line string) FirewallRule {
	rule := FirewallRule{}
	fields := strings.Split(line, ",")
	for _, f := range fields {
		parts := strings.SplitN(strings.TrimSpace(f), "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case ".id":
			rule.ID = val
		case "chain":
			rule.Chain = val
		case "action":
			rule.Action = val
		case "src-address":
			rule.SrcAddress = val
		case "dst-address":
			rule.DstAddress = val
		case "protocol":
			rule.Protocol = val
		case "src-port":
			rule.SrcPort = val
		case "dst-port":
			rule.DstPort = val
		case "in-interface":
			rule.InInterface = val
		case "out-interface":
			rule.OutInterface = val
		case "comment":
			rule.Comment = val
		case "disabled":
			rule.Disabled = val == "true"
		}
	}
	return rule
}

func (c *RouterOSClient) AddFirewallRule(chain, action, comment string, disabled bool) error {
	cmd := fmt.Sprintf("/ip/firewall/filter/add=chain=%s,action=%s,comment=%s,disabled=%v",
		chain, action, comment, disabled)
	_, err := c.runCmd(cmd)
	if err != nil {
		return fmt.Errorf("add firewall rule: %w", err)
	}
	return nil
}

func (c *RouterOSClient) RemoveFirewallRule(id string) error {
	cmd := fmt.Sprintf("/ip/firewall/filter/removeumbers=%s", id)
	_, err := c.runCmd(cmd)
	if err != nil {
		return fmt.Errorf("remove firewall rule: %w", err)
	}
	return nil
}

func (c *RouterOSClient) EnableFirewallRule(id string) error {
	cmd := fmt.Sprintf("/ip/firewall/filter/set=numbers=%s,disabled=false", id)
	_, err := c.runCmd(cmd)
	if err != nil {
		return fmt.Errorf("enable firewall rule: %w", err)
	}
	return nil
}

func (c *RouterOSClient) DisableFirewallRule(id string) error {
	cmd := fmt.Sprintf("/ip/firewall/filter/set=numbers=%s,disabled=true", id)
	_, err := c.runCmd(cmd)
	if err != nil {
		return fmt.Errorf("disable firewall rule: %w", err)
	}
	return nil
}

type BackupInfo struct {
	Name      string
	Size      int
	CreatedAt string
	Password  bool
}

func (c *RouterOSClient) CreateBackup(name, password string) error {
	var cmd string
	if password != "" {
		cmd = fmt.Sprintf("/system/backup/save=name=%s,password=%s", name, password)
	} else {
		cmd = fmt.Sprintf("/system/backup/save=name=%s", name)
	}
	_, err := c.runCmd(cmd)
	if err != nil {
		return fmt.Errorf("create backup: %w", err)
	}
	return nil
}

func (c *RouterOSClient) ListBackups() ([]BackupInfo, error) {
	lines, err := c.runCmd("/system/backup/print")
	if err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}

	var backups []BackupInfo
	for _, line := range lines {
		if !strings.Contains(line, "name=") {
			continue
		}
		bi := BackupInfo{}
		fields := strings.Split(line, ",")
		for _, f := range fields {
			parts := strings.SplitN(strings.TrimSpace(f), "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			switch key {
			case "name":
				bi.Name = val
			case "size":
				fmt.Sscanf(val, "%d", &bi.Size)
			case "created":
				bi.CreatedAt = val
			case "backup-password":
				bi.Password = val != ""
			}
		}
		if bi.Name != "" {
			backups = append(backups, bi)
		}
	}
	return backups, nil
}

func (c *RouterOSClient) LoadBackup(name string) error {
	cmd := fmt.Sprintf("/system/backup/load=name=%s", name)
	_, err := c.runCmd(cmd)
	if err != nil {
		return fmt.Errorf("load backup: %w", err)
	}
	return nil
}

func (c *RouterOSClient) ExportBackup(name string) (string, error) {
	cmd := fmt.Sprintf("/export terse=no,file=%s", name)
	lines, err := c.runCmd(cmd)
	if err != nil {
		return "", fmt.Errorf("export backup: %w", err)
	}
	return strings.Join(lines, "\n"), nil
}

type InterfaceInfo struct {
	ID       string
	Name     string
	Type     string
	RXBytes  int64
	TXBytes  int64
	Running  bool
	Disabled bool
}

func (c *RouterOSClient) GetInterfaces() ([]InterfaceInfo, error) {
	lines, err := c.runCmd("/interface/print")
	if err != nil {
		return nil, fmt.Errorf("get interfaces: %w", err)
	}

	var ifaces []InterfaceInfo
	for _, line := range lines {
		if !strings.Contains(line, "name=") {
			continue
		}
		iface := InterfaceInfo{}
		fields := strings.Split(line, ",")
		for _, f := range fields {
			parts := strings.SplitN(strings.TrimSpace(f), "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			switch key {
			case ".id":
				iface.ID = val
			case "name":
				iface.Name = val
			case "type":
				iface.Type = val
			case "rx-byte":
				fmt.Sscanf(val, "%d", &iface.RXBytes)
			case "tx-byte":
				fmt.Sscanf(val, "%d", &iface.TXBytes)
			case "running":
				iface.Running = val == "true"
			case "disabled":
				iface.Disabled = val == "true"
			}
		}
		if iface.Name != "" {
			ifaces = append(ifaces, iface)
		}
	}
	return ifaces, nil
}

type SystemResource struct {
	CPUCount    int
	CPULoad     float64
	MemoryTotal int64
	MemoryUsed  int64
	DiskTotal   int64
	DiskUsed    int64
	Uptime      string
}

func (c *RouterOSClient) GetSystemResources() (*SystemResource, error) {
	lines, err := c.runCmd("/system/resource/print")
	if err != nil {
		return nil, fmt.Errorf("get system resources: %w", err)
	}

	res := &SystemResource{}
	for _, line := range lines {
		fields := strings.Split(line, ",")
		for _, f := range fields {
			parts := strings.SplitN(strings.TrimSpace(f), "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			switch key {
			case "cpu-count":
				fmt.Sscanf(val, "%d", &res.CPUCount)
			case "cpu-load":
				fmt.Sscanf(val, "%f", &res.CPULoad)
			case "total-memory":
				fmt.Sscanf(val, "%d", &res.MemoryTotal)
			case "used-memory":
				fmt.Sscanf(val, "%d", &res.MemoryUsed)
			case "total-hdd-space":
				fmt.Sscanf(val, "%d", &res.DiskTotal)
			case "used-hdd-space":
				fmt.Sscanf(val, "%d", &res.DiskUsed)
			case "uptime":
				res.Uptime = val
			}
		}
	}
	return res, nil
}

func (c *RouterOSClient) runCmd(cmd string) ([]string, error) {
	var err error
	conn, err := c.connect()
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	if err := c.login(conn); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	if err := c.sendWord(conn, cmd); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	var lines []string
	for {
		word, err := c.readWord(conn)
		if err != nil {
			return nil, fmt.Errorf("read: %w", err)
		}
		if word == "!done" || word == "!re" {
			break
		}
		if word == "!trap" {
			break
		}
		if word != "" {
			lines = append(lines, word)
		}
	}
	return lines, nil
}

func (c *RouterOSClient) connect() (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *RouterOSClient) login(conn net.Conn) error {
	c.sendWord(conn, "/login")
	reply, err := c.readWord(conn)
	if err != nil {
		return err
	}
	if reply == "!done" {
		return nil
	}
	challenge := ""
	if strings.HasPrefix(reply, "=ret=") {
		challenge = strings.TrimPrefix(reply, "=ret=")
	}

	hash := c.hashPassword(challenge)
	loginCmd := fmt.Sprintf("/login/name=%s/password=%s", c.username, hash)
	c.sendWord(conn, loginCmd)
	_, err = c.readWord(conn)
	return err
}

func (c *RouterOSClient) hashPassword(challenge string) string {
	if challenge == "" {
		return c.password
	}
	h := md5.New()
	h.Write([]byte(c.password))
	h.Write([]byte(challenge))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *RouterOSClient) sendWord(conn net.Conn, word string) error {
	w := make([]byte, 4)
	b := []byte(word)
	binary.BigEndian.PutUint32(w, uint32(len(b)))
	_, err := conn.Write(w)
	if err != nil {
		return err
	}
	_, err = conn.Write(b)
	return err
}

func (c *RouterOSClient) readWord(conn net.Conn) (string, error) {
	lenBuf := make([]byte, 4)
	_, err := conn.Read(lenBuf)
	if err != nil {
		return "", err
	}
	length := binary.BigEndian.Uint32(lenBuf)
	if length == 0 {
		return "", nil
	}
	buf := make([]byte, length)
	_, err = conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

type SyslogTarget struct {
	Name    string
	Remote  string
	SrcAddr string
	Port    int
	Enabled bool
}

func (c *RouterOSClient) AddSyslogTarget(remoteAddr string, remotePort int) error {
	cmd := fmt.Sprintf("/system/logging/action/add=name=remote.targets=remote=%s:%d", remoteAddr, remotePort)
	_, err := c.runCmd(cmd)
	if err != nil {
		return fmt.Errorf("add syslog target: %w", err)
	}
	return nil
}

func (c *RouterOSClient) RemoveSyslogTarget(remoteAddr string) error {
	targets, err := c.ListSyslogTargets()
	if err != nil {
		return fmt.Errorf("list targets: %w", err)
	}
	for _, t := range targets {
		if t.Remote == remoteAddr {
			cmd := fmt.Sprintf("/system/logging/removeumbers=%s", t.Name)
			_, err := c.runCmd(cmd)
			if err != nil {
				return fmt.Errorf("remove target: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("target not found: %s", remoteAddr)
}

func (c *RouterOSClient) ListSyslogTargets() ([]SyslogTarget, error) {
	lines, err := c.runCmd("/system/logging/print")
	if err != nil {
		return nil, fmt.Errorf("list syslog targets: %w", err)
	}

	var targets []SyslogTarget
	for _, line := range lines {
		if !strings.Contains(line, "remote=") {
			continue
		}
		st := SyslogTarget{}
		fields := strings.Split(line, ",")
		for _, f := range fields {
			parts := strings.SplitN(strings.TrimSpace(f), "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			switch key {
			case ".id":
				st.Name = val
			case "remote":
				st.Remote = val
			case "src-address":
				st.SrcAddr = val
			case "enabled":
				st.Enabled = val == "true"
			}
		}
		if st.Name != "" && st.Remote != "" {
			targets = append(targets, st)
		}
	}
	return targets, nil
}
