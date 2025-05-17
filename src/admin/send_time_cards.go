package admin

import (
	"fmt"
	"os"
	"time"

	"github.com/tpoulsen/pcc-timebot/src/database"
	"github.com/tpoulsen/pcc-timebot/src/email"
	"github.com/tpoulsen/pcc-timebot/src/helpers"
	"github.com/tpoulsen/pcc-timebot/src/timecard"
)

func SendTimeCards() {
	// Initialize database connection
	if err := database.Initialize(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Date format:  YYYY-MM-DD")
	startDate := helpers.GetUserInput("Enter pay period start date:  ")
	endDate := helpers.GetUserInput("Enter pay period end date:    ")

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

		if err := tc.AddHours(entry.Date.Format("2006-01-02"), entry.Time); err != nil {
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
			[]byte(tc.String()),
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
}
