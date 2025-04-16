# Multi-stage Dockerfile for both development and production

# Base stage with common dependencies
FROM golang:1.23.4 AS base
WORKDIR /usr/src/api
COPY . /usr/src/api
RUN go mod tidy

# Development stage with air for hot-reloading
FROM base AS development
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

# Production stage with compiled binary
FROM base AS production
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./src/cmd/main.go
# Run the compiled binary directly
CMD ["./app"]
