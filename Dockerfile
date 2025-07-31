# Этап сборки
FROM golang:1.24-alpine AS builder

# Установка необходимых зависимостей
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование файлов go.mod и go.sum
COPY go.mod ./
COPY go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY docker .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/api

# Финальный этап
FROM alpine:latest

# Установка необходимых пакетов
RUN apk --no-cache add ca-certificates tzdata

# Установка рабочей директории
WORKDIR /app

# Копирование исполняемого файла из этапа сборки
COPY --from=builder /app/notification-service /app/notification-service

# Копирование шаблонов
COPY --from=builder /app/templates /app/templates

# Запуск приложения
CMD ["/app/notification-service"]