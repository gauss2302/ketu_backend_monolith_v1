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
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/internal/pkg/database/migrations ./internal/pkg/database/migrations

# Create non-root user
RUN adduser -D appuser
USER appuser

# Expose port
EXPOSE 8090

# Set environment variables
ENV APP_ENV=production

# Command to run the application
CMD ["./main"]