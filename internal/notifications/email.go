package notifications

import (
	"fmt"
	"net/smtp"
)

func (ns *NotificationService) SendEmail(to, subject, body string) error {
	// Set up authentication information.
	from := "youremail@example.com"
	password := "yourpassword"
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Compose the email message.
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// Connect to the server, authenticate, and send the email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}
