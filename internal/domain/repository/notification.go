package repository

import "github.com/Alias1177/Notification/internal/domain/entity"

// NotificationRepository определяет интерфейс для работы с уведомлениями
type NotificationRepository interface {
	Save(notification *entity.Notification) error
	GetByID(id string) (*entity.Notification, error)
	GetByEmail(email string) ([]*entity.Notification, error)
	GetPending() ([]*entity.Notification, error)
	MarkAsSent(id string) error
	DeleteExpired() error
}
