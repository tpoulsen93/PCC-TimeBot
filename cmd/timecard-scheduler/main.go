package timecardscheduler
package main

import (
	"fmt"














































}	fmt.Println("Timecard scheduler completed successfully")	admin.SendTimeCards(startDate, endDate, false)	// Send the timecards	fmt.Printf("Sending timecards for period: %s to %s\n", startDate, endDate)	endDate := lastSunday.Format("2006-01-02")	startDate := lastMonday.Format("2006-01-02")	lastSunday := lastMonday.AddDate(0, 0, 6) // 6 days after Monday = Sunday	lastMonday := now.AddDate(0, 0, -daysToMonday-7).Truncate(24 * time.Hour)	// Last week's Monday is current Monday minus 7 days		}		daysToMonday = 6	if currentWeekday == 0 { // Sunday is 0	daysToMonday := currentWeekday - 1 // Monday is 1	currentWeekday := int(now.Weekday())	// Get the most recent Monday (could be today if today is Monday)		now := time.Now().In(loc)	}		os.Exit(1)		fmt.Printf("Failed to load timezone: %v\n", err)	if err != nil {	loc, err := time.LoadLocation("America/Denver")	// Calculate last week's Monday-Sunday	}		os.Exit(1)		fmt.Printf("Failed to initialize database: %v\n", err)	if err := database.Initialize(); err != nil {	// Initialize database connection	fmt.Println("Starting timecard scheduler...")func main() {)	"github.com/tpoulsen/pcc-timebot/shared/database"	"github.com/tpoulsen/pcc-timebot/internal/admin"	"time"	"os"