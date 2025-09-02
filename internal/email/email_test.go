package email

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSMTPConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "testpassword")
	defer func() {
		os.Unsetenv("SMTP_USERNAME")
		os.Unsetenv("SMTP_PASSWORD")
	}()

	config := NewSMTPConfig()

	assert.Equal(t, "smtp.gmail.com", config.Host)
	assert.Equal(t, 465, config.Port)
	assert.Equal(t, "test@example.com", config.Username)
	assert.Equal(t, "testpassword", config.Password)
}

func TestNewSMTPConfig_EmptyEnvVars(t *testing.T) {
	// Ensure environment variables are not set
	os.Unsetenv("SMTP_USERNAME")
	os.Unsetenv("SMTP_PASSWORD")

	config := NewSMTPConfig()

	assert.Equal(t, "smtp.gmail.com", config.Host)
	assert.Equal(t, 465, config.Port)
	assert.Empty(t, config.Username)
	assert.Empty(t, config.Password)
}

func TestSMTPConfig_Fields(t *testing.T) {
	config := &SMTPConfig{
		Host:     "custom.smtp.com",
		Port:     587,
		Username: "user@test.com",
		Password: "secret",
	}

	assert.Equal(t, "custom.smtp.com", config.Host)
	assert.Equal(t, 587, config.Port)
	assert.Equal(t, "user@test.com", config.Username)
	assert.Equal(t, "secret", config.Password)
}

// Note: connectSMTP and SendTimeCard/SendPayrollSummary tests would require
// mocking SMTP connections, which is complex and typically done with integration tests.
// For unit tests, we focus on the testable parts like configuration and data preparation.
