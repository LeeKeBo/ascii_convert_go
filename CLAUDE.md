# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# 构建
go build ./...

# 构建并输出可执行文件
go build -o ascii_convert_go .

# 运行（需提供图片路径）
go run main.go <图片路径>
go run main.go -width 80 image.jpg
go run main.go -chars " .:-=+*#%@" image.jpg

# 测试
go test ./...

# 单个测试函数
go test -run TestConvertToASCII .
```

## 架构

单文件 CLI 工具，零外部依赖，仅使用 Go 标准库。

- **入口 `main()`**：用 `flag` 解析 `-width`（默认 100）和 `-chars`（默认 `@#S%?*+;:,.`）两个参数，位置参数为图片路径
- **`convertToASCII(img, width, charSet)`**：核心转换逻辑，按等比缩放（高度 ÷2 修正字符宽高比），将每个像素灰度值线性映射到 `charSet` 中的字符
- **`toGray(color)`**：使用感知加权公式 `0.299R + 0.587G + 0.114B` 计算灰度

支持格式：JPEG、PNG（通过 `image/jpeg`、`image/png` 的 blank import 注册解码器）。

## 字符集约定

`-chars` 参数的字符顺序必须是**从暗到亮**，程序按灰度值线性索引字符集。
