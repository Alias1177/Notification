package main

import (
	"Notification/config"
	"Notification/internal/api/handlers"
	"log/slog"
	"net/http"

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

	go cons.KafkaConnect(cfg, messageHandler)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mid.Recovery)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Post("/api/forgot", handlers.SendCodeForgotPassword)

	slog.Info("Сервер запущен на порту :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}

}
