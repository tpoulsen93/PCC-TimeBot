package admin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/internal/email"
	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
	"github.com/tpoulsen/pcc-timebot/shared/timecard"
)

// getLastEndDateFilePath returns the path to the .last_end_date file
// Uses the user's home directory to ensure consistent location regardless of where binary is run
func getLastEndDateFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".pcc-timebot-last-end-date"), nil
}

func SendTimeCards(startDateArg, endDateArg string, useLastPeriod bool) {
	// Initialize database connection
	if err := database.Initialize(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Date format:  YYYY-MM-DD")
	var startDate, endDate string

	if useLastPeriod {
		startDate, endDate = getDatesFromLastPeriod()
	} else if startDateArg != "" && endDateArg != "" {
		startDate, endDate = getDatesFromArgs(startDateArg, endDateArg)
	} else {
		startDate, endDate = getDatesFromUserInput()
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Printf("Invalid start date format: %v\n", err)
		os.Exit(1)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Printf("Invalid end date format: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nGetting all hours submitted between %s and %s...\n", startDate, endDate)

	// Get time cards from database
	entries, err := database.GetTimeCards(start, end)
	if err != nil {
		fmt.Printf("Failed to get time cards: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("No hours found for indicated dates...")
		os.Exit(0)
	}

	// Build time cards
	fmt.Println("Building time cards...")
	timeCards := make(map[int]*timecard.TimeCard)
	var payday time.Time

	for _, entry := range entries {
		tc, exists := timeCards[entry.EmployeeID]
		if !exists {
			tc, err = timecard.NewTimeCard(entry.EmployeeID, startDate, endDate)
			if err != nil {
				fmt.Printf("Failed to create time card for employee %d: %v\n", entry.EmployeeID, err)
				continue
			}
			timeCards[entry.EmployeeID] = tc
			if payday.IsZero() {
				payday = tc.PayDay
			}
		}

		if err := tc.AddHours(entry.Date.Format("2006-01-02"), entry.Time, entry.Location); err != nil {
			fmt.Printf("Failed to add hours for employee %d: %v\n", entry.EmployeeID, err)
		}
	}

	// Setup SMTP configuration
	smtpConfig := email.NewSMTPConfig()
	if smtpConfig.Username == "" {
		fmt.Println("SMTP_USERNAME environment variable not set")
		os.Exit(1)
	}
	if smtpConfig.Password == "" {
		fmt.Println("SMTP_PASSWORD environment variable not set")
		os.Exit(1)
	}

	fmt.Println("Connecting to SMTP server...")

	// Send individual time cards
	for _, tc := range timeCards {
		fmt.Printf("Sending time card to %s...\n", tc.Name)

		if err := email.SendTimeCard(smtpConfig,
			smtpConfig.Username,
			tc.Email,
			tc.Name,
			tc.ToHTML(),
			tc.PayDay); err != nil {
			fmt.Printf("Failed to send time card to %s: %v\n", tc.Name, err)
		}
	}

	// Send summary to admin
	fmt.Println("Sending payroll totals to admin...")
	adminID, err := database.GetEmployeeID("taylor", "poulsen")
	if err != nil {
		fmt.Printf("Failed to get admin ID: %v\n", err)
		os.Exit(1)
	}
	if adminID == 0 {
		fmt.Println("Admin not found in database")
		os.Exit(1)
	}

	adminEmail, err := database.GetEmployeeEmail(adminID)
	if err != nil {
		fmt.Printf("Failed to get admin email: %v\n", err)
		os.Exit(1)
	}

	if err := email.SendPayrollSummary(smtpConfig,
		smtpConfig.Username,
		adminEmail,
		timeCards,
		start,
		end,
		payday); err != nil {
		fmt.Printf("Failed to send payroll summary: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Mission accomplished")

	// Save the end date for future use
	lastEndDateFile, err := getLastEndDateFilePath()
	if err != nil {
		fmt.Printf("Warning: Failed to get config file path: %v\n", err)
	} else {
		if err := os.WriteFile(lastEndDateFile, []byte(endDate), 0644); err != nil {
			fmt.Printf("Warning: Failed to save last end date: %v\n", err)
		} else {
			fmt.Printf("Saved last end date to %s\n", lastEndDateFile)
		}
	}
}

// getDatesFromLastPeriod reads the last end date from file and calculates the next 7-day period
func getDatesFromLastPeriod() (string, string) {
	// Get the path to the last end date file
	lastEndDateFile, err := getLastEndDateFilePath()
	if err != nil {
		fmt.Printf("Failed to get config file path: %v\n", err)
		fmt.Println("Please provide dates manually or run without -lastperiod first.")
		os.Exit(1)
	}

	// Read the last end date from file
	lastEndData, err := os.ReadFile(lastEndDateFile)
	if err != nil {
		fmt.Printf("Failed to read last end date from %s: %v\n", lastEndDateFile, err)
		fmt.Println("Please provide dates manually or run without -lastperiod first.")
		os.Exit(1)
	}
	lastEndStr := strings.TrimSpace(string(lastEndData))
	lastEnd, err := time.Parse("2006-01-02", lastEndStr)
	if err != nil {
		fmt.Printf("Invalid last end date format: %v\n", err)
		os.Exit(1)
	}

	// Calculate next 7-day period
	start := lastEnd.AddDate(0, 0, 1) // Next day after last end
	end := start.AddDate(0, 0, 6)     // 6 days later for 7-day period

	startDate := start.Format("2006-01-02")
	endDate := end.Format("2006-01-02")

	fmt.Printf("Using last period's end date %s\n", lastEndStr)
	fmt.Printf("Calculated next period: %s to %s\n", startDate, endDate)

	return startDate, endDate
}

// getDatesFromArgs returns dates from command line arguments, prompting for missing ones
func getDatesFromArgs(startDateArg, endDateArg string) (string, string) {
	var startDate, endDate string

	if startDateArg != "" {
		startDate = startDateArg
	} else {
		startDate = helpers.GetUserInput("Enter pay period start date:  ")
	}
	if endDateArg != "" {
		endDate = endDateArg
	} else {
		endDate = helpers.GetUserInput("Enter pay period end date:    ")
	}

	return startDate, endDate
}

// getDatesFromUserInput prompts the user for both start and end dates
func getDatesFromUserInput() (string, string) {
	startDate := helpers.GetUserInput("Enter pay period start date:  ")
	endDate := helpers.GetUserInput("Enter pay period end date:    ")

	return startDate, endDate
}
