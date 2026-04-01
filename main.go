package main

import (
	_ "embed"
	"bytes"
	"encoding/base64"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"ascii_convert_go/describer"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var indexHTML string

// desc 是全局 Claude Vision 客户端，启动时初始化一次
var desc *describer.Client

func main() {
	// 初始化 Claude 客户端（自动读取 ANTHROPIC_API_KEY）
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		desc = describer.New("")
		log.Println("Claude Vision 已启用")
	} else {
		log.Println("未设置 ANTHROPIC_API_KEY，Vision 描述功能已禁用")
	}

	r := gin.Default()
	r.MaxMultipartMemory = 10 << 20 // 10MB 上传限制

	r.GET("/", handleIndex)
	r.POST("/convert", handleConvert)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// handleIndex 返回前端页面
func handleIndex(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, indexHTML)
}

// handleConvert 接收上传图片，返回 ASCII 字符画 PNG 及 Claude 图片描述
func handleConvert(c *gin.Context) {
	// 1. 读取上传文件到内存（需要复用给两个并行任务）
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 image 字段: " + err.Error()})
		return
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败: " + err.Error()})
		return
	}

	// 2. 解码图片
	src, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片解码失败: " + err.Error()})
		return
	}

	// 3. 解析参数
	width, err := strconv.Atoi(c.DefaultPostForm("width", "100"))
	if err != nil || width <= 0 {
		width = 100
	}
	chars := c.DefaultPostForm("chars", "")
	colorful := c.DefaultPostForm("colorful", "false") == "true"

	// 4. 并行执行：ASCII 转换 + Claude Vision 描述
	var (
		pngBytes    []byte
		convertErr  error
		description string
		descErr     error
		wg          sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		pngBytes, convertErr = ConvertToASCII(src, ConvertOptions{Width: width, CharSet: chars, Colorful: colorful})
	}()

	if desc != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mediaType := header.Header.Get("Content-Type")
			description, descErr = desc.Describe(c.Request.Context(), imgBytes, mediaType)
			if descErr != nil {
				log.Printf("Vision 描述失败（不影响主流程）: %v", descErr)
			}
		}()
	}

	wg.Wait()

	// 5. ASCII 转换失败则报错
	if convertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "转换失败: " + convertErr.Error()})
		return
	}

	// 6. 返回 JSON（包含 base64 PNG + 描述）
	c.JSON(http.StatusOK, gin.H{
		"png":         "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes),
		"description": description,
	})
}
