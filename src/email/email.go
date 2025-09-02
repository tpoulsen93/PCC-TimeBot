// Package email provides email functionality for sending time cards and notifications.
// It supports sending HTML emails with PDF attachments via SMTP.
package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/src/helpers"
	"github.com/tpoulsen/pcc-timebot/src/timecard"
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
func NewSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     465,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

// connectSMTP establishes a TLS connection to the SMTP server and returns a ready-to-use client.
func connectSMTP(cfg *SMTPConfig, from, to string) (*smtp.Client, error) {
	tlsConfig := &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS connection: %w", err)
	}

	c, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create SMTP client: %w", err)
	}

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

// SendTimeCard sends a time card via email
func SendTimeCard(cfg *SMTPConfig, from, to, name string, body []byte, payday time.Time) error {
	var buf bytes.Buffer

	// Create multipart writer
	writer := multipart.NewWriter(&buf)
	boundary := writer.Boundary()

	// Write email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("TimeBot <%s>", from)
	headers["To"] = to
	headers["Subject"] = fmt.Sprintf("Time Card for payday: %s", payday.Format("2006-01-02"))
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)

	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")

	// Add attachment part
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Transfer-Encoding", "base64")
	header.Set("Content-Disposition", `attachment; filename="TimeCard.txt"`)

	part, err := writer.CreatePart(header)
	if err != nil {
		return fmt.Errorf("failed to create attachment part: %w", err)
	}

	encoder := base64.NewEncoder(base64.StdEncoding, part)
	encoder.Write(body)
	encoder.Close()
	writer.Close()

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
