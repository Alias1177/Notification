package repository

import "github.com/Alias1177/Notification/internal/domain/entity"

// NotificationRepository определяет интерфейс для работы с уведомлениями
type NotificationRepository interface {
	// Save сохраняет уведомление
	Save(notification *entity.Notification) error

	// GetByID получает уведомление по ID
	GetByID(id string) (*entity.Notification, error)

	// GetByEmail получает уведомления по email
	GetByEmail(email string) ([]*entity.Notification, error)

	// GetPending получает все ожидающие отправки уведомления
	GetPending() ([]*entity.Notification, error)

	// MarkAsSent помечает уведомление как отправленное
	MarkAsSent(id string) error

	// DeleteExpired удаляет истекшие уведомления
	DeleteExpired() error
}
