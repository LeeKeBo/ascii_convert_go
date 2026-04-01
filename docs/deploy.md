# ECS 部署指南

## 前提条件

- 阿里云 ECS 实例（Ubuntu 22.04 推荐），公网 IP 已绑定域名
- 安装 Docker（>= 24）和 docker compose plugin
- 安装 curl、git

## 服务器首次配置步骤

1. SSH 到服务器

2. 克隆项目：
   ```bash
   git clone <repo_url> /opt/ascii_convert_go
   ```

3. 进入目录：
   ```bash
   cd /opt/ascii_convert_go
   ```

4. 复制 `.env.example` → `.env`，填入 `ANTHROPIC_API_KEY`：
   ```bash
   cp .env.example .env
   # 编辑 .env，填写真实的 ANTHROPIC_API_KEY
   nano .env
   ```

5. 创建 Let's Encrypt 目录：
   ```bash
   mkdir -p /etc/letsencrypt/live/<your-domain>
   mkdir -p /var/lib/letsencrypt
   ```

6. 申请证书（以 certbot 为例）：
   ```bash
   apt install certbot
   certbot certonly --standalone -d <your-domain>
   ```

7. 首次启动：
   ```bash
   VERSION=v1.0.0 docker compose up -d --build
   ```

8. 验证：
   ```bash
   curl -sf http://localhost/health
   ```

## GitHub Actions Secrets 配置

在 GitHub 仓库 **Settings → Secrets and variables → Actions** 中添加：

| Secret 名称 | 说明 |
|-------------|------|
| `ECS_HOST` | ECS 公网 IP |
| `ECS_USER` | SSH 用户名（如 `ubuntu`） |
| `ECS_SSH_KEY` | 私钥内容（`cat ~/.ssh/id_rsa`） |
| `ECS_DEPLOY_PATH` | 服务器上的项目路径，如 `/opt/ascii_convert_go` |

## 手动部署命令

```bash
cd /opt/ascii_convert_go
git pull origin main
VERSION=$(git describe --tags --always 2>/dev/null || echo dev)
VERSION=${VERSION} docker compose up -d --build --remove-orphans
```

## 验证部署

```bash
# 健康检查
curl -sf http://localhost/health

# 查看日志
docker compose logs -f app

# 查看容器状态
docker compose ps
```

## 证书续签

```bash
# 手动续签（certbot 已安装）
certbot renew
docker compose restart nginx
```

## 常见问题

- **health check 失败**：等待 nginx 启动（约 15s），检查 `docker compose ps`
- **证书路径错误**：确保 `/etc/letsencrypt/live/<domain>/` 下有 `fullchain.pem` 和 `privkey.pem`
- **端口冲突**：确保宿主机 80/443 未被占用
