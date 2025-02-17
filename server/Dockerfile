# server/Dockerfile
FROM golang:latest AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /Server

COPY . .

RUN go mod download

RUN go build -o my_server server/main.go

FROM ubuntu:latest
COPY server/Settings/config.yaml ./Settings/config.yaml
COPY server/wait-for-it.sh ./wait-for-it.sh
COPY server/templates /templates
COPY server/static /static
COPY --from=builder /Server/my_server /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		 netcat-openbsd; \
        chmod 755 wait-for-it.sh

#EXPOSE 8970
#
#CMD ["./server"]
