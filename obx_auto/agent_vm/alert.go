package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type AlertEvaluator struct {
	rules   []AlertRule
	pushers []Pusher
}

func NewAlertEvaluator() *AlertEvaluator {
	return &AlertEvaluator{
		rules:   []AlertRule{},
		pushers: []Pusher{},
	}
}

func (e *AlertEvaluator) AddRule(rule AlertRule) {
	e.rules = append(e.rules, rule)
}

func (e *AlertEvaluator) AddPusher(p Pusher) {
	e.pushers = append(e.pushers, p)
}

func (e *AlertEvaluator) Evaluate(stats HostStatsPayload) []AlertEvent {
	var events []AlertEvent

	for _, rule := range e.rules {
		if event := rule.Evaluate(stats); event != nil {
			events = append(events, *event)
		}
	}

	return events
}

func (e *AlertEvaluator) SendAlerts(events []AlertEvent, notifier Notifier) error {
	for _, event := range events {
		for _, p := range e.pushers {
			if err := p.Send(event); err != nil {
				return fmt.Errorf("send alert via %s: %w", p.Name(), err)
			}
		}
		if notifier != nil {
			notifier.Notify(event)
		}
	}
	return nil
}

type AlertRule struct {
	ID          string
	Name        string
	Condition   string
	Severity    string
	Description string
}

func (r *AlertRule) Evaluate(stats HostStatsPayload) *AlertEvent {
	switch {
	case strings.Contains(r.Condition, "cpu_high"):
		threshold := parseThreshold(r.Condition)
		if stats.CPUPercent > threshold {
			return &AlertEvent{
				RuleID:      r.ID,
				RuleName:    r.Name,
				Severity:    r.Severity,
				Message:     fmt.Sprintf("CPU usage %.1f%% exceeds threshold %.0f%%", stats.CPUPercent, threshold),
				Metrics:     map[string]float64{"cpu_percent": stats.CPUPercent, "threshold": threshold},
				TriggeredAt: time.Now().UTC().Format(time.RFC3339),
			}
		}
	case strings.Contains(r.Condition, "memory_high"):
		threshold := parseThreshold(r.Condition)
		memPercent := float64(stats.MemUsed) / float64(stats.MemTotal) * 100
		if memPercent > threshold {
			return &AlertEvent{
				RuleID:      r.ID,
				RuleName:    r.Name,
				Severity:    r.Severity,
				Message:     fmt.Sprintf("Memory usage %.1f%% exceeds threshold %.0f%%", memPercent, threshold),
				Metrics:     map[string]float64{"memory_percent": memPercent, "threshold": threshold},
				TriggeredAt: time.Now().UTC().Format(time.RFC3339),
			}
		}
	case strings.Contains(r.Condition, "disk_high"):
		threshold := parseThreshold(r.Condition)
		diskPercent := float64(stats.DiskUsed) / float64(stats.DiskTotal) * 100
		if diskPercent > threshold {
			return &AlertEvent{
				RuleID:      r.ID,
				RuleName:    r.Name,
				Severity:    r.Severity,
				Message:     fmt.Sprintf("Disk usage %.1f%% exceeds threshold %.0f%%", diskPercent, threshold),
				Metrics:     map[string]float64{"disk_percent": diskPercent, "threshold": threshold},
				TriggeredAt: time.Now().UTC().Format(time.RFC3339),
			}
		}
	case strings.Contains(r.Condition, "container_down"):
		re := regexp.MustCompile(`container_down\[(\w+)\]`)
		matches := re.FindStringSubmatch(r.Condition)
		if len(matches) >= 2 {
			containerName := matches[1]
			for _, c := range stats.Containers {
				if c.Name == containerName && c.Status != "running" {
					return &AlertEvent{
						RuleID:      r.ID,
						RuleName:    r.Name,
						Severity:    r.Severity,
						Message:     fmt.Sprintf("Container %s is %s", containerName, c.Status),
						Metrics:     map[string]float64{},
						TriggeredAt: time.Now().UTC().Format(time.RFC3339),
					}
				}
			}
		}
	}
	return nil
}

func parseThreshold(condition string) float64 {
	re := regexp.MustCompile(`>(\d+)`)
	matches := re.FindStringSubmatch(condition)
	if len(matches) >= 2 {
		var v float64
		fmt.Sscanf(matches[1], "%f", &v)
		return v
	}
	return 80
}

type AlertEvent struct {
	RuleID      string             `json:"rule_id"`
	RuleName    string             `json:"rule_name"`
	Severity    string             `json:"severity"`
	Message     string             `json:"message"`
	Metrics     map[string]float64 `json:"metrics"`
	TriggeredAt string             `json:"triggered_at"`
}

type Pusher interface {
	Name() string
	Send(event AlertEvent) error
}

type Notifier interface {
	Notify(event AlertEvent)
}

type LogNotifier struct{}

func (n *LogNotifier) Notify(event AlertEvent) {
	data, _ := json.Marshal(event)
	_ = data
}

func LoadAlertRules(path string) ([]AlertRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rules file: %w", err)
	}

	var rules []AlertRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("parse rules json: %w", err)
	}
	return rules, nil
}
