# ---------- Stage 1: Build ----------
FROM golang:1.24.3 AS builder

# Ensure static binary for Alpine compatibility
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory inside the build container
WORKDIR /build

# Copy dependency files first (leverages Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy full source code
COPY . .

# Navigate to the folder where main.go lives
WORKDIR /build/cmd/server

# Build the binary and name it 'main'
RUN go build -o main . && ls -lh main

# ---------- Stage 2: Run ----------
FROM alpine:latest

# Add SSL certs (needed for external HTTP requests)
RUN apk --no-cache add ca-certificates

# Set working directory inside runtime container
WORKDIR /app

# Copy only the binary from the build stage
COPY --from=builder /build/cmd/server/main .

# Ensure binary is executable (belt and suspenders)
RUN chmod +x ./main

# Expose app port
EXPOSE 3000

# Run the binary
CMD ["./main"]
