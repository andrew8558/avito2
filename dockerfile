FROM golang:1.24.2

WORKDIR /app
COPY . .

RUN go install github.com/pressly/goose/cmd/goose@latest

RUN go build -o /build ./cmd \
    && go clean -cache -modcache

EXPOSE 8080

CMD goose -dir /app/internal/db/migrations postgres "$DATABASE_URL" up && /build
