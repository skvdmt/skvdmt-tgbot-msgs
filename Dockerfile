FROM golang:alpine AS preper
ARG NAME
WORKDIR /usr/src/${NAME}
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /var/log/${NAME}
RUN mkdir -p /etc/${NAME}
COPY ./config /etc/${NAME}
COPY ./fonts /usr/local/share/fonts

FROM preper AS testing
ARG MODE
RUN go test --tags=unit -v ./...

FROM preper AS builder
RUN go build -v -o /usr/local/bin/${NAME} ./cmd/main.go

FROM alpine AS release
ARG NAME
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN mkdir -p /var/log/${NAME}
RUN mkdir -p /etc/${NAME}
COPY ./config /etc/${NAME}
COPY ./fonts /usr/local/share/fonts
WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/${NAME} ./${NAME}
ENTRYPOINT ["skvdmt-tgbot-msgs"]
