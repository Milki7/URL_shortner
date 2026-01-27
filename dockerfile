# Step 1: Build the Go binary
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/server/main.go

# Step 2: Create the small runtime image
FROM alpine:latest
WORKDIR /root/
# Install sqlite for the DB file to work
RUN apk add --no-cache sqlite
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]