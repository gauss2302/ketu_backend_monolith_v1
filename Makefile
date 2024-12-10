.PHONY: migrate-up migrate-down docker-up docker-down

# Database migrations
migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5430/myapp?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5430/myapp?sslmode=disable" down

# Docker commands
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down