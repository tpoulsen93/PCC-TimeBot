package admin

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
)

func UpdateEmployee() {
	// Get database connection string from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL environment variable not set")
		os.Exit(1)
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Get user input
	first := helpers.GetUserInput("Enter employee first name:       ")
	last := helpers.GetUserInput("Enter employee last name:        ")
	target := helpers.GetUserInput("Enter target: <email | phone | supervisor_id> ")
	value := helpers.GetUserInput("Enter new value:                 ")

	// Print confirmation
	fmt.Printf("\nname:   %s %s\n", helpers.Title.String(first), helpers.Title.String(last))
	fmt.Printf("target: %s\n", target)
	fmt.Printf("value:  %s\n\n", value)

	// Confirm submission
	confirm := helpers.GetUserInput("Submit? (y/n)   ")
	fmt.Println()

	if strings.HasPrefix(strings.ToLower(confirm), "y") {
		err := database.UpdateEmployee(first, last, target, value)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("%s %s's %s was changed to %s\n",
			helpers.Title.String(first), helpers.Title.String(last), target, value)
	} else {
		fmt.Println("Cancelled. See you in the next life...")
	}
}
