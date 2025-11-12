package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

// Client represents an email client
type Client interface {
	Send(to, subject, body string, data any, isSandbox bool) (int, error)
}

// EmailData contains common data for email templates
type EmailData struct {
	Username      string
	ActivationURL string
	ExpiryTime    time.Time
	AppName       string
}

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// SMTPClient implements email sending via SMTP (Gmail, etc.)
type SMTPClient struct {
	config SMTPConfig
}

// NewSMTPClient creates a new SMTP email client
func NewSMTPClient(config SMTPConfig) *SMTPClient {
	return &SMTPClient{
		config: config,
	}
}

// Send sends an email using SMTP
func (c *SMTPClient) Send(to, subject, templateName string, data any, isSandbox bool) (int, error) {
	// Parse and execute template
	tmpl, err := template.ParseFS(templates, fmt.Sprintf("templates/%s.html", templateName))
	if err != nil {
		return 0, fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return 0, fmt.Errorf("failed to execute template: %w", err)
	}

	// Build email message
	message := c.buildMessage(to, subject, body.String())

	// Setup authentication
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	// Send email
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	err = smtp.SendMail(addr, auth, c.config.From, []string{to}, []byte(message))
	if err != nil {
		return 0, fmt.Errorf("failed to send email: %w", err)
	}

	return 200, nil
}

// buildMessage constructs the email message with headers
func (c *SMTPClient) buildMessage(to, subject, body string) string {
	message := fmt.Sprintf("From: %s\r\n", c.config.From)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	message += "\r\n"
	message += body
	return message
}
