package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/Alias1177/Notification/internal/domain/entity"
	"github.com/Alias1177/Notification/internal/domain/repository"
)

// NotificationService содержит бизнес-логику для работы с уведомлениями
type NotificationService struct {
	repo repository.NotificationRepository
}

// NewNotificationService создает новый экземпляр NotificationService
func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

// GeneratePasswordResetCode генерирует безопасный код для сброса пароля
func (s *NotificationService) GeneratePasswordResetCode() string {
	// Используем crypto/rand для криптографически безопасной генерации
	max := big.NewInt(10000) // 4-значный код
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback к менее безопасному методу в случае ошибки
		return fmt.Sprintf("%04d", time.Now().UnixNano()%10000)
	}
	return fmt.Sprintf("%04d", n.Int64())
}

// CreatePasswordResetNotification создает уведомление для сброса пароля
func (s *NotificationService) CreatePasswordResetNotification(email string) (*entity.Notification, error) {
	// Проверяем, нет ли уже активного кода для этого email
	existingNotifications, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing notifications: %w", err)
	}

	// Инвалидируем старые коды для этого email
	for _, notification := range existingNotifications {
		if notification.Type == entity.PasswordReset && !notification.IsExpired() {
			// Помечаем как отправленное, чтобы деактивировать
			notification.MarkAsSent()
			if err := s.repo.Save(notification); err != nil {
				return nil, fmt.Errorf("failed to invalidate old notification: %w", err)
			}
		}
	}

	// Генерируем новый код
	code := s.GeneratePasswordResetCode()

	// Создаем новое уведомление
	notification := entity.NewPasswordResetNotification(email, code)

	// Сохраняем в репозиторий
	if err := s.repo.Save(notification); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	return notification, nil
}

// CreateRegistrationNotification создает уведомление для регистрации
func (s *NotificationService) CreateRegistrationNotification(email string) (*entity.Notification, error) {
	notification := entity.NewRegistrationNotification(email)

	if err := s.repo.Save(notification); err != nil {
		return nil, fmt.Errorf("failed to save registration notification: %w", err)
	}

	return notification, nil
}

// GetNotificationByID получает уведомление по ID
func (s *NotificationService) GetNotificationByID(id string) (*entity.Notification, error) {
	return s.repo.GetByID(id)
}

// GetNotificationsByEmail получает уведомления по email
func (s *NotificationService) GetNotificationsByEmail(email string) ([]*entity.Notification, error) {
	return s.repo.GetByEmail(email)
}

// MarkNotificationAsSent помечает уведомление как отправленное
func (s *NotificationService) MarkNotificationAsSent(id string) error {
	return s.repo.MarkAsSent(id)
}

// CleanupExpiredNotifications удаляет истекшие уведомления
func (s *NotificationService) CleanupExpiredNotifications() error {
	return s.repo.DeleteExpired()
}
