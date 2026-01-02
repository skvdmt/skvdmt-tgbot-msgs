FROM golang:1.25.3-alpine AS builder
WORKDIR /usr/src/skvdmt-tgbot-msgs
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test --tags=unit -v ./...
RUN go build -v -o /usr/local/bin/skvdmt-tgbot-msgs ./cmd/main.go

FROM ubuntu:latest
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN apt-get update
RUN apt-get install ca-certificates -y
RUN update-ca-certificates
WORKDIR /usr/local/bin
COPY ./schema /usr/src/skvdmt-tgbot-msgs/schema
COPY --from=builder /usr/local/bin/skvdmt-tgbot-msgs ./skvdmt-tgbot-msgs
RUN mkdir -p /var/log/skvdmt-tgbot-msgs
ENTRYPOINT ["skvdmt-tgbot-msgs"]
