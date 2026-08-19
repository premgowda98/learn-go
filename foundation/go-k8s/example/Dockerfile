# Start from the official Golang image for building
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app main.go

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]
