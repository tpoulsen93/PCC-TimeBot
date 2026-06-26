package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tpoulsen/pcc-timebot/internal/auth"
	"github.com/tpoulsen/pcc-timebot/internal/email"
	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
)

// AuthHandler handles passwordless magic-link authentication.
type AuthHandler struct{}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type requestLinkBody struct {
	Email string `json:"email"`
}

// RequestLink generates a magic-link token for the given email and emails it.
// To avoid account enumeration, it always returns 200 regardless of whether
// the email matches an employee.
//
// In development (APP_ENV=dev), the email step is skipped entirely:
// a session is created immediately and the cookie is set on the response so
// the browser is logged in without any email round-trip.
func (h *AuthHandler) RequestLink(c *gin.Context) {
	var body requestLinkBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	emailAddr := strings.TrimSpace(strings.ToLower(body.Email))
	genericResponse := gin.H{"message": "If that email is registered, a sign-in link has been sent."}

	if emailAddr == "" {
		c.JSON(http.StatusOK, genericResponse)
		return
	}

	employee, err := database.GetEmployeeByEmail(emailAddr)
	if err != nil {
		// Log server-side, but don't leak detail to the client.
		fmt.Printf("RequestLink lookup error: %v\n", err)
		c.JSON(http.StatusOK, genericResponse)
		return
	}
	if employee == nil || employee.Email == "" {
		c.JSON(http.StatusOK, genericResponse)
		return
	}

	// Dev shortcut: skip the email and log in immediately.
	if strings.ToLower(os.Getenv("APP_ENV")) == "dev" {
		sessionID, _, err := auth.CreateSession(employee.ID)
		if err != nil {
			fmt.Printf("Dev login session error: %v\n", err)
			c.JSON(http.StatusOK, genericResponse)
			return
		}
		auth.SetSessionCookie(c, sessionID, int(auth.SessionTTL.Seconds()))
		c.JSON(http.StatusOK, gin.H{"dev": true})
		return
	}

	if err := h.sendLoginLink(employee); err != nil {
		fmt.Printf("Failed to send login link to %s: %v\n", employee.Email, err)
	}

	c.JSON(http.StatusOK, genericResponse)
}

func (h *AuthHandler) sendLoginLink(employee *database.Employee) error {
	token, err := auth.CreateLoginToken(employee.ID)
	if err != nil {
		return err
	}

	loginURL := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", appBaseURL(), token)

	cfg := email.NewSMTPConfig()
	if cfg.Username == "" || cfg.Password == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	name := helpers.Title.String(employee.FirstName)
	minutes := int(auth.LoginTokenTTL.Minutes())
	return email.SendLoginLink(cfg, cfg.Username, employee.Email, name, loginURL, minutes)
}

// Verify consumes a magic-link token, creates a session, sets the session
// cookie, and redirects to the SPA root.
func (h *AuthHandler) Verify(c *gin.Context) {
	token := c.Query("token")

	employeeID, err := auth.ConsumeLoginToken(token)
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=invalid_link")
		return
	}

	sessionID, expiresAt, err := auth.CreateSession(employeeID)
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=server_error")
		return
	}

	maxAge := int(auth.SessionTTL.Seconds())
	_ = expiresAt
	auth.SetSessionCookie(c, sessionID, maxAge)
	c.Redirect(http.StatusFound, "/")
}

// Logout deletes the current session and clears the cookie.
func (h *AuthHandler) Logout(c *gin.Context) {
	if sessionID, err := c.Cookie(auth.SessionCookieName); err == nil && sessionID != "" {
		_ = database.DeleteSession(sessionID)
	}
	auth.ClearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Me returns the currently authenticated employee.
func (h *AuthHandler) Me(c *gin.Context) {
	employee, ok := auth.CurrentEmployee(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	c.JSON(http.StatusOK, employeeToJSON(employee))
}

// appBaseURL returns the public base URL used to build magic links.
func appBaseURL() string {
	if url := strings.TrimRight(os.Getenv("APP_BASE_URL"), "/"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

// employeeToJSON converts an employee into the public JSON shape used by the API.
func employeeToJSON(e *database.Employee) gin.H {
	return gin.H{
		"id":           e.ID,
		"firstName":    helpers.Title.String(e.FirstName),
		"lastName":     helpers.Title.String(e.LastName),
		"email":        e.Email,
		"phone":        e.Phone,
		"supervisorId": e.SupervisorID,
		"isAdmin":      e.IsAdmin,
	}
}
