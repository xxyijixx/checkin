# 使用官方 Golang 镜像作为构建环境
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 设置 Go 代理环境变量
ENV GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 构建主程序
RUN go build -o main main.go

# 构建数据库迁移工具
RUN go build -o migrate ./query/migrate/migrate.go

# # 运行时镜像
# FROM alpine:latest

# # 设置工作目录
# WORKDIR /app

# # 复制构建好的二进制文件
# COPY --from=builder /app/main .
# COPY --from=builder /app/migrate .


# RUN chmod +x /app/main
# RUN chmod +x /app/migrate

# 指定默认命令
CMD ["sh", "-c", "./migrate && ./main"]