# ── 构建阶段 ──────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 先复制 go.mod/go.sum，利用层缓存加速依赖下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译静态二进制
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o /ascii_convert_go .

# ── 运行阶段 ──────────────────────────────────────────────
FROM alpine:3.19

# 安装 CA 证书（用于 HTTPS 出向请求，如调用 Anthropic API）
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /ascii_convert_go .

# 非 root 用户运行
RUN adduser -D -H appuser && chown appuser /app/ascii_convert_go
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/ascii_convert_go"]
