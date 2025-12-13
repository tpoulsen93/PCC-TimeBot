package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tpoulsen/pcc-timebot/internal/admin"
	"github.com/tpoulsen/pcc-timebot/shared/database"
)

func main() {
	fmt.Println("Starting timecard scheduler...")

	// Initialize database connection
	if err := database.Initialize(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Calculate last week's Monday-Sunday
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Printf("Failed to load timezone: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().In(loc)
	
	// Get the most recent Monday (could be today if today is Monday)
	currentWeekday := int(now.Weekday())
	daysToMonday := currentWeekday - 1 // Monday is 1
	if currentWeekday == 0 { // Sunday is 0
		daysToMonday = 6
	}
	
	// Last week's Monday is current Monday minus 7 days
	lastMonday := now.AddDate(0, 0, -daysToMonday-7).Truncate(24 * time.Hour)
	lastSunday := lastMonday.AddDate(0, 0, 6) // 6 days after Monday = Sunday

	startDate := lastMonday.Format("2006-01-02")
	endDate := lastSunday.Format("2006-01-02")

	fmt.Printf("Sending timecards for period: %s to %s\n", startDate, endDate)

	// Send the timecards
	admin.SendTimeCards(startDate, endDate, false)

	fmt.Println("Timecard scheduler completed successfully")
}