# 构建阶段
FROM golang:1.23.4-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o main .

# 运行阶段
FROM alpine:latest

# 安装基本工具和时区数据
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"] 