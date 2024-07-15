build:
	@echo "==> Building Docker images..."
	docker-compose build
run:
	@echo "==> Running Docker containers..."
	docker-compose up -d

stop:
	@echo "==> Stopping Docker containers..."
	docker-compose down

.PHONY: build run stop