restart:
	make go-build && docker container restart skvdmt-tgbot-msgs

go-build:
	go build -v -o ./build/skvdmt-tgbot-msgs ./cmd/main.go

docker-build:
	docker build -t skvdmt-tgbot-msgs .
