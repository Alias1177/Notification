package templates

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func SendEmail(val string) {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем значения из переменных окружения
	senderEmail := os.Getenv("MAIL")
	secret := os.Getenv("SECRET")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Создание нового письма
	e := email.NewEmail()

	// Отправитель, получатель, тема и текст письма
	e.From = senderEmail // Адрес отправителя
	e.To = []string{val} // Адрес получателя
	e.Subject = "Test Email from Mail.ru"
	e.Text = []byte("This is a plain text email sent via Mail.ru SMTP!")

	// Если нужно отправить HTML-письмо
	e.HTML = []byte("<h1>This is an HTML email sent via Mail.ru SMTP!</h1>")

	// Настройки для подключения к Mail.ru SMTP
	auth := smtp.PlainAuth("", senderEmail, secret, smtpHost)

	// Отправка письма через SMTP-сервер Mail.ru
	err = e.Send(smtpHost+":"+smtpPort, auth)
	if err != nil {
		log.Fatalf("Ошибка при отправке письма: %v", err)
	}

	fmt.Println("Email успешно отправлен!")
}
