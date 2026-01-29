# Step 1: Build the Go binary
FROM golang:alpine AS builder

# Install build essentials
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy the module files first to leverage Docker caching
COPY go.mod go.sum ./

# Try downloading; if it fails, 'go build' will catch it anyway
RUN go mod download || true

COPY . .
RUN go build -o main cmd/server/main.go