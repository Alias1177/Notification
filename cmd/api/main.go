package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alias1177/Notification/internal/delivery/http/handlers"
	"github.com/Alias1177/Notification/internal/delivery/http/middlware"
	"github.com/Alias1177/Notification/internal/delivery/kafka"
	"github.com/Alias1177/Notification/internal/domain/service"
	"github.com/Alias1177/Notification/internal/infra/config"
	"github.com/Alias1177/Notification/internal/infra/kafka/consumer"
	"github.com/Alias1177/Notification/internal/infra/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Настройка логгера
	middlware.ColorLogger()

	// Загрузка конфигурации
	cfg, err := config.Loading(".env")
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Инициализация зависимостей
	notificationRepo := repository.NewMemoryRepository()
	notificationService := service.NewNotificationService(notificationRepo)
	emailService := service.NewEmailService()

	// Создание handlers
	notificationHandler := handlers.NewNotificationHandler(notificationService, emailService)
	kafkaHandler := kafka.NewKafkaHandler(notificationService, emailService)

	// Запуск Kafka consumer
	go func() {
		messageHandler := func(value []byte) {
			if err := kafkaHandler.HandleMessage(value); err != nil {
				slog.Error("Failed to handle Kafka message", "error", err, "message", string(value))
			}
		}
		consumer.KafkaConnect(cfg, messageHandler)
	}()

	// Настройка HTTP роутера
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middlware.Recovery)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{}, // Пустой список (разрешим динамически)
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", notificationHandler.HealthCheck)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/forgot", notificationHandler.SendPasswordResetCode)
		r.Post("/validate", notificationHandler.ValidatePasswordResetCode)
	})

	// Настройка HTTP сервера
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		slog.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown error", "error", err)
		}
	}()

	slog.Info("Server started", "port", "8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}
