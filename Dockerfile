# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required build tools
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/app

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Copy the binary and migrations
COPY --from=builder /app/main .
COPY --from=builder /app/internal/pkg/database/migrations ./internal/pkg/database/migrations

# Expose port
EXPOSE 8090

# Command to run the application
CMD ["./main"]