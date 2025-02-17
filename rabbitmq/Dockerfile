# rabbitmq/Dockerfile
FROM golang:latest AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /rabbitmq

COPY . .
RUN go mod download

RUN go build -o mq-server rabbitmq/main.go

FROM ubuntu:latest

COPY rabbitmq/Settings/config.yaml ./Settings/config.yaml
COPY rabbitmq/wait-for-it.sh ./wait-for-it.sh

COPY --from=builder /rabbitmq/mq-server /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		 netcat-openbsd; \
        chmod 755 wait-for-it.sh




#EXPOSE 8971
#
#CMD ["./mq-server","rabbitmq/Settings/config.yaml"]
