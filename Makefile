go-build:
	go build -v -o ./build/skvdmt-tgbot-msgs ./cmd/main.go

docker-build:
	docker build -t skvdmt/skvdmt-tgbot-msgs:latest .

docker-restart:
	make go-build && docker container restart skvdmt-tgbot-msgs
