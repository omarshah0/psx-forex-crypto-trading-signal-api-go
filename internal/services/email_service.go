package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

// EmailSender is the interface for sending emails
type EmailSender interface {
	SendVerificationEmail(email, name, token string) error
	SendPasswordResetEmail(email, name, token string) error
	SendPasswordChangedEmail(email, name string) error
	SendWelcomeEmail(email, name string) error
	SendSubscriptionConfirmation(email, name string, subscriptions []models.SubscriptionWithPackage, totalAmount float64) error
}

// EmailService wraps the email sender implementation
type EmailService struct {
	sender EmailSender
}

// NewEmailService creates a new email service with the appropriate sender
func NewEmailService(enabled bool, provider, frontendURL, fromAddress, fromName, resendAPIKey, smtpHost, smtpUsername, smtpPassword string, smtpPort int) *EmailService {
	if !enabled {
		return &EmailService{
			sender: NewMockEmailService(frontendURL, fromAddress, fromName),
		}
	}

	switch provider {
	case "resend":
		return &EmailService{
			sender: NewResendEmailService(resendAPIKey, frontendURL, fromAddress, fromName),
		}
	case "smtp":
		return &EmailService{
			sender: NewSMTPEmailService(smtpHost, smtpPort, smtpUsername, smtpPassword, frontendURL, fromAddress, fromName),
		}
	default:
		log.Printf("Unknown email provider '%s', using mock email service", provider)
		return &EmailService{
			sender: NewMockEmailService(frontendURL, fromAddress, fromName),
		}
	}
}

func (s *EmailService) SendVerificationEmail(email, name, token string) error {
	return s.sender.SendVerificationEmail(email, name, token)
}

func (s *EmailService) SendPasswordResetEmail(email, name, token string) error {
	return s.sender.SendPasswordResetEmail(email, name, token)
}

func (s *EmailService) SendPasswordChangedEmail(email, name string) error {
	return s.sender.SendPasswordChangedEmail(email, name)
}

func (s *EmailService) SendWelcomeEmail(email, name string) error {
	return s.sender.SendWelcomeEmail(email, name)
}

func (s *EmailService) SendSubscriptionConfirmation(email, name string, subscriptions []models.SubscriptionWithPackage, totalAmount float64) error {
	return s.sender.SendSubscriptionConfirmation(email, name, subscriptions, totalAmount)
}

// MockEmailService simulates email sending by logging
type MockEmailService struct {
	frontendURL      string
	verificationURL  string
	resetPasswordURL string
	fromAddress      string
	fromName         string
}

func NewMockEmailService(frontendURL, fromAddress, fromName string) *MockEmailService {
	return &MockEmailService{
		frontendURL:      frontendURL,
		verificationURL:  frontendURL + "/verify-email",
		resetPasswordURL: frontendURL + "/reset-password",
		fromAddress:      fromAddress,
		fromName:         fromName,
	}
}

func (s *MockEmailService) SendVerificationEmail(email, name, token string) error {
	verificationLink := fmt.Sprintf("%s?token=%s", s.verificationURL, token)
	log.Printf("[EMAIL SIMULATION] Verification email to %s\n", email)
	log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
	log.Printf("[EMAIL SIMULATION] Link: %s\n", verificationLink)
	log.Printf("[EMAIL SIMULATION] Token: %s\n", token)
	return nil
}

func (s *MockEmailService) SendPasswordResetEmail(email, name, token string) error {
	resetLink := fmt.Sprintf("%s?token=%s", s.resetPasswordURL, token)
	log.Printf("[EMAIL SIMULATION] Password reset email to %s\n", email)
	log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
	log.Printf("[EMAIL SIMULATION] Link: %s\n", resetLink)
	log.Printf("[EMAIL SIMULATION] Token: %s\n", token)
	return nil
}

func (s *MockEmailService) SendPasswordChangedEmail(email, name string) error {
	log.Printf("[EMAIL SIMULATION] Password changed notification to %s\n", email)
	log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
	return nil
}

func (s *MockEmailService) SendWelcomeEmail(email, name string) error {
	log.Printf("[EMAIL SIMULATION] Welcome email to %s\n", email)
	log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
	return nil
}

func (s *MockEmailService) SendSubscriptionConfirmation(email, name string, subscriptions []models.SubscriptionWithPackage, totalAmount float64) error {
	log.Printf("[EMAIL SIMULATION] Subscription confirmation to %s\n", email)
	log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
	log.Printf("[EMAIL SIMULATION] Total Amount: $%.2f\n", totalAmount)
	log.Printf("[EMAIL SIMULATION] Subscriptions: %d packages\n", len(subscriptions))
	return nil
}

// ResendEmailService sends emails using Resend API
type ResendEmailService struct {
	apiKey           string
	frontendURL      string
	verificationURL  string
	resetPasswordURL string
	fromAddress      string
	fromName         string
	httpClient       *http.Client
}

func NewResendEmailService(apiKey, frontendURL, fromAddress, fromName string) *ResendEmailService {
	return &ResendEmailService{
		apiKey:           apiKey,
		frontendURL:      frontendURL,
		verificationURL:  frontendURL + "/verify-email",
		resetPasswordURL: frontendURL + "/reset-password",
		fromAddress:      fromAddress,
		fromName:         fromName,
		httpClient:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *ResendEmailService) sendEmail(to, subject, htmlBody string) error {
	reqBody := map[string]interface{}{
		"from":    fmt.Sprintf("%s <%s>", s.fromName, s.fromAddress),
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resend API returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *ResendEmailService) SendVerificationEmail(email, name, token string) error {
	verificationLink := fmt.Sprintf("%s?token=%s", s.verificationURL, token)
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Thank you for signing up! Please verify your email address by clicking the link below:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>This link will expire in 24 hours.</p>
		<p>If you didn't create an account, please ignore this email.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, verificationLink, s.fromName)
	return s.sendEmail(email, subject, body)
}

func (s *ResendEmailService) SendPasswordResetEmail(email, name, token string) error {
	resetLink := fmt.Sprintf("%s?token=%s", s.resetPasswordURL, token)
	subject := "Reset Your Password"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>We received a request to reset your password. Click the link below to reset it:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you didn't request a password reset, please ignore this email or contact support if you have concerns.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, resetLink, s.fromName)
	return s.sendEmail(email, subject, body)
}

func (s *ResendEmailService) SendPasswordChangedEmail(email, name string) error {
	subject := "Your Password Was Changed"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>This is a confirmation that your password was successfully changed.</p>
		<p>If you didn't make this change, please contact support immediately.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, s.fromName)
	return s.sendEmail(email, subject, body)
}

func (s *ResendEmailService) SendWelcomeEmail(email, name string) error {
	subject := "Welcome!"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Welcome to %s! We're excited to have you on board.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, s.fromName, s.fromName)
	return s.sendEmail(email, subject, body)
}

func (s *ResendEmailService) SendSubscriptionConfirmation(email, name string, subscriptions []models.SubscriptionWithPackage, totalAmount float64) error {
	subject := "Subscription Confirmation"

	var packagesHTML string
	for _, sub := range subscriptions {
		if sub.Package != nil {
			packagesHTML += fmt.Sprintf("<li>%s - $%.2f (expires: %s)</li>",
				sub.Package.Name, sub.PricePaid, sub.ExpiresAt.Format("January 2, 2006"))
		}
	}

	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Thank you for subscribing! Your subscription has been confirmed.</p>
		<h3>Subscription Details:</h3>
		<ul>%s</ul>
		<p><strong>Total Amount: $%.2f</strong></p>
		<p>You now have access to all signals in your subscribed packages.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, packagesHTML, totalAmount, s.fromName)
	return s.sendEmail(email, subject, body)
}

// SMTPEmailService sends emails using SMTP
type SMTPEmailService struct {
	host             string
	port             int
	username         string
	password         string
	frontendURL      string
	verificationURL  string
	resetPasswordURL string
	fromAddress      string
	fromName         string
}

func NewSMTPEmailService(host string, port int, username, password, frontendURL, fromAddress, fromName string) *SMTPEmailService {
	return &SMTPEmailService{
		host:             host,
		port:             port,
		username:         username,
		password:         password,
		frontendURL:      frontendURL,
		verificationURL:  frontendURL + "/verify-email",
		resetPasswordURL: frontendURL + "/reset-password",
		fromAddress:      fromAddress,
		fromName:         fromName,
	}
}

func (s *SMTPEmailService) sendEmail(to, subject, body string) error {
	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromAddress)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", from, to, subject, body))

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// For TLS connections
	tlsConfig := &tls.Config{
		ServerName: s.host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		// Try without TLS (for port 587 with STARTTLS)
		return smtp.SendMail(addr, auth, s.fromAddress, []string{to}, msg)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err := client.Mail(s.fromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

func (s *SMTPEmailService) SendVerificationEmail(email, name, token string) error {
	verificationLink := fmt.Sprintf("%s?token=%s", s.verificationURL, token)
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Thank you for signing up! Please verify your email address by clicking the link below:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>This link will expire in 24 hours.</p>
		<p>If you didn't create an account, please ignore this email.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, verificationLink, s.fromName)
	return s.sendEmail(email, subject, body)
}

func (s *SMTPEmailService) SendPasswordResetEmail(email, name, token string) error {
	resetLink := fmt.Sprintf("%s?token=%s", s.resetPasswordURL, token)
	subject := "Reset Your Password"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>We received a request to reset your password. Click the link below to reset it:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you didn't request a password reset, please ignore this email or contact support if you have concerns.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, resetLink, s.fromName)
	return s.sendEmail(email, subject, body)
}

func (s *SMTPEmailService) SendPasswordChangedEmail(email, name string) error {
	subject := "Your Password Was Changed"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>This is a confirmation that your password was successfully changed.</p>
		<p>If you didn't make this change, please contact support immediately.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, s.fromName)
	return s.sendEmail(email, subject, body)
}

func (s *SMTPEmailService) SendWelcomeEmail(email, name string) error {
	subject := "Welcome!"
	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Welcome to %s! We're excited to have you on board.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, s.fromName, s.fromName)
	return s.sendEmail(email, subject, body)
}

func (s *SMTPEmailService) SendSubscriptionConfirmation(email, name string, subscriptions []models.SubscriptionWithPackage, totalAmount float64) error {
	subject := "Subscription Confirmation"

	var packagesHTML string
	for _, sub := range subscriptions {
		if sub.Package != nil {
			packagesHTML += fmt.Sprintf("<li>%s - $%.2f (expires: %s)</li>",
				sub.Package.Name, sub.PricePaid, sub.ExpiresAt.Format("January 2, 2006"))
		}
	}

	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Thank you for subscribing! Your subscription has been confirmed.</p>
		<h3>Subscription Details:</h3>
		<ul>%s</ul>
		<p><strong>Total Amount: $%.2f</strong></p>
		<p>You now have access to all signals in your subscribed packages.</p>
		<p>Best regards,<br>%s Team</p>
	`, name, packagesHTML, totalAmount, s.fromName)
	return s.sendEmail(email, subject, body)
}
