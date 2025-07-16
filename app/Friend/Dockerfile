# friend-service/Dockerfile
FROM golang:latest AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app/friend

COPY . .

RUN go mod download

RUN go build -o friend-service app/friend/main.go

#EXPOSE 8974
#
#CMD ["./friend-service","app/friend/viper/config.yaml"]
FROM ubuntu:latest
# 定义构建参数，默认使用本地配置
ARG ENV=local

# 根据参数复制对应环境的配置（覆盖目标文件）
COPY app/friend/settings/${ENV}/config.yaml ./settings/config.yaml
COPY app/friend/wait-for-it.sh  /

COPY --from=builder /app/friend/friend-service /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		 netcat-openbsd; \
        chmod 755 wait-for-it.sh