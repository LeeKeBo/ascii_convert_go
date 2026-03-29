package main

import (
	"image"
	"image/color"
	"testing"

	"ascii_convert_go/converter"
)

// newSolidImage 创建纯色测试图片
func newSolidImage(w, h int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// ── ToGray 测试 ──

func TestToGray_Black(t *testing.T) {
	if got := converter.ToGray(color.Black); got != 0 {
		t.Errorf("纯黑应为 0，got %d", got)
	}
}

func TestToGray_White(t *testing.T) {
	if got := converter.ToGray(color.White); got != 255 {
		t.Errorf("纯白应为 255，got %d", got)
	}
}

func TestToGray_Gray(t *testing.T) {
	c := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	got := converter.ToGray(c)
	if got < 126 || got > 130 {
		t.Errorf("纯灰期望约 128，got %d", got)
	}
}

func TestToGray_GreenBrighter(t *testing.T) {
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	if converter.ToGray(green) <= converter.ToGray(red) {
		t.Error("纯绿感知亮度应高于纯红")
	}
}

// ── ConvertToASCII 测试 ──

func TestConvertToASCII_Width(t *testing.T) {
	img := newSolidImage(200, 100, color.White)
	// 通过 ConvertToASCII 验证整体流程
	data, err := ConvertToASCII(img, ConvertOptions{Width: 40})
	if err != nil {
		t.Fatalf("ConvertToASCII 失败: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("返回空字节流")
	}
}

func TestConvertToASCII_ReturnsPNG(t *testing.T) {
	img := newSolidImage(100, 80, color.Gray{Y: 128})
	data, err := ConvertToASCII(img, ConvertOptions{Width: 40})
	if err != nil {
		t.Fatalf("ConvertToASCII 返回错误: %v", err)
	}
	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i, b := range pngMagic {
		if data[i] != b {
			t.Fatalf("输出不是合法 PNG，第 %d 字节期望 %02x，实际 %02x", i, b, data[i])
		}
	}
}

func TestConvertToASCII_Defaults(t *testing.T) {
	img := newSolidImage(50, 50, color.White)
	_, err := ConvertToASCII(img, ConvertOptions{})
	if err != nil {
		t.Fatalf("默认参数不应报错: %v", err)
	}
}
