package main

import (
	"image"

	"ascii_convert_go/converter"
)

// 以下是对 converter 包的薄封装，保持 main 包的 API 不变，测试文件可继续使用

// DefaultCharSet 默认字符集
const DefaultCharSet = converter.DefaultCharSet

// ConvertOptions 转换参数
type ConvertOptions = converter.Options

// ConvertToASCII 将图片转换为 ASCII 字符画，返回 PNG 字节流
func ConvertToASCII(src image.Image, opt ConvertOptions) ([]byte, error) {
	return converter.Convert(src, opt)
}
