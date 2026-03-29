package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"ascii_convert_go/converter"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"ASCII Art Converter",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// 注册工具
	tool := mcp.NewTool("convert_to_ascii",
		mcp.WithDescription("将图片转换为 ASCII 字符画，返回 base64 编码的 PNG 图片。支持 JPG、PNG、GIF 格式。"),
		mcp.WithString("image_base64",
			mcp.Required(),
			mcp.Description("base64 编码的图片内容（纯 base64，不含 data URI 前缀）"),
		),
		mcp.WithNumber("width",
			mcp.Description("输出字符宽度（40-200），默认 100。值越大图片越精细"),
		),
		mcp.WithString("chars",
			mcp.Description(`ASCII 字符集，从暗到亮排列，默认 "@#S%?*+;:,. "；白底黑字风格可用 " .:-=+*#%@"`),
		),
	)
	s.AddTool(tool, handleConvert)

	log.Println("ASCII Art Converter MCP Server 启动（stdio transport）")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("MCP Server 错误: %v", err)
	}
}

func handleConvert(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 获取必填参数
	imageBase64, err := req.RequireString("image_base64")
	if err != nil {
		return mcp.NewToolResultError("缺少 image_base64 参数"), nil
	}

	// 2. base64 解码（兼容标准和 URL-safe 两种编码）
	imgData, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		imgData, err = base64.URLEncoding.DecodeString(imageBase64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("base64 解码失败: %v", err)), nil
		}
	}

	// 3. 解码图片
	src, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("图片解码失败（支持 JPG/PNG/GIF）: %v", err)), nil
	}
	_ = format

	// 4. 读取可选参数
	width := 100
	if w := mcp.ParseInt(req, "width", 100); w > 0 {
		width = w
		if width < 40 {
			width = 40
		} else if width > 200 {
			width = 200
		}
	}
	chars := mcp.ParseString(req, "chars", converter.DefaultCharSet)

	// 5. 转换
	pngBytes, err := converter.Convert(src, converter.Options{
		Width:   width,
		CharSet: chars,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("转换失败: %v", err)), nil
	}

	// 6. 返回 base64 编码结果（data URI 格式，支持直接在 Markdown 中显示）
	result := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
	return mcp.NewToolResultText(result), nil
}
