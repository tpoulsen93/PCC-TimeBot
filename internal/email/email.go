// Package email provides email functionality for sending time cards and notifications.
// It supports sending HTML emails with PDF attachments via SMTP.
package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/shared/helpers"
	"github.com/tpoulsen/pcc-timebot/shared/timecard"
)

// SMTPConfig holds SMTP server configuration for sending emails.
type SMTPConfig struct {
	Host     string // SMTP server hostname
	Port     int    // SMTP server port
	Username string // SMTP username
	Password string // SMTP password
}

// NewSMTPConfig creates a new SMTP configuration using environment variables.
// It expects SMTP_USERNAME and SMTP_PASSWORD to be set in the environment.
// Uses port 587 with STARTTLS which works on Heroku (port 465 is blocked).
func NewSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

// connectSMTP establishes a connection to the SMTP server using STARTTLS and returns a ready-to-use client.
func connectSMTP(cfg *SMTPConfig, from, to string) (*smtp.Client, error) {
	// Connect to the SMTP server (plain connection first)
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	if err = c.StartTLS(tlsConfig); err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if err = c.Auth(auth); err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = c.Mail(from); err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to set FROM address: %w", err)
	}

	if err = c.Rcpt(to); err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to set TO address: %w", err)
	}

	return c, nil
}

// SendLoginLink emails a magic-link sign-in URL to an employee.
// The link is single-use and time-limited; this function only delivers it.
func SendLoginLink(cfg *SMTPConfig, from, to, name, loginURL string, expiresMinutes int) error {
	displayName := strings.TrimSpace(name)
	if displayName == "" {
		displayName = "there"
	}

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<style>
body { font-family: Arial, sans-serif; margin: 0; padding: 24px; background: #f3f4f6; color: #111827; }
.card { max-width: 480px; margin: 0 auto; background: #ffffff; border-radius: 12px; padding: 32px; }
.btn { display: inline-block; padding: 14px 24px; background: #4f46e5; color: #ffffff !important;
       text-decoration: none; border-radius: 10px; font-weight: 600; }
.muted { color: #6b7280; font-size: 13px; line-height: 1.5; }
.link { word-break: break-all; color: #4f46e5; font-size: 13px; }
</style>
</head>
<body>
<div class="card">
  <h2>PCC TimeBot sign-in</h2>
  <p>Hi %s,</p>
  <p>Click the button below to sign in. This link expires in %d minutes and can only be used once.</p>
  <p style="margin: 24px 0;"><a class="btn" href="%s">Sign in to PCC TimeBot</a></p>
  <p class="muted">If the button doesn't work, copy and paste this URL into your browser:</p>
  <p class="link">%s</p>
  <p class="muted">If you didn't request this, you can safely ignore this email.</p>
</div>
</body>
</html>`, displayName, expiresMinutes, loginURL, loginURL)

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: PCC TimeBot <%s>\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString("Subject: Your PCC TimeBot sign-in link\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(htmlBody)

	c, err := connectSMTP(cfg, from, to)
	if err != nil {
		return err
	}
	defer c.Close()

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to create message writer: %w", err)
	}
	if _, err = w.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return c.Quit()
}

// SendTimeCard sends a time card via email as HTML
func SendTimeCard(cfg *SMTPConfig, from, to, name string, htmlBody string, payday time.Time) error {
	var buf bytes.Buffer

	// Write email headers
	buf.WriteString(fmt.Sprintf("From: TimeBot <%s>\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: Time Card for payday: %s\r\n", payday.Format("2006-01-02")))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("\r\n")

	// Write HTML body
	buf.WriteString(htmlBody)

	// Connect to SMTP server and send email
	c, err := connectSMTP(cfg, from, to)
	if err != nil {
		return err
	}
	defer c.Close()

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to create message writer: %w", err)
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return c.Quit()
}

// SendPayrollSummary sends a summary of all time cards to the admin
func SendPayrollSummary(cfg *SMTPConfig, from, to string, timeCards map[int]*timecard.TimeCard, startDate, endDate, payday time.Time) error {
	var body strings.Builder

	// Write email headers
	body.WriteString(fmt.Sprintf("From: TimeBot <%s>\n", from))
	body.WriteString(fmt.Sprintf("To: TP <%s>\n", to))
	body.WriteString(fmt.Sprintf("Subject: PCC Payroll totals for payday -> %s\n\n", payday.Format("2006-01-02")))
	body.WriteString(fmt.Sprintf("Pay period: %s  <->  %s\nPayday: %s\n\n",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		payday.Format("2006-01-02")))

	// Collect employee IDs for sorting
	ids := make([]int, 0, len(timeCards))
	for id := range timeCards {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	var hoursSum, costSum float64
	for _, id := range ids {
		tc := timeCards[id]
		if tc == nil {
			continue
		}
		hoursSum += tc.TotalHours
		costSum += tc.TotalHours * helpers.GlobalWage
		body.WriteString(fmt.Sprintf("%s  -->  %.2f\n", tc.Name, tc.TotalHours))
	}

	body.WriteString(fmt.Sprintf("\nTotal Hours  -->  %.2f\n", hoursSum))
	body.WriteString(fmt.Sprintf("Estimated Total Cost  -->  $%.2f\n", costSum))

	// Print the summary to standard output
	fmt.Printf("\nPayroll Summary for %s to %s\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	fmt.Printf("Payday: %s\n\n", payday.Format("2006-01-02"))
	fmt.Printf("%-30s %s\n", "Employee", "Hours")
	fmt.Printf("%s %s\n", strings.Repeat("-", 30), strings.Repeat("-", 10))

	for _, id := range ids {
		tc := timeCards[id]
		if tc == nil {
			continue
		}
		fmt.Printf("%-30s %7.2f\n", tc.Name, tc.TotalHours)
	}
	fmt.Printf("\n%-30s %7.2f\n", "Total Hours:", hoursSum)
	fmt.Printf("%-30s $%7.2f\n\n", "Estimated Total Cost:", costSum)

	// Connect to SMTP server and send email
	c, err := connectSMTP(cfg, from, to)
	if err != nil {
		return err
	}
	defer c.Close()

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to create message writer: %w", err)
	}

	_, err = w.Write([]byte(body.String()))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return c.Quit()
}
