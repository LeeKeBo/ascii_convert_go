package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// sendAlert 发送告警到钉钉/企微 Webhook
// 钉钉 Webhook 格式，企微只需改 content 字段结构
func sendAlert(title, content string) error {
	webhookURL := os.Getenv("ALERT_WEBHOOK_URL")
	if webhookURL == "" {
		// 未配置时只打印日志
		fmt.Printf("[ALERT] %s\n%s\n", title, content)
		return nil
	}

	// 钉钉 markdown 消息格式
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  fmt.Sprintf("## %s\n\n%s", title, content),
		},
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("发送告警失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Webhook 返回非200: %d", resp.StatusCode)
	}
	return nil
}
