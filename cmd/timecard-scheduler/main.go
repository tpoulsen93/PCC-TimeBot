package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tpoulsen/pcc-timebot/internal/admin"
	"github.com/tpoulsen/pcc-timebot/shared/database"
)

// calculateLastWeekDates calculates the Sunday-Saturday date range for the last completed pay period
// based on the provided current time in the specified timezone.
func calculateLastWeekDates(now time.Time, timezone string) (startDate, endDate string, err error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", "", fmt.Errorf("failed to load timezone: %w", err)
	}

	localNow := now.In(loc)

	// IMPORTANT: Do not use Truncate(24h) to get "midnight".
	// Truncate rounds relative to the Unix epoch (effectively UTC day boundaries),
	// which can shift the local date backwards and cause off-by-one errors
	// (e.g., returning Saturday instead of Sunday and making payday Thursday instead of Friday).
	localDate := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)

	// We want the last *completed* Sunday-Saturday pay period.
	// Find the start of the current pay period (most recent Sunday, possibly today).
	currentWeekday := int(localDate.Weekday()) // Sunday=0 ... Saturday=6
	currentPeriodStart := localDate.AddDate(0, 0, -currentWeekday)

	// Last period is the 7 days immediately before the current period.
	lastPeriodStart := currentPeriodStart.AddDate(0, 0, -7) // Sunday
	lastPeriodEnd := lastPeriodStart.AddDate(0, 0, 6)       // Saturday

	startDate = lastPeriodStart.Format("2006-01-02")
	endDate = lastPeriodEnd.Format("2006-01-02")

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
	fmt.Printf("Current time in Denver: %s (%s)\n", now.Format("2006-01-02 15:04:05 MST"), now.Weekday())

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

	// Calculate last completed pay period (Sunday-Saturday)
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
