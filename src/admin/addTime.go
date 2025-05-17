package admin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/src/database"
	"github.com/tpoulsen/pcc-timebot/src/helpers"
)

// AddTime runs the interactive command-line interface for adding time manually
func AddTime() error {
	// Initialize database connection
	if err := database.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Get user input
	firstName := helpers.GetUserInput("Enter employee first name:       ")
	lastName := helpers.GetUserInput("Enter employee last name:        ")
	date := helpers.GetUserInput("Enter date:     <YYYY-MM-DD>     ")
	hours := helpers.GetUserInput("Enter hours:                     ")

	// Print confirmation
	fmt.Printf("\nname:   %s %s\n", helpers.Title.String(firstName), helpers.Title.String(lastName))
	fmt.Printf("date:   %s\n", date)
	fmt.Printf("hours:  %s\n\n", hours)

	// Confirm submission
	confirm := helpers.GetUserInput("Submit? (y/n)   ")
	fmt.Println()

	if !strings.HasPrefix(strings.ToLower(confirm), "y") {
		fmt.Println("Cancelled. See you in the next life...")
		return nil
	}

	// Get employee ID
	id, err := database.GetEmployeeID(firstName, lastName)
	if err != nil {
		return fmt.Errorf("failed to get employee ID: %w", err)
	}
	if id == 0 {
		return fmt.Errorf("employee not found")
	}

	// Parse the date string
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	// Convert hours string to float64
	hoursFloat, err := strconv.ParseFloat(hours, 64)
	if err != nil {
		return fmt.Errorf("invalid hours format: %w", err)
	}

	// Add time to database
	result, err := database.AddTime(id, parsedDate, hoursFloat)
	if err != nil {
		return fmt.Errorf("failed to add time: %w", err)
	}

	fmt.Println(result)
	return nil
}
