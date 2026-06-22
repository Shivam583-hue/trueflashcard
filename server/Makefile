# Server build tooling (used by Render)
.PHONY: build start run

## build: compile the server binary — use as Render's Build Command
build:
	go build -o bin/server ./cmd/server

## start: run the compiled binary — use as Render's Start Command
start:
	./bin/server

## run: build then start
run: build start
