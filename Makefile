.PHONY: build run

build:
	npm --prefix web run build
	go build -o bin/the-chat-server ./cmd/the-chat-server

run: build
	./bin/the-chat-server
