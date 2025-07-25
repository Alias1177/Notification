package consumer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alias1177/Notification/internal/infra/config"

	"github.com/segmentio/kafka-go"
)

func KafkaConnect(cfg *config.Config, messageHandler func([]byte)) {
	brokerAddress := cfg.KafkaConnect
	topic := cfg.KafkaTopic
	groupID := cfg.KafkaGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал завершения, закрываем подключение...")
		cancel()
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{brokerAddress},
		Topic:           topic,
		GroupID:         groupID,
		MinBytes:        10e3, // 10KB
		MaxBytes:        10e6, // 10MB
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Читаем сообщение
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				fmt.Printf("Ошибка чтения сообщения: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}

			fmt.Printf("Значение: %s\n", string(message.Value))

			// Вызываем обработчик, передавая ему значение сообщения
			messageHandler(message.Value)
		}
	}
}
