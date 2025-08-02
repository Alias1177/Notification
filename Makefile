# Переменные
IMAGE_NAME = notification-service
TAG = latest

# Основные команды
.PHONY: help build build-fast build-optimized clean test docker-build docker-push

# Помощь
help:
	@echo "Доступные команды для Notification Service:"
	@echo "  build-fast        - Быстрая сборка для разработки"
	@echo "  build-optimized   - Оптимизированная сборка для production"
	@echo "  build-cached      - Сборка с кэшированием (BuildKit)"
	@echo "  clean            - Очистка Docker кэша"
	@echo "  test             - Запуск тестов"
	@echo "  dev              - Локальная разработка"
	@echo "  deploy           - Production деплой"
	@echo "  logs             - Просмотр логов"
	@echo "  status           - Статус сервисов"
	@echo "  health           - Проверка health endpoints"
	@echo "  size             - Размер образа"

# Быстрая сборка (для разработки)
build-fast:
	@echo "🚀 Быстрая сборка Notification Service..."
	DOCKER_BUILDKIT=1 docker-compose build notification-service
	@echo "Сборка завершена"

# Оптимизированная сборка (для production)
build-optimized:
	@echo "⚡ Оптимизированная сборка Notification Service..."
	DOCKER_BUILDKIT=1 docker-compose build --no-cache notification-service
	@echo "Сборка завершена"

# Сборка с кэшированием (требует BuildKit)
build-cached:
	@echo "💾 Сборка с кэшированием..."
	DOCKER_BUILDKIT=1 docker-compose build \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		notification-service
	@echo "Сборка завершена"

# Очистка
clean:
	@echo "🧹 Очистка Docker кэша..."
	docker system prune -f
	docker builder prune -f
	@echo "Кэш очищен"

# Тесты
test:
	@echo "🧪 Запуск тестов..."
	go test -v ./...

# Локальная разработка
dev:
	@echo "🔧 Запуск в режиме разработки..."
	docker-compose up --build

# Production деплой
deploy:
	@echo "🚀 Production деплой..."
	docker-compose up -d --build

# Мониторинг
logs:
	docker-compose logs -f notification-service

# Статус
status:
	@echo "📊 Статус сервисов:"
	docker-compose ps

# Health check
health:
	@echo "🏥 Проверка health endpoints:"
	@echo "Notification Service: http://localhost:8080/health"
	@curl -s http://localhost:8080/health | jq . || echo "Notification service недоступен"

# Размер образа
size:
	@echo "📊 Размер образа:"
	docker images $(IMAGE_NAME) --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

# Перезапуск
restart:
	@echo "🔄 Перезапуск сервиса..."
	docker-compose restart notification-service
	@echo "Сервис перезапущен"

# Остановка
stop:
	@echo "⏹️ Остановка сервисов..."
	docker-compose down

# Полная очистка
clean-all:
	@echo "🧹 Полная очистка всех контейнеров и образов..."
	docker-compose down -v
	docker system prune -af
	docker volume prune -f
	@echo "Полная очистка завершена"

# Проверка портов
ports:
	@echo "🔍 Проверка занятых портов:"
	@netstat -tlnp | grep -E ':(8080)' || echo "Порты свободны" 