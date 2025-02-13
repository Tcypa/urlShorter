FROM golang:1.23rc1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o url_shortener ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/url_shortener .

COPY config/ config/

EXPOSE 8080

CMD ["./url_shortener"]
