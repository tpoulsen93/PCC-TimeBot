package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tpoulsen/pcc-timebot/internal/auth"
	"github.com/tpoulsen/pcc-timebot/internal/handlers"
	"github.com/tpoulsen/pcc-timebot/internal/middleware"
	"github.com/tpoulsen/pcc-timebot/shared/database"
	web "github.com/tpoulsen/pcc-timebot/web"
)

func main() {
	// Load .env if present (local dev). In production (Heroku) the file won't
	// exist and env vars are injected by the platform — the error is ignored.
	_ = godotenv.Overload()

	// Initialize database connection
	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := gin.Default()

	// Same-origin SPA + API: lock CORS down instead of allowing all origins.
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "timebot-web-api",
		})
	})

	authHandler := handlers.NewAuthHandler()
	timecardHandler := handlers.NewTimecardHandler()
	adminHandler := handlers.NewAdminHandler()

	api := r.Group("/api/v1")
	{
		// Public auth endpoints (no session required).
		api.POST("/auth/request-link", authHandler.RequestLink)
		api.GET("/auth/verify", authHandler.Verify)
		api.POST("/auth/logout", authHandler.Logout)

		// Authenticated endpoints.
		authed := api.Group("")
		authed.Use(auth.RequireAuth())
		{
			authed.GET("/me", authHandler.Me)
			authed.POST("/timecards", timecardHandler.SubmitHours)
			authed.GET("/timecards/history", timecardHandler.GetHistory)
			authed.GET("/timecards/summary", timecardHandler.GetWeeklySummary)
		}

		// Admin-only endpoints.
		adminGroup := api.Group("/admin")
		adminGroup.Use(auth.RequireAuth(), auth.RequireAdmin())
		{
			adminGroup.GET("/employees", adminHandler.ListEmployees)
			adminGroup.POST("/employees", adminHandler.CreateEmployee)
			adminGroup.PUT("/employees/:id", adminHandler.UpdateEmployee)
			adminGroup.GET("/timecards", adminHandler.GetAllTimecards)
			adminGroup.POST("/timecards/send", adminHandler.SendTimecards)
		}
	}

	// Serve the embedded React SPA for everything else.
	registerSPA(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting web API server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// registerSPA serves the embedded single-page application. Real asset paths are
// served directly; any unknown non-API path falls back to index.html so that
// client-side routing works on deep links and refreshes.
func registerSPA(r *gin.Engine) {
	dist, err := fs.Sub(web.DistFS(), "app/dist")
	if err != nil {
		log.Fatalf("Failed to load embedded SPA: %v", err)
	}
	fileServer := http.FileServer(http.FS(dist))

	indexHTML, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		log.Fatalf("Failed to read embedded index.html: %v", err)
	}

	serveIndex := func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	}

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API routes must never fall through to the SPA.
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// Serve a real static asset if it exists; otherwise serve the SPA shell.
		trimmed := strings.TrimPrefix(path, "/")
		if trimmed == "" {
			serveIndex(c)
			return
		}
		if f, err := dist.Open(trimmed); err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}
		serveIndex(c)
	})
}
