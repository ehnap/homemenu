VERSION ?= latest

.PHONY: dev-backend dev-frontend dev build clean docker-build docker-up docker-down

# Development
dev-backend:
	cd backend && export PATH="/home/ehnap/go-sdk/go/bin:$$PATH" && go run .

dev-frontend:
	cd frontend && npm run dev

dev:
	@echo "Run 'make dev-backend' and 'make dev-frontend' in separate terminals"

# Build
build-frontend:
	cd frontend && npm ci && npm run build
	rm -rf backend/static/*
	cp -r frontend/dist/* backend/static/

build-backend: build-frontend
	cd backend && export PATH="/home/ehnap/go-sdk/go/bin:$$PATH" && CGO_ENABLED=1 go build -o ../homemenu .

build: build-backend

# Run
run:
	./homemenu

# Clean
clean:
	rm -f homemenu
	rm -rf frontend/dist
	rm -rf backend/static/*

# Docker
docker-build:
	VERSION=$(VERSION) docker compose build

docker-up:
	VERSION=$(VERSION) docker compose up -d

docker-down:
	VERSION=$(VERSION) docker compose down
