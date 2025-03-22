package main

import (
	"Notification/config"

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
	cons.KafkaConnect(cfg)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mid.Recovery)

}
