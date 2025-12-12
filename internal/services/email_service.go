package services

import (
	"fmt"
	"log"
)

type EmailService struct {
	enabled           bool
	frontendURL       string
	verificationURL   string
	resetPasswordURL  string
	fromAddress       string
	fromName          string
}

func NewEmailService(enabled bool, frontendURL, fromAddress, fromName string) *EmailService {
	return &EmailService{
		enabled:          enabled,
		frontendURL:      frontendURL,
		verificationURL:  frontendURL + "/verify-email",
		resetPasswordURL: frontendURL + "/reset-password",
		fromAddress:      fromAddress,
		fromName:         fromName,
	}
}

// SendVerificationEmail sends an email verification link
func (s *EmailService) SendVerificationEmail(email, name, token string) error {
	verificationLink := fmt.Sprintf("%s?token=%s", s.verificationURL, token)
	
	if !s.enabled {
		log.Printf("[EMAIL SIMULATION] Verification email to %s\n", email)
		log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
		log.Printf("[EMAIL SIMULATION] Link: %s\n", verificationLink)
		log.Printf("[EMAIL SIMULATION] Token: %s\n", token)
		return nil
	}
	
	// TODO: Implement actual email sending
	// Use SendGrid, Mailgun, AWS SES, or SMTP
	// Example structure:
	/*
		subject := "Verify Your Email Address"
		body := fmt.Sprintf(`
			Hi %s,
			
			Thank you for signing up! Please verify your email address by clicking the link below:
			
			%s
			
			This link will expire in 24 hours.
			
			If you didn't create an account, please ignore this email.
			
			Best regards,
			%s Team
		`, name, verificationLink, s.fromName)
		
		return s.sendEmail(email, subject, body)
	*/
	
	return nil
}

// SendPasswordResetEmail sends a password reset link
func (s *EmailService) SendPasswordResetEmail(email, name, token string) error {
	resetLink := fmt.Sprintf("%s?token=%s", s.resetPasswordURL, token)
	
	if !s.enabled {
		log.Printf("[EMAIL SIMULATION] Password reset email to %s\n", email)
		log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
		log.Printf("[EMAIL SIMULATION] Link: %s\n", resetLink)
		log.Printf("[EMAIL SIMULATION] Token: %s\n", token)
		return nil
	}
	
	// TODO: Implement actual email sending
	// Example structure:
	/*
		subject := "Reset Your Password"
		body := fmt.Sprintf(`
			Hi %s,
			
			We received a request to reset your password. Click the link below to reset it:
			
			%s
			
			This link will expire in 1 hour.
			
			If you didn't request a password reset, please ignore this email or contact support if you have concerns.
			
			Best regards,
			%s Team
		`, name, resetLink, s.fromName)
		
		return s.sendEmail(email, subject, body)
	*/
	
	return nil
}

// SendPasswordChangedEmail sends a notification that password was changed
func (s *EmailService) SendPasswordChangedEmail(email, name string) error {
	if !s.enabled {
		log.Printf("[EMAIL SIMULATION] Password changed notification to %s\n", email)
		log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
		return nil
	}
	
	// TODO: Implement actual email sending
	// Example structure:
	/*
		subject := "Your Password Was Changed"
		body := fmt.Sprintf(`
			Hi %s,
			
			This is a confirmation that your password was successfully changed.
			
			If you didn't make this change, please contact support immediately.
			
			Best regards,
			%s Team
		`, name, s.fromName)
		
		return s.sendEmail(email, subject, body)
	*/
	
	return nil
}

// SendWelcomeEmail sends a welcome email (optional)
func (s *EmailService) SendWelcomeEmail(email, name string) error {
	if !s.enabled {
		log.Printf("[EMAIL SIMULATION] Welcome email to %s\n", email)
		log.Printf("[EMAIL SIMULATION] Name: %s\n", name)
		return nil
	}
	
	// TODO: Implement actual email sending
	return nil
}

// TODO: Implement actual email sending method
// func (s *EmailService) sendEmail(to, subject, body string) error {
//     // Implement based on your email provider
//     // SendGrid, Mailgun, AWS SES, SMTP, etc.
//     return nil
// }

