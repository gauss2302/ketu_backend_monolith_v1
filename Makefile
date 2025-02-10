.PHONY: dev prod docker-dev docker-prod migrate-up migrate-down

# Development commands
dev:
	APP_ENV=development go run cmd/app/main.go

docker-dev:
	APP_ENV=development docker-compose up --build -d
	@echo "Waiting for services to start..."
	@sleep 10
	@echo "Creating database and running migrations..."
	@docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE ke2;" || true
	@echo "Database created (or already exists). Running migrations..."
	@migrate -path internal/pkg/database/migrations \
		-database "postgresql://postgres:postgres@localhost:5432/ke2?sslmode=disable" up || \
		(echo "Migration failed. Retrying in 5s..." && sleep 5 && \
		migrate -path internal/pkg/database/migrations \
		-database "postgresql://postgres:postgres@localhost:5432/ke2?sslmode=disable" up)
	@echo "Checking service status..."
	@docker-compose ps
	@echo "\nChecking application logs..."
	@docker-compose logs app

# Production commands
prod:
	APP_ENV=production go run cmd/app/main.go

docker-prod:
	APP_ENV=production docker-compose up --build -d

# Docker commands
up:
	docker-compose up --build -d
	@echo "Waiting for PostgreSQL to start..."
	@sleep 5

down:
	docker-compose down

r:
	docker compose down && docker-compose up -d

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	@docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE ke2;" 2>/dev/null || true
	migrate -path internal/pkg/database/migrations \
		-database "postgresql://postgres:postgres@localhost:5432/ke2?sslmode=disable" up

migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path internal/pkg/database/migrations \
		-database "postgresql://postgres:postgres@localhost:5432/ke2?sslmode=disable" down

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
	docker system prune -f
	rm -rf postgres_data

logs:
	docker-compose logs -f

test:
	go test -v ./...

# Debug commands
status:
	@echo "Docker containers status:"
	@docker ps
	@echo "\nApplication logs:"
	@docker-compose logs app
	@echo "\nPostgres logs:"
	@docker-compose logs postgres

# Add new commands for debugging
check-app:
	@echo "Checking application status..."
	@curl -s http://localhost:8090/health || echo "Application is not responding"
	@echo "\nApplication logs:"
	@docker-compose logs app

check-db:
	@echo "Checking database status..."
	@docker-compose exec -T postgres pg_isready -U postgres || echo "Database is not responding"
	@echo "\nDatabase logs:"
	@docker-compose logs postgres

debug:
	@make status
	@echo "\nChecking individual services:"
	@make check-app
	@make check-db