
# user-service/Dockerfile
FROM golang:latest AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app/user

COPY . .


RUN go mod download

RUN go build -o user-service app/user/main.go

FROM ubuntu:latest
# 定义构建参数，默认使用本地配置
ARG ENV=local

# 根据参数复制对应环境的配置（覆盖目标文件）
COPY app/user/settings/${ENV}/config.yaml ./settings/config.yaml
COPY app/user/wait-for-it.sh ./wait-for-it.sh

COPY --from=builder /app/user/user-service /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		 netcat-openbsd; \
        chmod 755 wait-for-it.sh

#EXPOSE 8972
#
#CMD ["./user-service","app/user/viper/config.yaml"]
