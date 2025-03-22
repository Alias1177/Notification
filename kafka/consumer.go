package consumer

import (
	"Notification/config"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func KafkaConnect(cfg *config.Config) {
	brokerAddress := cfg.KafkaConnect // Адрес вашего внешнего Kafka-брокера
	topic := cfg.KafkaTopic           // Имя топика для чтения
	groupID := cfg.KafkaGroup         // ID группы потребителей

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
				// Если контекст отменен, выходим без ошибки
				if ctx.Err() != nil {
					return
				}
				fmt.Printf("Ошибка чтения сообщения: %v\n", err)
				time.Sleep(1 * time.Second) // Пауза перед повторной попыткой
				continue
			}

			fmt.Printf("Значение: %s\n", string(message.Value)) //затычка

		}
	}
}
