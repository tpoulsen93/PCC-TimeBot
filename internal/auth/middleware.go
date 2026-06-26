package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tpoulsen/pcc-timebot/shared/database"
)

const contextEmployeeKey = "employee"

// SetSessionCookie writes the session cookie on the response. The cookie is
// HttpOnly (not readable by JS) and Secure in production, with SameSite=Lax to
// allow top-level navigation from the magic-link email while blocking CSRF on
// cross-site POSTs.
func SetSessionCookie(c *gin.Context, sessionID string, maxAgeSeconds int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		SessionCookieName,
		sessionID,
		maxAgeSeconds,
		"/",
		"",         // default to the request host
		isSecure(), // Secure flag (HTTPS only) in production
		true,       // HttpOnly
	)
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SessionCookieName, "", -1, "/", "", isSecure(), true)
}

// isSecure reports whether cookies should carry the Secure flag. It is enabled
// unless APP_ENV is explicitly set to "dev".
func isSecure() bool {
	return strings.ToLower(os.Getenv("APP_ENV")) != "dev"
}

// RequireAuth is gin middleware that rejects unauthenticated requests. On
// success it stores the authenticated employee in the request context.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(SessionCookieName)
		if err != nil || sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		employeeID, err := database.GetSessionEmployeeID(sessionID, time.Now().Add(SessionTTL))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "session lookup failed"})
			return
		}
		if employeeID == 0 {
			ClearSessionCookie(c)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
			return
		}

		// Refresh the browser cookie's Max-Age on every authenticated request
		// so the expiry slides forward silently while the user is active.
		SetSessionCookie(c, sessionID, int(SessionTTL.Seconds()))

		employee, err := database.GetEmployee(employeeID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
			return
		}

		c.Set(contextEmployeeKey, employee)
		c.Next()
	}
}

// RequireAdmin is gin middleware that requires an authenticated admin. It must
// be chained after RequireAuth.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		employee, ok := CurrentEmployee(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}
		if !employee.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}

// CurrentEmployee returns the authenticated employee stored in the context.
func CurrentEmployee(c *gin.Context) (*database.Employee, bool) {
	value, exists := c.Get(contextEmployeeKey)
	if !exists {
		return nil, false
	}
	employee, ok := value.(*database.Employee)
	return employee, ok
}
