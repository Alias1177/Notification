# Этап сборки зависимостей
FROM golang:1.24-alpine AS deps
WORKDIR /app

# Устанавливаем git для приватных репозиториев (если нужно)
RUN apk add --no-cache git

COPY go.mod go.sum ./

# Используем кэширование для модулей
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Этап сборки
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Копируем зависимости из предыдущего этапа
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY go.mod go.sum ./

# Копируем только исходный код (без тестов)
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Сборка с максимальными оптимизациями
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
        -ldflags="-w -s -extldflags=-static" \
        -trimpath \
        -o notification-service ./cmd/api

# Финальный этап - используем distroless для минимального размера
FROM gcr.io/distroless/static-debian12:nonroot

# Копируем SSL сертификаты для SMTP/TLS соединений
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем бинарник
COPY --from=builder /app/notification-service /app/notification-service

# Копируем только необходимые шаблоны
COPY --from=builder /app/templates /app/templates

# Устанавливаем рабочую директорию
WORKDIR /app

# Запуск от непривилегированного пользователя
USER 65532:65532

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/notification-service"] || exit 1

# Запуск приложения
ENTRYPOINT ["/app/notification-service"]