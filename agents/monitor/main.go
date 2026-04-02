package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	healthURL := os.Getenv("HEALTH_URL")
	if healthURL == "" {
		healthURL = "http://localhost:8080/health"
	}

	db, err := openDB()
	if err != nil {
		log.Fatalf("无法打开 SQLite: %v", err)
	}
	defer db.Close()

	log.Printf("开始健康检查: %s", healthURL)

	// 1. 健康检查
	healthy, statusDetail := checkHealth(healthURL)

	// 2. 读取应用日志（最近50行）
	logPath := os.Getenv("APP_LOG_PATH") // 如 /var/log/ascii-app.log
	var errorLines []string
	if logPath != "" {
		errorLines = readErrorLines(logPath, 50)
	}

	// 3. 如果服务不健康或有错误日志，调 Claude 分析
	if !healthy || len(errorLines) > 0 {
		summary := buildSummary(healthy, statusDetail, errorLines)

		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		var analysis string
		if apiKey != "" {
			analysis, err = analyzeWithClaude(summary)
			if err != nil {
				analysis = fmt.Sprintf("（Claude 分析失败: %v）", err)
			}
		} else {
			analysis = "（未配置 ANTHROPIC_API_KEY，跳过 AI 分析）"
		}

		alertContent := fmt.Sprintf("**服务状态：** %s\n\n**错误摘要：** %s\n\n**AI 分析：** %s",
			statusDetail, summary, analysis)

		isNew, _ := isNewAlert(db, summary)
		if isNew {
			title := "⚠️ ASCII 转换服务告警"
			if !healthy {
				title = "🔴 ASCII 转换服务宕机"
			}
			if err := sendAlert(title, alertContent); err != nil {
				log.Printf("发送告警失败: %v", err)
			} else {
				log.Printf("告警已发送: %s", title)
			}
		} else {
			log.Printf("告警已去重（24小时内已通知过）")
		}
	} else {
		log.Printf("服务健康，无告警")
	}
}

func checkHealth(url string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Sprintf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	return true, string(body)
}

func readErrorLines(logPath string, maxLines int) []string {
	data, err := os.ReadFile(logPath)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")

	// 取最后 maxLines 行中包含 ERROR/WARN/panic 的行
	var result []string
	start := len(lines) - maxLines
	if start < 0 {
		start = 0
	}
	for _, line := range lines[start:] {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "error") || strings.Contains(lower, "panic") || strings.Contains(lower, "fatal") {
			result = append(result, line)
		}
	}
	return result
}

func buildSummary(healthy bool, statusDetail string, errorLines []string) string {
	var sb strings.Builder
	if !healthy {
		sb.WriteString(fmt.Sprintf("服务不可用: %s\n", statusDetail))
	}
	if len(errorLines) > 0 {
		sb.WriteString(fmt.Sprintf("发现 %d 条错误日志:\n", len(errorLines)))
		for i, line := range errorLines {
			if i >= 10 {
				sb.WriteString(fmt.Sprintf("... 还有 %d 条\n", len(errorLines)-i))
				break
			}
			sb.WriteString(line + "\n")
		}
	}
	return sb.String()
}

func analyzeWithClaude(summary string) (string, error) {
	client := anthropic.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prompt := fmt.Sprintf(`你是一位 Go 服务运维专家。以下是一个 ASCII 图片转换服务的异常摘要，请用中文分析：
1. 判断是否是已知常见问题（如内存溢出、API Key 失效、图片格式不支持等）
2. 给出可能原因（1-2条）
3. 给出处理建议（1-2条，具体可操作）

字数限制100字以内。

异常摘要：
%s`, summary)

	msg, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
		MaxTokens: anthropic.F(int64(256)),
		Messages:  anthropic.F([]anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(prompt))}),
	})
	if err != nil {
		return "", err
	}

	for _, block := range msg.Content {
		if block.Type == "text" && block.Text != "" {
			return block.Text, nil
		}
	}
	return "", fmt.Errorf("响应无文字内容")
}
