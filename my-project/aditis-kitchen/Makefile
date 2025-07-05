PROJECT := aditis-kitchen
SERVICES := admin-service customer-service gateway-service menu-service order-service payment-service
BACKEND_DIR := backend

.PHONY: all build run clean test docker-build docker-run docker-push

all: build

build:
	@for svc in $(SERVICES); do \
		echo "🔧 Building $$svc..."; \
		$(MAKE) -C $(BACKEND_DIR)/$$svc build; \
	done

run:
	@for svc in $(SERVICES); do \
		echo "🚀 Running $$svc..."; \
		$(MAKE) -C $(BACKEND_DIR)/$$svc run & \
	done
	@wait

clean:
	@for svc in $(SERVICES); do \
		echo "🧹 Cleaning $$svc..."; \
		$(MAKE) -C $(BACKEND_DIR)/$$svc clean; \
	done

test:
	@for svc in $(SERVICES); do \
		echo "🧪 Testing $$svc..."; \
		$(MAKE) -C $(BACKEND_DIR)/$$svc test; \
	done

docker-build:
	@for svc in $(SERVICES); do \
		echo "🐳 Docker build $$svc..."; \
		$(MAKE) -C $(BACKEND_DIR)/$$svc docker-build; \
	done

docker-run:
	@for svc in $(SERVICES); do \
		echo "🐳 Docker run $$svc..."; \
		$(MAKE) -C $(BACKEND_DIR)/$$svc docker-run & \
	done
	@wait

docker-push:
	@for svc in $(SERVICES); do \
		echo "📦 Pushing Docker image for $$svc..."; \
		$(MAKE) -C $(BACKEND_DIR)/$$svc docker-push; \
	done
