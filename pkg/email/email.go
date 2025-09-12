package email

import (
	"fmt"
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

func (s *EmailService) SendEmailChangePassword(toEmail string, resetToken string) error {
	m := gomail.NewMessage()

	from := m.FormatAddress(os.Getenv("SMTP_USER"), "Your App Name")
	m.SetHeader("From", from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Reset your password")

	resetLink := fmt.Sprintf("https://your-frontend.com/reset-password?token=%s", resetToken)

	body := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				margin: 0;
				padding: 0;
			}
			.container {
				max-width: 600px;
				margin: 20px auto;
				background: #ffffff;
				padding: 30px;
				border-radius: 8px;
				box-shadow: 0 2px 6px rgba(0,0,0,0.1);
			}
			.header {
				font-size: 22px;
				font-weight: bold;
				margin-bottom: 20px;
				color: #333333;
			}
			.button {
				display: inline-block;
				padding: 12px 20px;
				margin-top: 20px;
				background-color: #007bff;
				color: #ffffff;
				text-decoration: none;
				border-radius: 5px;
				font-weight: bold;
			}
			.footer {
				margin-top: 30px;
				font-size: 12px;
				color: #777777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">Reset Your Password</div>
			<p>Hello,</p>
			<p>You recently requested to reset your password. Click the button below to set a new one:</p>
			<a href="%s" class="button">Reset Password</a>
			<p>If you didn’t request this, please ignore this email.</p>
			<div class="footer">
				<p>This link will expire in 15 minutes.</p>
				<p>© 2025 Your App Name</p>
			</div>
		</div>
	</body>
	</html>
	`, resetLink)

	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
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
