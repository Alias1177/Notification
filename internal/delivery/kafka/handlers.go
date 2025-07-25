package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Alias1177/Notification/internal/domain/service"
)

// KafkaHandler обрабатывает сообщения из Kafka
type KafkaHandler struct {
	notificationService *service.NotificationService
	emailService        *service.EmailService
}

// NewKafkaHandler создает новый экземпляр KafkaHandler
func NewKafkaHandler(
	notificationService *service.NotificationService,
	emailService *service.EmailService,
) *KafkaHandler {
	return &KafkaHandler{
		notificationService: notificationService,
		emailService:        emailService,
	}
}

// PasswordResetRequest представляет запрос на восстановление пароля из Kafka
type PasswordResetRequest struct {
	Email  string `json:"email"`
	UserID string `json:"user_id,omitempty"`
}

// RegistrationRequest представляет запрос на регистрацию из Kafka
type RegistrationRequest struct {
	Email    string `json:"email"`
	Username string `json:"username,omitempty"`
}

// HandleMessage обрабатывает входящие сообщения из Kafka
func (h *KafkaHandler) HandleMessage(message []byte) error {
	// Пытаемся определить тип сообщения
	// Сначала пробуем как PasswordResetRequest
	var passwordResetReq PasswordResetRequest
	if err := json.Unmarshal(message, &passwordResetReq); err == nil && passwordResetReq.Email != "" {
		return h.handlePasswordResetRequest(passwordResetReq)
	}

	// Если не получилось, пробуем как RegistrationRequest
	var registrationReq RegistrationRequest
	if err := json.Unmarshal(message, &registrationReq); err == nil && registrationReq.Email != "" {
		return h.handleRegistrationRequest(registrationReq)
	}

	// Если не получилось распарсить как структуру, пробуем как простой email
	email := string(message)
	if email != "" {
		// По умолчанию считаем это запросом на регистрацию
		return h.handleRegistrationRequest(RegistrationRequest{Email: email})
	}

	return fmt.Errorf("unable to parse message: %s", string(message))
}

// handlePasswordResetRequest обрабатывает запрос на восстановление пароля
func (h *KafkaHandler) handlePasswordResetRequest(req PasswordResetRequest) error {
	log.Printf("Processing password reset request for email: %s", req.Email)

	// Создаем уведомление для восстановления пароля
	notification, err := h.notificationService.CreatePasswordResetNotification(req.Email)
	if err != nil {
		log.Printf("Failed to create password reset notification: %v", err)
		return fmt.Errorf("failed to create password reset notification: %w", err)
	}

	// Отправляем email с кодом
	if err := h.emailService.SendPasswordResetEmail(req.Email, notification.Code); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	// Помечаем уведомление как отправленное
	if err := h.notificationService.MarkNotificationAsSent(notification.ID); err != nil {
		log.Printf("Failed to mark notification as sent: %v", err)
		// Не возвращаем ошибку, так как email уже отправлен
	}

	log.Printf("Password reset code successfully sent to %s", req.Email)
	return nil
}

// handleRegistrationRequest обрабатывает запрос на регистрацию
func (h *KafkaHandler) handleRegistrationRequest(req RegistrationRequest) error {
	log.Printf("Processing registration request for email: %s", req.Email)

	// Создаем уведомление для регистрации
	notification, err := h.notificationService.CreateRegistrationNotification(req.Email)
	if err != nil {
		log.Printf("Failed to create registration notification: %v", err)
		return fmt.Errorf("failed to create registration notification: %w", err)
	}

	// Отправляем email подтверждения регистрации
	if err := h.emailService.SendRegistrationEmail(req.Email); err != nil {
		log.Printf("Failed to send registration email: %v", err)
		return fmt.Errorf("failed to send registration email: %w", err)
	}

	// Помечаем уведомление как отправленное
	if err := h.notificationService.MarkNotificationAsSent(notification.ID); err != nil {
		log.Printf("Failed to mark notification as sent: %v", err)
		// Не возвращаем ошибку, так как email уже отправлен
	}

	log.Printf("Registration confirmation email successfully sent to %s", req.Email)
	return nil
}
