build:
	docker-compose build

up:
	docker-compose up

down:
	docker-compose down

restart: down up

refresh: down
	docker-compose up --build -d

clean: down
	docker system prune -f

ps:
	docker-compose ps

load-test:
	k6 run load-test.js

help:
	@echo "Available commands:"
	@echo "  make build    - Build Docker images"
	@echo "  make up       - Start all services in foreground"
	@echo "  make down     - Stop and remove all services"
	@echo "  make restart  - Restart services (down + up)"
	@echo "  make refresh  - Rebuild and restart services"
	@echo "  make clean    - Stop services and clean Docker system"
	@echo "  make ps       - Show status of running services"
	@echo "  make help     - Show this help message"
	@echo "  make load-test          - Run load test locally"
