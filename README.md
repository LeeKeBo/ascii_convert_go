# ascii_convert_go

将图片转换为 ASCII 艺术字的 CLI 工具与 Web 服务，零外部依赖，仅使用 Go 标准库。

## 功能

- 图片转 ASCII 字符画（JPEG / PNG）
- 支持彩色模式（ANSI 转义码）
- 文字内容转 ASCII 艺术字
- Claude Vision 智能参数推荐
- HTTP API 服务

## 快速开始（CLI）

```bash
go run main.go <图片路径>
go run main.go -width 80 image.jpg
go run main.go -chars " .:-=+*#%@" image.jpg
```

## Docker 部署

参见 [docs/deploy.md](docs/deploy.md) 获取完整的 ECS 部署指南。

### 快速启动（本地）

```bash
cp .env.example .env
# 编辑 .env 填入 ANTHROPIC_API_KEY
docker compose up -d --build
# 访问 http://localhost
```
