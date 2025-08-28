package email

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
}

func NewEmailService() *EmailService {
	
	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPortStr := os.Getenv("EMAIL_PORT")
	smtpPort := 587
	if portNum, err := strconv.Atoi(smtpPortStr); err == nil && portNum > 0 {
		smtpPort = portNum
	}

	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	return &EmailService{
		dialer: dialer,
	}

}

func (s *EmailService) SendEmail(to string, subject string, body string) error {
    smtpUser := os.Getenv("SMTP_USER") 

    m := gomail.NewMessage()

    fromName := "Football Shop" 
    if fromName != "" {
        m.SetAddressHeader("From", smtpUser, fromName)
    } else {
        m.SetHeader("From", smtpUser)
    }

    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    return s.dialer.DialAndSend(m)
}
