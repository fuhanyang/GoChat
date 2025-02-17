# message-service/Dockerfile
FROM golang:latest AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app/Message

COPY . .


RUN go mod download

RUN go build -o message-service app/Message/main.go

FROM ubuntu:latest

COPY app/Message/Settings/config.yaml ./Settings/config.yaml
COPY app/Message/wait-for-it.sh ./wait-for-it.sh

COPY --from=builder /app/Message/message-service /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		 netcat-openbsd; \
        chmod 755 wait-for-it.sh

#EXPOSE 8973
#
#CMD ["./message-service","app/Message/Settings/config.yaml"]
