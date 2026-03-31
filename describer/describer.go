// Package describer 使用 Claude Vision API 对图片生成一句话描述
package describer

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Client 封装 Claude Vision API 调用
type Client struct {
	inner anthropic.Client
}

// New 创建 Client，apiKey 为空时自动读取 ANTHROPIC_API_KEY 环境变量
func New(apiKey string) *Client {
	opts := []option.RequestOption{}
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}
	return &Client{inner: anthropic.NewClient(opts...)}
}

// Describe 对原始图片字节（JPEG/PNG）调用 Claude Vision，返回一句话描述
func (c *Client) Describe(ctx context.Context, imgBytes []byte, mediaType string) (string, error) {
	if mediaType == "" {
		mediaType = "image/jpeg"
	}

	b64 := base64.StdEncoding.EncodeToString(imgBytes)

	msg, err := c.inner.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_6,
		MaxTokens: 256,
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					// Vision：传入 base64 图片
					anthropic.NewImageBlockBase64(mediaType, b64),
					// 文字指令
					anthropic.NewTextBlock("用一句话描述这张图片的主要内容，中文回答，不超过 30 个字。"),
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("Claude API 调用失败: %w", err)
	}

	// 取第一个 text block
	for _, block := range msg.Content {
		if text := block.AsText(); text.Text != "" {
			return text.Text, nil
		}
	}
	return "", fmt.Errorf("响应中没有文字内容")
}
