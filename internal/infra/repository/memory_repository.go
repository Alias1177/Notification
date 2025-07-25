package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/Alias1177/Notification/internal/domain/entity"
	"github.com/Alias1177/Notification/internal/domain/repository"
)

// MemoryRepository реализует NotificationRepository с in-memory хранилищем
type MemoryRepository struct {
	notifications map[string]*entity.Notification
	mutex         sync.RWMutex
}

// NewMemoryRepository создает новый экземпляр MemoryRepository
func NewMemoryRepository() repository.NotificationRepository {
	return &MemoryRepository{
		notifications: make(map[string]*entity.Notification),
	}
}

// Save сохраняет уведомление
func (r *MemoryRepository) Save(notification *entity.Notification) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.notifications[notification.ID] = notification
	return nil
}

// GetByID получает уведомление по ID
func (r *MemoryRepository) GetByID(id string) (*entity.Notification, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	notification, exists := r.notifications[id]
	if !exists {
		return nil, fmt.Errorf("notification not found: %s", id)
	}

	return notification, nil
}

// GetByEmail получает уведомления по email
func (r *MemoryRepository) GetByEmail(email string) ([]*entity.Notification, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var notifications []*entity.Notification
	for _, notification := range r.notifications {
		if notification.Email == email {
			notifications = append(notifications, notification)
		}
	}

	return notifications, nil
}

// GetPending получает все ожидающие отправки уведомления
func (r *MemoryRepository) GetPending() ([]*entity.Notification, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var notifications []*entity.Notification
	for _, notification := range r.notifications {
		if !notification.Sent && !notification.IsExpired() {
			notifications = append(notifications, notification)
		}
	}

	return notifications, nil
}

// MarkAsSent помечает уведомление как отправленное
func (r *MemoryRepository) MarkAsSent(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	notification, exists := r.notifications[id]
	if !exists {
		return fmt.Errorf("notification not found: %s", id)
	}

	notification.MarkAsSent()
	return nil
}

// DeleteExpired удаляет истекшие уведомления
func (r *MemoryRepository) DeleteExpired() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	for id, notification := range r.notifications {
		if notification.ExpiresAt.Before(now) {
			delete(r.notifications, id)
		}
	}

	return nil
}
