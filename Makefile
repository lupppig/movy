# Variables
COMPOSE_FILE=infra/docker-compose.dev.yml
ENV_FILE=.env

DB_URL=postgres://postgres:postgres@localhost:5433/booking?sslmode=disable

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

go-gen:
	cd src && go generate ./...


migrate-create:
	migrate create -ext sql -dir ./src/migrations -seq $(name)

migrate-up:
	migrate -path src/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path src/migrations -database "$(DB_URL)" down 1

migrate-force:
	migrate -path src/migrations -database "$(DB_URL)" force $(version)