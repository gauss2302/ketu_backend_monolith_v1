FROM golang:1.23-alpine

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Copy migrations
COPY pkg/database/migrations /app/pkg/database/migrations/

# Build the application
RUN go build -o main cmd/app/main.go

# Expose port
EXPOSE 8090

# Run the application
CMD ["./main"]