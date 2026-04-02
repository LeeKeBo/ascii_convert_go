# ascii_convert_go

将图片转换为 ASCII 艺术字的 CLI 工具与 Web 服务，零外部依赖，仅使用 Go 标准库。

## 功能

- 图片转 ASCII 字符画（JPEG / PNG），支持彩色模式
- 文字内容转 ASCII 艺术字
- Claude Vision 智能参数推荐
- HTTP API 服务 + 前端界面
- 微信小程序客户端
- Docker 一键部署
- GitHub Actions AI 代码审查 + 自动 CHANGELOG
- Go 监控 Agent（健康检查 + AI 分析 + 告警）

## 快速开始（CLI）

```bash
go run main.go <图片路径>
go run main.go -width 80 image.jpg
go run main.go -chars " .:-=+*#%@" image.jpg
```

## Docker 部署

```bash
cp .env.example .env
# 编辑 .env 填入 ANTHROPIC_API_KEY
docker compose up -d --build
# 访问 http://localhost
```

参见 [docs/deploy.md](docs/deploy.md) 获取完整的 ECS 部署指南。

## 微信小程序

1. 用微信开发者工具打开 `miniprogram/` 目录
2. 修改 `miniprogram/pages/index/index.js` 中的 `API_BASE` 为你的服务器地址
3. 在 `miniprogram/project.config.json` 中填入你的 AppID

功能：图片上传转换、历史记录、保存到相册、分享好友。

## 监控 Agent

健康检查 + 错误日志扫描 + Claude AI 分析 + 钉钉/企微告警。

```bash
cd agents/monitor
go build -o ascii-monitor .

# 配置环境变量
export HEALTH_URL=http://localhost:8080/health
export APP_LOG_PATH=/var/log/ascii-app.log
export ANTHROPIC_API_KEY=sk-ant-...
export ALERT_WEBHOOK_URL=https://oapi.dingtalk.com/robot/send?access_token=...

./ascii-monitor
```

**systemd 定时运行（每5分钟）：**

```bash
sudo cp agents/monitor/deploy/monitor.service /etc/systemd/system/
sudo cp agents/monitor/deploy/monitor.timer /etc/systemd/system/
# 创建 /etc/ascii-monitor.env 填入环境变量
sudo systemctl enable --now monitor.timer
```

## GitHub Actions

| Workflow | 触发时机 | 功能 |
|---|---|---|
| `ai-review.yml` | PR 创建/更新 | Claude Haiku 审查 Go 代码，评论到 PR |
| `changelog.yml` | push to main | 自动在 CHANGELOG.md 插入一行变更记录 |

需要在 GitHub Secrets 中设置 `ANTHROPIC_API_KEY`。

## HTTP API

| 接口 | 方法 | 说明 |
|---|---|---|
| `POST /convert` | multipart | 图片转 ASCII，返回 PNG |
| `POST /convert/text` | multipart | 文字转 ASCII，返回文本 |
| `POST /convert/suggest` | multipart | Claude 推荐最优参数 |
| `GET /health` | - | 健康检查 |
