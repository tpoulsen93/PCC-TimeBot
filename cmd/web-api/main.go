package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tpoulsen/pcc-timebot/internal/handlers"
	"github.com/tpoulsen/pcc-timebot/internal/middleware"
	"github.com/tpoulsen/pcc-timebot/shared/database"
)

func main() {
	// Initialize database connection
	err := database.Initialize()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "timebot-web-api",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Initialize handlers with database
		db := database.GetDB()
		timecardHandler := handlers.NewTimecardHandler(db)
		employeeHandler := handlers.NewEmployeeHandler(db)

		// Timecard routes
		api.GET("/timecards", timecardHandler.GetTimecards)
		api.POST("/timecards", timecardHandler.CreateTimecard)
		api.GET("/timecards/:id", timecardHandler.GetTimecard)
		api.PUT("/timecards/:id", timecardHandler.UpdateTimecard)
		api.DELETE("/timecards/:id", timecardHandler.DeleteTimecard)

		// Employee routes
		api.GET("/employees", employeeHandler.GetEmployees)
		api.GET("/employees/:id", employeeHandler.GetEmployee)
		api.PUT("/employees/:id", employeeHandler.UpdateEmployee)

		// Reports routes
		api.GET("/reports/payroll", timecardHandler.GetPayrollReport)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting web API server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
