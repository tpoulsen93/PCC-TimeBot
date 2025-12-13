package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tpoulsen/pcc-timebot/internal/admin"
	"github.com/tpoulsen/pcc-timebot/shared/database"
)

// calculateLastWeekDates calculates the Monday-Sunday date range for the previous week
// based on the provided current time in the specified timezone.
func calculateLastWeekDates(now time.Time, timezone string) (startDate, endDate string, err error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", "", fmt.Errorf("failed to load timezone: %w", err)
	}

	localNow := now.In(loc)
	
	// Get the most recent Monday (could be today if today is Monday)
	currentWeekday := int(localNow.Weekday())
	daysToMonday := currentWeekday - 1 // Monday is 1
	if currentWeekday == 0 { // Sunday is 0
		daysToMonday = 6
	}
	
	// Last week's Monday is current Monday minus 7 days
	lastMonday := localNow.AddDate(0, 0, -daysToMonday-7).Truncate(24 * time.Hour)
	lastSunday := lastMonday.AddDate(0, 0, 6) // 6 days after Monday = Sunday

	startDate = lastMonday.Format("2006-01-02")
	endDate = lastSunday.Format("2006-01-02")
	
	return startDate, endDate, nil
}

func main() {
	fmt.Println("Starting timecard scheduler...")

	// Check if today is Monday
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Printf("Failed to load timezone: %v\n", err)
		os.Exit(1)
	}
	
	now := time.Now().In(loc)
	if now.Weekday() != time.Monday {
		fmt.Printf("Today is %s, not Monday. Skipping timecard send.\n", now.Weekday())
		return
	}

	fmt.Println("Today is Monday. Proceeding with timecard send...")

	// Initialize database connection
	if err := database.Initialize(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Calculate last week's Monday-Sunday
	startDate, endDate, err := calculateLastWeekDates(now, "America/Denver")
	if err != nil {
		fmt.Printf("Failed to calculate dates: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sending timecards for period: %s to %s\n", startDate, endDate)

	// Send the timecards
	admin.SendTimeCards(startDate, endDate, false)

	fmt.Println("Timecard scheduler completed successfully")
}