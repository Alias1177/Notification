package entity

import "time"

// NotificationType определяет тип уведомления
type NotificationType string

const (
	PasswordReset NotificationType = "password_reset"
	Registration  NotificationType = "registration"
)

// Notification представляет уведомление
type Notification struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Email     string           `json:"email"`
	Code      string           `json:"code,omitempty"` // Код только для внутреннего использования
	CreatedAt time.Time        `json:"created_at"`
	ExpiresAt time.Time        `json:"expires_at"`
	Sent      bool             `json:"sent"`
}

// NewPasswordResetNotification создает новое уведомление для сброса пароля
func NewPasswordResetNotification(email, code string) *Notification {
	now := time.Now()
	return &Notification{
		ID:        generateID(),
		Type:      PasswordReset,
		Email:     email,
		Code:      code,
		CreatedAt: now,
		ExpiresAt: now.Add(10 * time.Minute), // Код истекает через 10 минут
		Sent:      false,
	}
}

// NewRegistrationNotification создает новое уведомление для регистрации
func NewRegistrationNotification(email string) *Notification {
	now := time.Now()
	return &Notification{
		ID:        generateID(),
		Type:      Registration,
		Email:     email,
		CreatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour), // Регистрационное письмо истекает через 24 часа
		Sent:      false,
	}
}

// IsExpired проверяет, истекло ли уведомление
func (n *Notification) IsExpired() bool {
	return time.Now().After(n.ExpiresAt)
}

// MarkAsSent помечает уведомление как отправленное
func (n *Notification) MarkAsSent() {
	n.Sent = true
}

// generateID генерирует уникальный ID (упрощенная версия)
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString генерирует случайную строку заданной длины
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
