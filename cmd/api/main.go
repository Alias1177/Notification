package main

import (
	"Notification/config"
	"log"
)

func main() {
	cfg, err := config.Loading(".env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
}
