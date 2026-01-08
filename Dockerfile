FROM golang:alpine AS preper
WORKDIR /usr/src/skvdmt-tgbot-msgs
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /var/log/skvdmt-tgbot-msgs
RUN mkdir -p /etc/skvdmt-tgbot-msgs
COPY ./config /etc/skvdmt-tgbot-msgs
COPY ./fonts /usr/local/share/fonts

FROM preper AS testing
ARG MODE
RUN go test --tags=unit -v ./...

FROM preper AS builder
RUN go build -v -o /usr/local/bin/skvdmt-tgbot-msgs ./cmd/main.go

FROM alpine AS release
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN mkdir -p /var/log/skvdmt-tgbot-msgs
RUN mkdir -p /etc/skvdmt-tgbot-msgs
COPY ./config /etc/skvdmt-tgbot-msgs
COPY ./fonts /usr/local/share/fonts
WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/skvdmt-tgbot-msgs ./skvdmt-tgbot-msgs
ENTRYPOINT ["skvdmt-tgbot-msgs"]
