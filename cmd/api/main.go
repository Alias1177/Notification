package main

import (
	"Notification/config"
	"log/slog"
	"net/http"

	//connect "Notification/service"
	"Notification/templates"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mid "Notification/internal/middlware"
	cons "Notification/kafka"
)

func main() {
	mid.ColorLogger()
	cfg, err := config.Loading(".env")
	if err != nil {
		panic(err)
	}

	messageHandler := func(value []byte) {
		templates.SendEmail(string(value))
	}

	// Запускаем Kafka в отдельной горутине
	go cons.KafkaConnect(cfg, messageHandler)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mid.Recovery)

	// Добавляем эндпоинт для проверки работоспособности
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Запускаем HTTP-сервер
	slog.Info("Сервер запущен на порту :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
