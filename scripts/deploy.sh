#!/bin/bash

# Скрипт для быстрого деплоя Notification Service
# Использование: ./scripts/deploy.sh [fast|optimized|cached]

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функции для логирования
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Проверка аргументов
BUILD_TYPE=${1:-fast}

case $BUILD_TYPE in
    fast|optimized|cached)
        log_info "Тип сборки: $BUILD_TYPE"
        ;;
    *)
        log_error "Неверный тип сборки. Используйте: fast, optimized, cached"
        exit 1
        ;;
esac

# Проверка наличия необходимых файлов
log_info "Проверка файлов проекта..."
if [ ! -f "docker-compose.yml" ]; then
    log_error "docker-compose.yml не найден!"
    exit 1
fi

if [ ! -f "Dockerfile" ]; then
    log_error "Dockerfile не найден!"
    exit 1
fi

# Проверка .env файла
if [ ! -f ".env" ]; then
    log_warning ".env файл не найден, копируем из примера..."
    if [ -f "env.example" ]; then
        cp env.example .env
        log_success ".env файл создан из env.example"
    else
        log_warning "env.example не найден, создаем базовый .env"
        cat > .env << EOF
# Notification Service Environment Variables
MAIL=your-email@gmail.com
SECRET=your-app-password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
APP_ENV=production
EOF
    fi
fi

# Остановка существующих контейнеров
log_info "Остановка существующих контейнеров..."
docker-compose down --timeout 30 || true

# Очистка
log_info "Очистка Docker кэша..."
docker system prune -f || true

# Сборка в зависимости от типа
log_info "Начинаем сборку (тип: $BUILD_TYPE)..."
case $BUILD_TYPE in
    fast)
        DOCKER_BUILDKIT=1 docker-compose build notification-service
        ;;
    optimized)
        DOCKER_BUILDKIT=1 docker-compose build --no-cache notification-service
        ;;
    cached)
        DOCKER_BUILDKIT=1 docker-compose build \
            --build-arg BUILDKIT_INLINE_CACHE=1 \
            notification-service
        ;;
esac

# Запуск сервисов
log_info "Запуск сервисов..."
docker-compose up -d

# Ожидание запуска
log_info "Ожидание запуска сервисов..."
sleep 15

# Проверка статуса
log_info "Проверка статуса контейнеров..."
docker-compose ps

# Health check
log_info "Проверка health endpoint..."
for i in {1..10}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_success "✅ Notification Service здоров!"
        break
    else
        log_warning "⏳ Ожидание готовности сервиса... (попытка $i/10)"
        if [ $i -eq 10 ]; then
            log_error "❌ Сервис не стал готовым после 10 попыток"
            log_info "Проверка логов..."
            docker-compose logs --tail=20 notification-service
            exit 1
        fi
        sleep 10
    fi
done

# Финальная проверка
log_info "Финальная проверка..."
docker-compose ps

# Проверка размера образа
log_info "Размер образа:"
docker images notification-service --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

# Проверка логов на ошибки
log_info "Проверка логов на ошибки..."
if docker-compose logs --tail=10 notification-service | grep -i error; then
    log_warning "Найдены ошибки в логах"
else
    log_success "Ошибок в логах не найдено"
fi

log_success "🎉 Деплой Notification Service завершен успешно!"
log_info "Сервис доступен по адресу: http://localhost:8080"
log_info "Health endpoint: http://localhost:8080/health" 