#!/bin/bash

set -e

BACKEND_DIR="./backend"
SERVICES=("admin-service" "customer-service" "gateway-service" "menu-service" "order-service" "payment-service")

echo "🔧 Generating Makefiles for services..."

for svc in "${SERVICES[@]}"; do
  SVC_DIR="${BACKEND_DIR}/${svc}"
  SVC_NAME=$(echo "$svc" | sed 's/-service//')
  DOCKER_IMG="aditis-kitchen/${svc}"

  echo "📁 Creating Makefile in $SVC_DIR..."

  cat > "${SVC_DIR}/Makefile" <<EOF
SERVICE_NAME := ${svc}
BINARY_NAME := ${SVC_NAME}
DOCKER_IMAGE := ${DOCKER_IMG}
PORT := 8080

build:
\t@echo "🔨 Building \$(SERVICE_NAME)..."
\tgo build -o \$(BINARY_NAME) main.go

run:
\t@echo "🚀 Running \$(SERVICE_NAME)..."
\t./\$(BINARY_NAME)

test:
\t@echo "🧪 Running tests for \$(SERVICE_NAME)..."
\tgo test ./...

clean:
\t@echo "🧹 Cleaning up \$(SERVICE_NAME)..."
\trm -f \$(BINARY_NAME)

docker-build:
\t@echo "🐳 Building Docker image for \$(SERVICE_NAME)..."
\tdocker build -t \$(DOCKER_IMAGE) .

docker-run:
\t@echo "🐳 Running Docker container for \$(SERVICE_NAME)..."
\tdocker run -d -p \$(PORT):\$(PORT) --name \$(SERVICE_NAME) \$(DOCKER_IMAGE)

docker-push:
\t@echo "📦 Pushing Docker image to registry for \$(SERVICE_NAME)..."
\tdocker push \$(DOCKER_IMAGE)
EOF

  chmod +x "${SVC_DIR}/Makefile"
done

echo "✅ All service Makefiles generated!"
