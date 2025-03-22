package main

import (
	"Notification/config"
	//connect "Notification/service"
	"Notification/templates"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mid "Notification/internal/middlware"
	cons "Notification/kafka"
)

func main() {
	//TODO :настроить конект к бд
	mid.ColorLogger()
	cfg, err := config.Loading(".env")
	if err != nil {
		panic(err)
	}
	messageHandler := func(value []byte) {
		// Здесь вы можете использовать value (Message.Value) как вам нужно
		// Например, передать его в другие функции или сохранить где-то
		templates.SendEmail(string(value))
	}

	cons.KafkaConnect(cfg, messageHandler)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mid.Recovery)
	//db, err := connect.NewPostgresDB(cfg.DSN)
	//if err != nil {
	//	panic(err)
	//}

}
