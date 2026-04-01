package converter

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// DefaultCharSet 默认字符集（从暗到亮）
const DefaultCharSet = `@#S%?*+;:,. `

// UnicodeCharSet 精细字符集，使用 Unicode 块状图形字符（需按 rune 索引，Task 2 后生效）（从暗到亮）
const UnicodeCharSet = "█▓▒░▐▌▄▀+;:,. "

// basicfont.Face7x13 每个字符的像素尺寸
const charW = 7
const charH = 13

// charCell 存储一个字符格的 ASCII 字符和对应的原图平均颜色
type charCell struct {
	ch    rune
	color color.RGBA
}

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
	out := renderToImage(grid, opt.Colorful)
	return encodePNG(out)
}

// buildCharGrid 将图片映射为二维字符格网格（区域均值采样）
func buildCharGrid(img image.Image, width int, charSet string) [][]charCell {
	bounds := img.Bounds()
	srcW := bounds.Max.X - bounds.Min.X
	srcH := bounds.Max.Y - bounds.Min.Y

	// 宽高比修正：basicfont 字符为 7×13，非正方形
	height := int(math.Round(float64(srcH) * float64(width) / float64(srcW) * float64(charW) / float64(charH)))
	if height < 1 {
		height = 1
	}

	runes := []rune(charSet)
	last := len(runes) - 1
	grid := make([][]charCell, height)

	for row := 0; row < height; row++ {
		grid[row] = make([]charCell, width)
		for col := 0; col < width; col++ {
			// 计算该格对应的原图区域
			x0 := bounds.Min.X + col*srcW/width
			x1 := bounds.Min.X + (col+1)*srcW/width
			y0 := bounds.Min.Y + row*srcH/height
			y1 := bounds.Min.Y + (row+1)*srcH/height
			if x1 <= x0 {
				x1 = x0 + 1
			}
			if y1 <= y0 {
				y1 = y0 + 1
			}

			// 对区域内所有像素求 R/G/B 均值
			var rSum, gSum, bSum, count uint64
			for py := y0; py < y1; py++ {
				for px := x0; px < x1; px++ {
					r, g, b, _ := img.At(px, py).RGBA()
					rSum += uint64(r >> 8)
					gSum += uint64(g >> 8)
					bSum += uint64(b >> 8)
					count++
				}
			}
			avgR := uint8(rSum / count)
			avgG := uint8(gSum / count)
			avgB := uint8(bSum / count)

			// 由平均色计算灰度，映射到字符集
			gray := ToGray(color.RGBA{R: avgR, G: avgG, B: avgB, A: 255})
			idx := int(gray) * last / 255

			grid[row][col] = charCell{
				ch:    runes[idx],
				color: color.RGBA{R: avgR, G: avgG, B: avgB, A: 255},
			}
		}
	}
	return grid
}

// renderToImage 将字符网格渲染为图片
// colorful=false：黑底白字；colorful=true：黑底彩字（彩色在 Task 3 实现）
func renderToImage(grid [][]charCell, colorful bool) *image.RGBA {
	if len(grid) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	cols := len(grid[0])
	rows := len(grid)

	out := image.NewRGBA(image.Rect(0, 0, cols*charW, rows*charH))
	draw.Draw(out, out.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  out,
		Face: basicfont.Face7x13,
	}
	for row, line := range grid {
		for col, cell := range line {
			// 暂时全部用白色（彩色在 Task 3 实现）
			d.Src = image.NewUniform(color.White)
			d.Dot = fixed.P(col*charW, row*charH+(charH-2))
			d.DrawString(string(cell.ch))
		}
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
