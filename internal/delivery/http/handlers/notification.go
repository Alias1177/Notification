package handlers

import (
	"net/http"

	"github.com/Alias1177/Notification/internal/domain/service"

	"github.com/go-chi/render"
)

// NotificationHandler обрабатывает HTTP запросы для уведомлений
type NotificationHandler struct {
	notificationService *service.NotificationService
	emailService        *service.EmailService
}

// NewNotificationHandler создает новый экземпляр NotificationHandler
func NewNotificationHandler(
	notificationService *service.NotificationService,
	emailService *service.EmailService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		emailService:        emailService,
	}
}

// PasswordResetRequest представляет запрос на восстановление пароля
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetResponse представляет ответ на запрос восстановления пароля
type PasswordResetResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// SendPasswordResetCode обрабатывает запрос на отправку кода восстановления пароля
func (h *NotificationHandler) SendPasswordResetCode(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Валидация email
	if req.Email == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Email is required"})
		return
	}

	// Создаем уведомление для восстановления пароля
	notification, err := h.notificationService.CreatePasswordResetNotification(req.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to create password reset notification"})
		return
	}

	// Отправляем email с кодом
	if err := h.emailService.SendPasswordResetEmail(req.Email, notification.Code); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to send password reset email"})
		return
	}

	// Помечаем уведомление как отправленное
	if err := h.notificationService.MarkNotificationAsSent(notification.ID); err != nil {
		// Логируем ошибку, но не возвращаем пользователю
		// Email уже отправлен, это не критично
	}

	// Возвращаем успешный ответ (БЕЗ КОДА!)
	render.JSON(w, r, PasswordResetResponse{
		Message: "Password reset code has been sent to your email",
		Status:  "success",
	})
}

// ValidateCodeRequest представляет запрос на валидацию кода
type ValidateCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

// ValidateCodeResponse представляет ответ на валидацию кода
type ValidateCodeResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}

// ValidatePasswordResetCode валидирует код восстановления пароля
func (h *NotificationHandler) ValidatePasswordResetCode(w http.ResponseWriter, r *http.Request) {
	var req ValidateCodeRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Валидация входных данных
	if req.Email == "" || req.Code == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Email and code are required"})
		return
	}

	// Получаем уведомления для данного email
	notifications, err := h.notificationService.GetNotificationsByEmail(req.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to validate code"})
		return
	}

	// Ищем активный код
	var isValid bool
	for _, notification := range notifications {
		if notification.Type == "password_reset" &&
			notification.Code == req.Code &&
			!notification.IsExpired() &&
			!notification.Sent {
			isValid = true
			break
		}
	}

	render.JSON(w, r, ValidateCodeResponse{
		Valid:   isValid,
		Message: "Code validation completed",
	})
}

// HealthCheck проверяет состояние сервиса
func (h *NotificationHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status":  "healthy",
		"service": "notification-service",
	})
}
