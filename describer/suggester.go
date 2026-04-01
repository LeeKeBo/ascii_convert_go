package describer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

// SuggestParams 智能推荐的转换参数
type SuggestParams struct {
	Width    int    `json:"width"`
	Chars    string `json:"chars"`
	Colorful bool   `json:"colorful"`
	Reason   string `json:"reason"`
}

// Suggest 分析图片内容，返回推荐的 ASCII 转换参数
func (c *Client) Suggest(ctx context.Context, imgBytes []byte, mediaType string) (*SuggestParams, error) {
	if mediaType == "" {
		mediaType = "image/jpeg"
	}
	b64 := base64.StdEncoding.EncodeToString(imgBytes)

	prompt := `分析这张图片，推荐最适合的 ASCII 字符画转换参数。

请只返回如下 JSON，不要有任何其他文字：
{
  "width": <40-200的整数，人像用120，风景用100，插画/文字用80>,
  "chars": <字符集字符串，默认"@#S%?*+;:,. "，精细用"█▓▒░▐▌▄▀+;:,. ">,
  "colorful": <true或false，彩色照片用true，黑白/线稿用false>,
  "reason": <一句话中文说明推荐理由，不超过20字>
}`

	msg, err := c.inner.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_6,
		MaxTokens: 256,
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					anthropic.NewImageBlockBase64(mediaType, b64),
					anthropic.NewTextBlock(prompt),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Claude API 调用失败: %w", err)
	}

	var raw string
	for _, block := range msg.Content {
		if text := block.AsText(); text.Text != "" {
			raw = text.Text
			break
		}
	}
	if raw == "" {
		return nil, fmt.Errorf("响应中没有文字内容")
	}

	// 提取 JSON（防止模型在 JSON 外附加文字）
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("响应不包含有效 JSON: %s", raw)
	}
	raw = raw[start : end+1]

	var params SuggestParams
	if err := json.Unmarshal([]byte(raw), &params); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w, raw: %s", err, raw)
	}

	// 边界保护
	if params.Width < 40 {
		params.Width = 40
	}
	if params.Width > 200 {
		params.Width = 200
	}
	if params.Chars == "" {
		params.Chars = "@#S%?*+;:,. "
	}
	return &params, nil
}
