package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// 默认 ASCII 字符集（从暗到亮）
const defaultChars = `@#S%?*+;:,.`

func main() {
	// 定义 flag 参数
	width := flag.Int("width", 100, "输出宽度（字符数）")
	chars := flag.String("chars", defaultChars, "自定义 ASCII 字符集（从暗到亮排列）")
	flag.Usage = func() {
		fmt.Println("用法: ascii_convert_go [选项] <图片路径>")
		fmt.Println("\n选项:")
		flag.PrintDefaults()
		fmt.Println("\n示例:")
		fmt.Println("  ascii_convert_go image.jpg")
		fmt.Println("  ascii_convert_go -width 80 image.jpg")
		fmt.Println(`  ascii_convert_go -chars " .:-=+*#%@" image.jpg`)
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	imagePath := flag.Arg(0)

	// 打开图片
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Printf("无法打开图片: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("无法解码图片: %v\n", err)
		os.Exit(1)
	}

	// 转换为 ASCII
	ascii := convertToASCII(img, *width, *chars)
	fmt.Println(ascii)
}

// convertToASCII 将图片转换为 ASCII 字符画
// width 参数控制输出宽度（字符数），charSet 为自定义字符集
func convertToASCII(img image.Image, width int, charSet string) string {
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y

	// 计算等比例的高度（字符高宽比约为 2:1，所以高度减半）
	height := int(math.Round(float64(imgHeight) * float64(width) / float64(imgWidth) / 2))

	result := ""

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 映射到原图坐标
			srcX := bounds.Min.X + x*imgWidth/width
			srcY := bounds.Min.Y + y*imgHeight/height

			// 获取像素颜色并转为灰度
			pixel := img.At(srcX, srcY)
			gray := toGray(pixel)

			// 映射到 ASCII 字符
			charIndex := int(gray) * (len(charSet) - 1) / 255
			result += string(charSet[charIndex])
		}
		result += "\n"
	}

	return result
}

// toGray 将颜色转换为灰度值 (0-255)
func toGray(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	// 使用人眼感知加权公式
	gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
	return uint8(gray)
}
