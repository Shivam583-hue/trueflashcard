# Flashcard app — dev tooling
SHELL := /bin/bash
.DEFAULT_GOAL := help
.PHONY: run backend frontend install build test migrate proto help

## run: start backend and frontend together (Ctrl-C stops both)
run:
	@echo "→ backend   http://localhost:8080 (Connect)  :50051 (gRPC)"
	@echo "→ frontend  http://localhost:3000"
	@trap 'kill 0' INT TERM EXIT; \
	$(MAKE) --no-print-directory backend & \
	$(MAKE) --no-print-directory frontend & \
	wait

## backend: run the Go server with server/.env loaded
backend:
	@cd server && \
	if [ -f .env ]; then \
	  while IFS= read -r line; do \
	    case "$$line" in \#*|'') ;; *=*) export "$$line" ;; esac; \
	  done < .env; \
	fi; \
	go run ./cmd/server

## frontend: run the Next.js dev server
frontend:
	@cd web && pnpm dev

## install: install frontend and backend dependencies
install:
	@cd web && pnpm install
	@cd server && go mod download

## build: build backend and frontend
build:
	@cd server && go build ./...
	@cd web && pnpm build

## test: run backend tests
test:
	@cd server && go test ./...

## migrate: apply database migrations (reads DATABASE_URL from server/.env)
migrate:
	@cd server && \
	DBURL=$$(grep '^DATABASE_URL=' .env | cut -d= -f2-) && \
	migrate -path db/migrations -database "$$DBURL" up

## proto: regenerate gRPC/Connect code (Go + TypeScript)
proto:
	@cd server && buf generate
	@cd web && buf generate

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | \
	  awk -F': ' '{printf "  \033[36m%-9s\033[0m %s\n", $$1, $$2}'
