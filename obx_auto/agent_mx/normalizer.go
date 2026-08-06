package main

import (
	"regexp"
	"strings"
)

type SyslogNormalizer struct {
	priRe       *regexp.Regexp
	timestampRe *regexp.Regexp
	hostRe      *regexp.Regexp
	appRe       *regexp.Regexp
	pidRe       *regexp.Regexp
	msgRe       *regexp.Regexp
}

func NewSyslogNormalizer() *SyslogNormalizer {
	return &SyslogNormalizer{
		priRe:       regexp.MustCompile(`^<(\d+)>`),
		timestampRe: regexp.MustCompile(`^(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}|\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})`),
		hostRe:      regexp.MustCompile(`^(\S+)`),
		appRe:       regexp.MustCompile(`(\S+?)(\[\d+\])?:`),
		pidRe:       regexp.MustCompile(`\[(\d+)\]`),
		msgRe:       regexp.MustCompile(`:\s*(.*)$`),
	}
}

func (n *SyslogNormalizer) Normalize(raw string, sourceIP string) *LogEntry {
	if raw == "" {
		return nil
	}

	entry := &LogEntry{
		Topics:   "system",
		Severity: "info",
		SourceIP: sourceIP,
		Message:  raw,
	}

	if strings.HasPrefix(raw, "<") {
		n.parsePriority(raw, entry)
		raw = n.stripPriority(raw)
	}

	n.parseTimestamp(raw, entry)
	n.parseHost(raw, entry)
	n.parseMessage(raw, entry)

	return entry
}

func (n *SyslogNormalizer) parsePriority(raw string, entry *LogEntry) {
	matches := n.priRe.FindStringSubmatch(raw)
	if len(matches) < 2 {
		return
	}

	pri := 0
	for _, c := range matches[1] {
		pri = pri*10 + int(c-'0')
	}

	facility := pri >> 3
	severity := pri & 7

	entry.Severity = severityToLevel(severity)
	entry.Topics = facilityToTopic(facility)
}

func (n *SyslogNormalizer) stripPriority(raw string) string {
	return n.priRe.ReplaceAllString(raw, "")
}

func (n *SyslogNormalizer) parseTimestamp(raw string, entry *LogEntry) {
	matches := n.timestampRe.FindStringSubmatch(raw)
	if len(matches) >= 2 {
		entry.Time = matches[1]
	}
}

func (n *SyslogNormalizer) parseHost(raw string, entry *LogEntry) {
	rest := strings.TrimLeft(raw, " ")
	rest = n.timestampRe.ReplaceAllString(rest, "")
	rest = strings.TrimLeft(rest, " ")

	parts := strings.SplitN(rest, " ", 2)
	if len(parts) >= 1 {
		host := parts[0]
		host = n.appRe.ReplaceAllString(host, "")
		if host != "" && !strings.Contains(host, "=") {
			entry.SourceIP = host
		}
	}
}

func (n *SyslogNormalizer) parseMessage(raw string, entry *LogEntry) {
	matches := n.msgRe.FindStringSubmatch(raw)
	if len(matches) >= 2 {
		msg := matches[1]
		msg = strings.TrimSpace(msg)
		entry.Message = msg

		if contains(wordlistFirewall, msg) {
			entry.Topics = "firewall"
		} else if contains(wordlistAuth, msg) {
			entry.Topics = "auth"
		} else if contains(wordlistSystem, msg) {
			entry.Topics = "system"
		} else if contains(wordlistWarning, msg) {
			entry.Severity = "warning"
		} else if contains(wordlistError, msg) {
			entry.Severity = "error"
		}
	}
}

var wordlistFirewall = []string{
	"firewall", "filter", "nat", "mangle", "forward", "input", "output",
	"drop", "reject", "accept", "src-nat", "dst-nat",
}

var wordlistAuth = []string{
	"login", "logout", "ssh", "telnet", "winbox", "password", "failed", "invalid",
}

var wordlistSystem = []string{
	"system", " reboot", "shutdown", "startup", "interface", "bridge",
}

var wordlistWarning = []string{
	"warning", "warn", "high cpu", "high memory", "disk full",
}

var wordlistError = []string{
	"error", "fail", "critical", "crash", "panic",
}

func contains(list []string, s string) bool {
	lower := strings.ToLower(s)
	for _, w := range list {
		if strings.Contains(lower, w) {
			return true
		}
	}
	return false
}

func severityToLevel(severity int) string {
	switch severity {
	case 0:
		return "emergency"
	case 1:
		return "alert"
	case 2:
		return "critical"
	case 3:
		return "error"
	case 4:
		return "warning"
	case 5:
		return "notice"
	case 6:
		return "info"
	case 7:
		return "debug"
	default:
		return "info"
	}
}

func facilityToTopic(facility int) string {
	switch facility {
	case 0:
		return "kern"
	case 1:
		return "user"
	case 2:
		return "mail"
	case 3:
		return "daemon"
	case 4:
		return "auth"
	case 5:
		return "syslog"
	case 6:
		return "lpr"
	case 7:
		return "news"
	case 8:
		return "uucp"
	case 9:
		return "cron"
	case 10:
		return "local0"
	case 11:
		return "local1"
	case 12:
		return "local2"
	case 13:
		return "local3"
	case 14:
		return "local4"
	case 15:
		return "local5"
	case 16:
		return "local6"
	case 17:
		return "local7"
	default:
		return "system"
	}
}
