FROM golang:1.24.3


WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/server

RUN go build -o main .

EXPOSE 3000

CMD ["./main"]