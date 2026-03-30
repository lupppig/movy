# Variables
COMPOSE_FILE=infra/docker-compose.dev.yml
ENV_FILE=.env

.PHONY: dev-up dev-down dev-build dev-logs

# Start the development environment
dev-up:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up

# Stop the development environment
dev-down:
	docker compose -f $(COMPOSE_FILE) down

# Force a rebuild (useful when you change go.mod)
dev-build:
	docker compose -f $(COMPOSE_FILE) build --no-cache

# View real-time logs from the Gin app
dev-logs:
	docker compose -f $(COMPOSE_FILE) logs -f