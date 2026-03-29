package main

import (
	_ "embed"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var indexHTML string

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 10 << 20 // 10MB 上传限制

	r.GET("/", handleIndex)
	r.POST("/convert", handleConvert)

	r.Run(":8080")
}

// handleIndex 返回前端页面
func handleIndex(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, indexHTML)
}

// handleConvert 接收上传图片，返回 ASCII 字符画 PNG
func handleConvert(c *gin.Context) {
	// 1. 读取上传文件
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 image 字段: " + err.Error()})
		return
	}
	defer file.Close()

	// 2. 解码图片
	src, _, err := image.Decode(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片解码失败: " + err.Error()})
		return
	}

	// 3. 解析参数
	width, _ := strconv.Atoi(c.DefaultPostForm("width", "100"))
	chars := c.DefaultPostForm("chars", "")

	// 4. 转换
	pngBytes, err := ConvertToASCII(src, ConvertOptions{Width: width, CharSet: chars})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "转换失败: " + err.Error()})
		return
	}

	// 5. 返回 PNG 字节流
	c.Data(http.StatusOK, "image/png", pngBytes)
}
