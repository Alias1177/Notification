# Используем официальный образ Go для сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o notification ./cmd/api

# Используем минимальный образ Alpine для финального контейнера
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем собранное приложение из стадии builder
COPY --from=builder /app/notification .

# Копируем папку templates (если она нужна в runtime)
COPY --from=builder /app/internal/templates ./templates

# Команда, которая выполняется при запуске контейнера
CMD ["./notification"]