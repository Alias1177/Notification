package email

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jordan-wright/email"
)

// EmailService отвечает за отправку email уведомлений
type EmailService struct {
	senderEmail string
	secret      string
	smtpHost    string
	smtpPort    string
}

// NewEmailService создает новый экземпляр EmailService
func NewEmailService() *EmailService {
	return &EmailService{
		senderEmail: os.Getenv("MAIL"),
		secret:      os.Getenv("SECRET"),
		smtpHost:    os.Getenv("SMTP_HOST"),
		smtpPort:    os.Getenv("SMTP_PORT"),
	}
}

// SendPasswordResetEmail отправляет код восстановления пароля
func (e *EmailService) SendPasswordResetEmail(recipientEmail, code string) error {
	// Валидация конфигурации
	if e.senderEmail == "" || e.secret == "" || e.smtpHost == "" || e.smtpPort == "" {
		return fmt.Errorf("email configuration is incomplete")
	}

	// Читаем HTML шаблон
	htmlContent, err := e.readHTMLTemplate("password_reset.html")
	if err != nil {
		return fmt.Errorf("failed to read password reset template: %w", err)
	}

	// Заменяем плейсхолдер кода
	htmlString := strings.Replace(string(htmlContent), "{{CODE}}", code, -1)

	// Создаем email
	emailMsg := email.NewEmail()
	emailMsg.From = fmt.Sprintf("Password Reset <%s>", e.senderEmail)
	emailMsg.To = []string{recipientEmail}
	emailMsg.Subject = "Password Reset Code"
	emailMsg.Text = []byte(fmt.Sprintf("Your password reset code is: %s\nThis code will expire in 10 minutes.", code))
	emailMsg.HTML = []byte(htmlString)

	// Отправляем email
	if err := e.sendEmail(emailMsg); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	log.Printf("Password reset code successfully sent to %s", recipientEmail)
	return nil
}

// SendRegistrationEmail отправляет письмо подтверждения регистрации
func (e *EmailService) SendRegistrationEmail(recipientEmail string) error {
	// Валидация конфигурации
	if e.senderEmail == "" || e.secret == "" || e.smtpHost == "" || e.smtpPort == "" {
		return fmt.Errorf("email configuration is incomplete")
	}

	// Читаем HTML шаблон
	htmlContent, err := e.readHTMLTemplate("registration_confirmation.html")
	if err != nil {
		return fmt.Errorf("failed to read registration template: %w", err)
	}

	// Создаем email
	emailMsg := email.NewEmail()
	emailMsg.From = fmt.Sprintf("Four-X Registration <%s>", e.senderEmail)
	emailMsg.To = []string{recipientEmail}
	emailMsg.Subject = "Welcome to Four-X: Registration Successful!"
	emailMsg.Text = []byte("Thank you for registering with four-x.com. Your registration has been successfully completed.")
	emailMsg.HTML = htmlContent

	// Отправляем email
	if err := e.sendEmail(emailMsg); err != nil {
		return fmt.Errorf("failed to send registration email: %w", err)
	}

	log.Printf("Registration confirmation email successfully sent to %s", recipientEmail)
	return nil
}

// readHTMLTemplate читает HTML шаблон из файла
func (e *EmailService) readHTMLTemplate(templateName string) ([]byte, error) {
	// Получаем путь к текущему файлу
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("unable to determine current file path")
	}

	// Строим путь к шаблону
	basePath := filepath.Dir(filename)
	htmlPath := filepath.Join(basePath, "..", "..", "..", "templates", "html", templateName)

	// Проверяем существование файла
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		// Альтернативный путь для Docker
		altPath := fmt.Sprintf("/app/templates/html/%s", templateName)
		if _, err := os.Stat(altPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("HTML template file not found: %s", templateName)
		}
		htmlPath = altPath
	}

	// Читаем файл
	htmlContent, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		return nil, fmt.Errorf("error reading HTML template file: %w", err)
	}

	return htmlContent, nil
}

// sendEmail отправляет email через SMTP
func (e *EmailService) sendEmail(emailMsg *email.Email) error {
	// Настройка TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         e.smtpHost,
	}

	// Адрес сервера
	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)

	// Аутентификация
	auth := smtp.PlainAuth("", e.senderEmail, e.secret, e.smtpHost)

	// Отправка через SSL
	return emailMsg.SendWithTLS(addr, auth, tlsConfig)
}
