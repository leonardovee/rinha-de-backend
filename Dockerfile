FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o rinha-de-backend ./cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/rinha-de-backend /usr/local/bin/rinha-de-backend

ENV DATABASE_URL=postgres://postgres:postgres@db/postgres
ENV REDIS_URL=redis:6379

CMD ["rinha-de-backend"]