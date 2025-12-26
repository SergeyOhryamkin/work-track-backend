# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies (gcc, musl-dev needed for SQLite)
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main ./cmd/api

# Runtime stage
FROM alpine:latest

# Install runtime dependencies for SQLite
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
