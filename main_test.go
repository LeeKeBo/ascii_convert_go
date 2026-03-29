package main

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

// --- toGray 测试 ---

func TestToGray_Black(t *testing.T) {
	c := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	if got := toGray(c); got != 0 {
		t.Errorf("纯黑应为 0，got %d", got)
	}
}

func TestToGray_White(t *testing.T) {
	c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	if got := toGray(c); got != 255 {
		t.Errorf("纯白应为 255，got %d", got)
	}
}

func TestToGray_WeightedFormula(t *testing.T) {
	// 验证感知加权公式 0.299R + 0.587G + 0.114B
	c := color.RGBA{R: 100, G: 150, B: 200, A: 255}
	expected := uint8(140) // 0.299*100 + 0.587*150 + 0.114*200 ≈ 140
	got := toGray(c)
	// 允许 ±1 的浮点误差
	if got < expected-1 || got > expected+1 {
		t.Errorf("期望灰度值约 %d，got %d", expected, got)
	}
}

func TestToGray_GreenChannel(t *testing.T) {
	// 绿色权重最高（0.587），纯绿应比纯红亮
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	if toGray(green) <= toGray(red) {
		t.Error("纯绿的感知亮度应高于纯红")
	}
}

// --- convertToASCII 测试 ---

// newSolidImage 创建一张纯色图片，方便测试
func newSolidImage(w, h int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestConvertToASCII_OutputWidth(t *testing.T) {
	img := newSolidImage(200, 100, color.White)
	result := convertToASCII(img, 40, defaultChars)
	for _, line := range strings.Split(strings.TrimRight(result, "\n"), "\n") {
		if len(line) != 40 {
			t.Errorf("期望每行 40 个字符，got %d", len(line))
			break
		}
	}
}

func TestConvertToASCII_AllBlack(t *testing.T) {
	// 纯黑图片 → 灰度为 0 → 映射到字符集第一个字符（最暗）
	img := newSolidImage(100, 50, color.Black)
	result := convertToASCII(img, 20, defaultChars)
	for _, ch := range strings.ReplaceAll(result, "\n", "") {
		if ch != rune(defaultChars[0]) {
			t.Errorf("纯黑图片应全为 '%c'，got '%c'", defaultChars[0], ch)
			break
		}
	}
}

func TestConvertToASCII_AllWhite(t *testing.T) {
	// 纯白图片 → 灰度为 255 → 映射到字符集最后一个字符（最亮）
	img := newSolidImage(100, 50, color.White)
	result := convertToASCII(img, 20, defaultChars)
	lastChar := rune(defaultChars[len(defaultChars)-1])
	for _, ch := range strings.ReplaceAll(result, "\n", "") {
		if ch != lastChar {
			t.Errorf("纯白图片应全为 '%c'，got '%c'", lastChar, ch)
			break
		}
	}
}

func TestConvertToASCII_CustomCharSet(t *testing.T) {
	img := newSolidImage(100, 50, color.Black)
	customChars := "AB"
	result := convertToASCII(img, 20, customChars)
	for _, ch := range strings.ReplaceAll(result, "\n", "") {
		if ch != 'A' && ch != 'B' {
			t.Errorf("只应出现自定义字符集中的字符，got '%c'", ch)
			break
		}
	}
}

func TestConvertToASCII_AspectRatio(t *testing.T) {
	// 正方形图片 → 输出高度应约为宽度的一半（字符宽高比修正）
	img := newSolidImage(100, 100, color.White)
	result := convertToASCII(img, 40, defaultChars)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	height := len(lines)
	// 高度应约为宽度的一半，允许 ±2 误差
	if height < 18 || height > 22 {
		t.Errorf("正方形图片输出高度应约为宽度的一半（~20），got %d", height)
	}
}
