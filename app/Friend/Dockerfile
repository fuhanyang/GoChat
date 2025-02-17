# friend-service/Dockerfile
FROM golang:latest AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app/Friend

COPY . .

RUN go mod download

RUN go build -o friend-service app/Friend/main.go

#EXPOSE 8974
#
#CMD ["./friend-service","app/Friend/Settings/config.yaml"]
FROM ubuntu:latest

COPY app/Friend/Settings/config.yaml /Settings/
COPY app/Friend/wait-for-it.sh  /

COPY --from=builder /app/Friend/friend-service /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		 netcat-openbsd; \
        chmod 755 wait-for-it.sh