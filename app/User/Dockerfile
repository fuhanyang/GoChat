# user-service/Dockerfile
FROM golang:latest AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app/User

COPY . .


RUN go mod download

RUN go build -o user-service app/User/main.go

FROM ubuntu:latest

COPY app/User/Settings/config.yaml ./Settings/config.yaml
COPY app/User/wait-for-it.sh ./wait-for-it.sh

COPY --from=builder /app/User/user-service /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		 netcat-openbsd; \
        chmod 755 wait-for-it.sh

#EXPOSE 8972
#
#CMD ["./user-service","app/User/Settings/config.yaml"]
