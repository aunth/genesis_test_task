FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /tmp/go-build && chmod 777 /tmp/go-build
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/internal/static ./internal/static
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["./main"]