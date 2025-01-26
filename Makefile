.PHONY: dev prod docker-dev docker-prod migrate-up migrate-down

# Development commands
dev:
	APP_ENV=development go run cmd/app/main.go

docker-dev:
	APP_ENV=development docker-compose up --build -d

# Production commands
prod:
	APP_ENV=production go run cmd/app/main.go

docker-prod:
	APP_ENV=production docker-compose up --build -d

# Docker commands
docker-compose-up:
	docker-compose up --build -d
	@echo "Waiting for PostgreSQL to start..."
	@sleep 5

docker-compose-down:
	docker-compose down

r:
	docker compose down && docker-compose up -d

# Database creation and migrations
create-db:
	docker exec -i postgres psql -U postgres -c "CREATE DATABASE myapp;" || true

# Database migrations
migrate-up:
	migrate -path internal/pkg/database/migrations -database "postgresql://$(APP_POSTGRES_USERNAME):$(APP_POSTGRES_PASSWORD)@$(APP_POSTGRES_HOST):$(APP_POSTGRES_PORT)/$(APP_POSTGRES_DBNAME)?sslmode=$(APP_POSTGRES_SSLMODE)" up

migrate-down-test:
	migrate -path internal/migrations -database "postgresql://postgres:postgres@localhost:5432/myapp?sslmode=disable" down

migrate-down:
	migrate -path internal/pkg/database/migrations -database "postgresql://$(APP_POSTGRES_USERNAME):$(APP_POSTGRES_PASSWORD)@$(APP_POSTGRES_HOST):$(APP_POSTGRES_PORT)/$(APP_POSTGRES_DBNAME)?sslmode=$(APP_POSTGRES_SSLMODE)" down

# Build commands
build:
	go build -o bin/app cmd/app/main.go

run:
	go run cmd/app/main.go

docker-build:
	docker build -t myapp .

docker-run:
	docker run -p 8090:8090 myapp

# Utility commands
clean:
	docker-compose down -v

logs:
	docker-compose logs -f

test:
	go test -v ./...