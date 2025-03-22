package templates

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jordan-wright/email"
)

func SendEmail(val string) {
	// Retrieve environment variables
	senderEmail := os.Getenv("MAIL")
	secret := os.Getenv("SECRET")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if senderEmail == "" || secret == "" || smtpHost == "" || smtpPort == "" {
		log.Printf("Error: not all environment variables are set. MAIL: %s, SMTP_HOST: %s, SMTP_PORT: %s",
			senderEmail, smtpHost, smtpPort)
		return
	}

	// Get the absolute path to the HTML template using Go's runtime package
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Printf("Error: Unable to determine current file path")
		return
	}

	// Navigate from the current file (in templates package) to the HTML template
	basePath := filepath.Dir(filename)
	htmlPath := filepath.Join(basePath, "html", "registration_confirmation.html")

	// Verify the file exists
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		// Try alternative path for Docker or different runtime environments
		altPath := "/app/templates/html/registration_confirmation.html"
		log.Printf("Template not found. Trying alternative path: %s", altPath)

		if _, err := os.Stat(altPath); os.IsNotExist(err) {
			log.Printf("Error: HTML template file not found at either location")
			return
		}

		htmlPath = altPath
	}

	// Read HTML template file
	htmlContent, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		log.Printf("Error reading HTML template file from %s: %v", htmlPath, err)
		return
	}

	// Create new email
	e := email.NewEmail()

	// Set sender, recipient, subject and content
	e.From = fmt.Sprintf("Four-X Registration <%s>", senderEmail)
	e.To = []string{val}
	e.Subject = "Welcome to Four-X: Registration Successful!"

	// Plain text alternative for email clients that don't support HTML
	e.Text = []byte("Thank you for registering with four-x.com. Your registration has been successfully completed.")

	// Set HTML content from template file
	e.HTML = htmlContent

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
		log.Printf("Error sending email: %v", err)
		return
	}

	log.Printf("Registration confirmation email successfully sent to %s", val)
}
