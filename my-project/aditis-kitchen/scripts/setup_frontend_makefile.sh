#!/bin/bash

# Navigate to the project root relative to scripts/
cd "$(dirname "$0")/.." || exit 1

FRONTEND_DIR="frontend"
MAKEFILE_PATH="$FRONTEND_DIR/Makefile"

echo "📂 Creating Makefile in $FRONTEND_DIR..."

cat > "$MAKEFILE_PATH" << 'EOF'
# Project Config
APP_NAME := aditis-kitchen-frontend
IMAGE_NAME := $(APP_NAME):latest
PORT := 3000
DOCKER_PORT := 80
CONTAINER_NAME := $(APP_NAME)

ifneq ("$(wildcard .env)","")
    include .env
    export
endif

install:
	@echo "📦 Installing dependencies..."
	npm install

dev:
	@echo "🚀 Starting Vite dev server on port $(PORT)..."
	npm run dev -- --port $(PORT)

build:
	@echo "⚙️  Building frontend for production..."
	npm run build

lint:
	@echo "🔍 Running ESLint..."
	npx eslint src --ext .js,.jsx

format:
	@echo "🎨 Formatting code using Prettier..."
	npx prettier --write src

test:
	@echo "🧪 Running tests..."
	npm test

docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(IMAGE_NAME) .

docker-run:
	@echo "🐳 Running Docker container on port $(PORT)..."
	docker run -d -p $(PORT):$(DOCKER_PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

docker-clean:
	@echo "🧹 Cleaning up Docker image and container..."
	-docker rm -f $(CONTAINER_NAME)
	-docker rmi $(IMAGE_NAME)

env:
	@echo "🌍 Displaying environment variables..."
	@cat .env || echo "⚠️  No .env file found."

clean:
	@echo "🧼 Cleaning node_modules and dist..."
	rm -rf node_modules dist

help:
	@echo "🚀 Makefile Targets for Frontend:"
	@echo "  install         Install dependencies"
	@echo "  dev             Start Vite dev server"
	@echo "  build           Build frontend for production"
	@echo "  lint            Lint code with ESLint"
	@echo "  format          Format code with Prettier"
	@echo "  test            Run frontend tests"
	@echo "  docker-build    Build Docker image"
	@echo "  docker-run      Run Docker container"
	@echo "  docker-clean    Remove Docker image and container"
	@echo "  env             Show environment variables"
	@echo "  clean           Remove node_modules and dist"
EOF

echo "✅ Makefile created at $MAKEFILE_PATH"
