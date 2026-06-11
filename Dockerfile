# ====== 构建阶段 ======
FROM golang:1.25-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 先复制依赖文件，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bluebell .

# ====== 运行阶段 ======
FROM alpine:3.19

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 创建非 root 用户运行应用
RUN adduser -D -h /app appuser

WORKDIR /app

# 从构建阶段复制编译产物
COPY --from=builder /app/bluebell .

# 复制配置文件和前端静态文件
COPY --from=builder /app/config.docker.yaml ./config.yaml
COPY --from=builder /app/bluebell_frontend/dist ./bluebell_frontend/dist

# 创建日志目录
RUN mkdir -p logs && chown -R appuser:appuser /app

USER appuser

EXPOSE 8084

CMD ["./bluebell"]
