package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type FCMClient struct {
	serverKey string
	token     string
}

func NewFCMClient(serverKey, token string) *FCMClient {
	return &FCMClient{
		serverKey: serverKey,
		token:     token,
	}
}

func (f *FCMClient) Send(notification FCMNotification) error {
	payload := map[string]interface{}{
		"to": f.token,
		"notification": map[string]string{
			"title": notification.Title,
			"body":  notification.Body,
		},
		"data": notification.Data,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal fcm payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create fcm request: %w", err)
	}

	req.Header.Set("Authorization", "key="+f.serverKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send fcm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("fcm error: status %d", resp.StatusCode)
	}
	return nil
}

type FCMNotification struct {
	Title string
	Body  string
	Data  map[string]string
}

type APNsClient struct {
	teamID   string
	keyID    string
	bundleID string
	token    string
}

func NewAPNsClient(teamID, keyID, bundleID, token string) *APNsClient {
	return &APNsClient{
		teamID:   teamID,
		keyID:    keyID,
		bundleID: bundleID,
		token:    token,
	}
}

func (a *APNsClient) Send(notification APNsNotification) error {
	return fmt.Errorf("apns not yet implemented")
}

type APNsNotification struct {
	Title string
	Body  string
	Data  map[string]string
}

type WebPushPusher struct {
	endpoint string
	key      string
	auth     string
}

func NewWebPushPusher(endpoint, key, auth string) *WebPushPusher {
	return &WebPushPusher{
		endpoint: endpoint,
		key:      key,
		auth:     auth,
	}
}

func (w *WebPushPusher) Name() string {
	return "webpush"
}

func (w *WebPushPusher) Send(event AlertEvent) error {
	return fmt.Errorf("webpush not yet implemented")
}

type SlackPusher struct {
	webhookURL string
}

func NewSlackPusher(webhookURL string) *SlackPusher {
	return &SlackPusher{webhookURL: webhookURL}
}

func (s *SlackPusher) Name() string {
	return "slack"
}

func (s *SlackPusher) Send(event AlertEvent) error {
	if s.webhookURL == "" {
		return fmt.Errorf("slack webhook url not configured")
	}

	attachment := map[string]interface{}{
		"color": severityColor(event.Severity),
		"title": event.RuleName,
		"text":  event.Message,
		"fields": []map[string]string{
			{"title": "Severity", "value": event.Severity, "short": "true"},
			{"title": "Time", "value": event.TriggeredAt, "short": "true"},
		},
	}

	payload := map[string]interface{}{
		"attachments": []map[string]interface{}{attachment},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal slack payload: %w", err)
	}

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create slack request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send slack: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("slack error: status %d", resp.StatusCode)
	}
	return nil
}

func severityColor(severity string) string {
	switch severity {
	case "critical":
		return "#FF0000"
	case "high":
		return "#FFA500"
	case "medium":
		return "#FFFF00"
	case "low":
		return "#00FF00"
	default:
		return "#808080"
	}
}
