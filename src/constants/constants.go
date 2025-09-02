// Package constants defines application-wide constants for the PCC-TimeBot application.
package constants

const (
	// PayrollConstants
	DefaultWage        = 30.0 // Default hourly wage rate
	PaydayOffsetDays   = 12   // Number of days after pay period end for payday
	HoursDecimalPlaces = 2    // Number of decimal places for hour calculations

	// Date formats
	DateFormat = "2006-01-02" // Standard date format used throughout the application

	// SMTP Configuration
	SMTPHost = "smtp.gmail.com"
	SMTPPort = 465

	// Environment variable names
	EnvDatabaseURL  = "DATABASE_URL"
	EnvSMTPUsername = "SMTP_USERNAME"
	EnvSMTPPassword = "SMTP_PASSWORD"
)
