# Server build tooling
.PHONY: build build-linux start run

## build: compile the server binary for the local platform
build:
	go build -o bin/server ./cmd/server

## build-linux: cross-compile a static linux binary for a VPS (set GOARCH=arm64 if your VPS is ARM)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=$${GOARCH:-amd64} go build -ldflags "-s -w" -o bin/server-linux ./cmd/server

## start: run the compiled binary
start:
	./bin/server

## run: build then start
run: build start
