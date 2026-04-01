package converter

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// DefaultCharSet 默认字符集（从暗到亮）
const DefaultCharSet = `@#S%?*+;:,. `

// UnicodeCharSet 精细字符集，使用 Unicode 半块字符（从暗到亮）
const UnicodeCharSet = "█▓▒░▐▌▄▀+;:,. "

// basicfont.Face7x13 每个字符的像素尺寸
const charW = 7
const charH = 13

// Options 转换参数
type Options struct {
	Width    int    // 输出字符列数，<=0 时使用默认值 100
	CharSet  string // 字符集（从暗到亮），空时使用 DefaultCharSet
	Colorful bool   // 是否启用彩色渲染（字符颜色与原图对应区域颜色一致）
}

// Convert 将图片转换为 ASCII 字符画，返回 PNG 字节流
func Convert(src image.Image, opt Options) ([]byte, error) {
	if opt.Width <= 0 {
		opt.Width = 100
	}
	if opt.CharSet == "" {
		opt.CharSet = DefaultCharSet
	}
	grid := buildCharGrid(src, opt.Width, opt.CharSet)
	out := renderToImage(grid)
	return encodePNG(out)
}

// buildCharGrid 将图片映射为二维 ASCII 字符网格
func buildCharGrid(img image.Image, width int, charSet string) []string {
	bounds := img.Bounds()
	srcW := bounds.Max.X - bounds.Min.X
	srcH := bounds.Max.Y - bounds.Min.Y

	// 宽高比修正：basicfont 字符为 7×13，非正方形
	height := int(math.Round(float64(srcH) * float64(width) / float64(srcW) * float64(charW) / float64(charH)))
	if height < 1 {
		height = 1
	}

	last := len(charSet) - 1
	grid := make([]string, height)

	for row := 0; row < height; row++ {
		var sb strings.Builder
		sb.Grow(width)
		for col := 0; col < width; col++ {
			px := bounds.Min.X + col*srcW/width
			py := bounds.Min.Y + row*srcH/height
			g := ToGray(img.At(px, py))
			idx := int(g) * last / 255
			sb.WriteByte(charSet[idx])
		}
		grid[row] = sb.String()
	}
	return grid
}

// renderToImage 将字符网格渲染为黑底白字的 RGBA 图片
func renderToImage(grid []string) *image.RGBA {
	if len(grid) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	cols := len(grid[0])
	rows := len(grid)

	out := image.NewRGBA(image.Rect(0, 0, cols*charW, rows*charH))
	draw.Draw(out, out.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  out,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
	}
	for row, line := range grid {
		d.Dot = fixed.P(0, row*charH+(charH-2))
		d.DrawString(line)
	}
	return out
}

// encodePNG 将图片编码为 PNG 字节流
func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ToGray 将颜色转换为灰度值（0-255），使用人眼感知加权公式
func ToGray(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
	return uint8(gray)
}
