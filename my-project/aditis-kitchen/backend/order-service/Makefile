SERVICE_NAME := order-service
BINARY_NAME := order
DOCKER_IMAGE := aditis-kitchen/order-service
PORT := 8080

build:
\t@echo "🔨 Building $(SERVICE_NAME)..."
\tgo build -o $(BINARY_NAME) main.go

run:
\t@echo "🚀 Running $(SERVICE_NAME)..."
\t./$(BINARY_NAME)

test:
\t@echo "🧪 Running tests for $(SERVICE_NAME)..."
\tgo test ./...

clean:
\t@echo "🧹 Cleaning up $(SERVICE_NAME)..."
\trm -f $(BINARY_NAME)

docker-build:
\t@echo "🐳 Building Docker image for $(SERVICE_NAME)..."
\tdocker build -t $(DOCKER_IMAGE) .

docker-run:
\t@echo "🐳 Running Docker container for $(SERVICE_NAME)..."
\tdocker run -d -p $(PORT):$(PORT) --name $(SERVICE_NAME) $(DOCKER_IMAGE)

docker-push:
\t@echo "📦 Pushing Docker image to registry for $(SERVICE_NAME)..."
\tdocker push $(DOCKER_IMAGE)
