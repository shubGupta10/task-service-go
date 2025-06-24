# Stage 1: Build Stage
FROM golang:1.24.3 AS builder

# Set working directory to project root
WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

# Go to the directory where main.go exists and build
WORKDIR /build/cmd/server
RUN go build -o main . && ls -l main

# Stage 2: Run Stage
FROM alpine:latest

# Add certificates (important for HTTPS and HTTP clients)
RUN apk --no-cache add ca-certificates

# Set working directory in runtime container
WORKDIR /app

# Copy only the built binary
COPY --from=builder /build/cmd/server/main .

# Expose port (adjust if your app uses a different one)
EXPOSE 3000

# Start the binary
CMD ["./main"]
