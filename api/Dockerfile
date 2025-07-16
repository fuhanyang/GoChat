# api/Dockerfile
FROM golang:latest AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /api

COPY . .

RUN go mod download

RUN go build -o my_api api/main.go

FROM ubuntu:latest
# 定义构建参数，默认使用本地配置
ARG ENV=local

# 根据参数复制对应环境的配置（覆盖目标文件）
COPY api/settings/${ENV}/config.yaml ./settings/config.yaml
COPY api/wait-for-it.sh ./wait-for-it.sh
COPY api/templates /templates
COPY api/static /static
COPY --from=builder /api/my_api /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		 netcat-openbsd; \
        chmod 755 wait-for-it.sh

#EXPOSE 8970
#
#CMD ["./api"]
