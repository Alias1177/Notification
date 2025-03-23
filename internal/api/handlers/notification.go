package handlers

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/jordan-wright/email"
)

type CodeResponse struct {
	CodeForFrontEnd string `json:"code_for_frontend"`
}

type CodeRequest struct {
	RequestCode string `json:"request_code"` // Поле для email
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendCodeForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req CodeRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	if req.RequestCode == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Email is required"})
		return
	}

	// Генерация кода
	code := generateCode()

	// Отправка кода на почту
	err = sendPasswordResetEmail(req.RequestCode, code)
	if err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to send password reset email"})
		return
	}

	// Отправляем успешный статус (без самого кода для безопасности)
	render.JSON(w, r, CodeResponse{
		CodeForFrontEnd: code,
	})
}

// Функция для отправки кода восстановления пароля по email
func sendPasswordResetEmail(recipientEmail, code string) error {
	// Retrieve environment variables
	senderEmail := os.Getenv("MAIL")
	secret := os.Getenv("SECRET")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if senderEmail == "" || secret == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("not all environment variables are set")
	}

	// Get the absolute path to the HTML template using Go's runtime package
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("unable to determine current file path")
	}

	// Navigate from the current file (in handlers package) to the HTML template
	basePath := filepath.Dir(filename)
	// Идем на 3 уровня выше (к корню проекта) и затем в templates/html
	htmlPath := filepath.Join(basePath, "..", "..", "..", "templates", "html", "password_reset.html")

	// Verify the file exists
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		// Try alternative path for Docker or different runtime environments
		altPath := "/app/templates/html/password_reset.html"
		log.Printf("Template not found. Trying alternative path: %s", altPath)

		if _, err := os.Stat(altPath); os.IsNotExist(err) {
			// Для отладки - выводим текущий рабочий каталог
			cwd, _ := os.Getwd()
			log.Printf("Current working directory: %s", cwd)
			return fmt.Errorf("HTML template file not found")
		}

		htmlPath = altPath
	}

	// Read HTML template file
	htmlContent, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		return fmt.Errorf("error reading HTML template file: %v", err)
	}

	// Заменяем плейсхолдер кода в HTML шаблоне
	htmlString := string(htmlContent)
	htmlString = strings.Replace(htmlString, "{{CODE}}", code, -1)

	// Create new email
	e := email.NewEmail()

	// Set sender, recipient, subject and content
	e.From = fmt.Sprintf("Password Reset <%s>", senderEmail)
	e.To = []string{recipientEmail}
	e.Subject = "Password Reset Code"

	// Plain text alternative for email clients that don't support HTML
	e.Text = []byte(fmt.Sprintf("Your password reset code is: %s\nThis code will expire in 10 minutes.", code))

	// Set HTML content from template file with replaced code
	e.HTML = []byte(htmlString)

	// Configure TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	// Server address with port
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Authentication
	auth := smtp.PlainAuth("", senderEmail, secret, smtpHost)

	// Send email via SSL
	err = e.SendWithTLS(addr, auth, tlsConfig)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	log.Printf("Password reset code successfully sent to %s", recipientEmail)
	return nil
}
