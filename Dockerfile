# Stage 1: Build Stage
FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/server
RUN go build -o main .

# Stage 2: Run Stage (Alpine is a tiny base image)
FROM alpine:latest

# Install certificates to make HTTP requests work (optional but important)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the compiled Go binary from the builder stage
COPY --from=builder /app/cmd/server/main .

EXPOSE 3000

CMD ["./main"]
