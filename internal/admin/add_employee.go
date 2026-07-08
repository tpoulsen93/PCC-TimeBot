package admin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
)

var nonDigit = regexp.MustCompile(`\D`)

// normalizePhone strips non-digit characters and prepends +1.
// Returns an error if the result is not exactly 10 digits.
func normalizePhone(input string) (string, error) {
	digits := nonDigit.ReplaceAllString(input, "")
	if len(digits) != 10 {
		return "", fmt.Errorf("phone number must be 10 digits (got %d)", len(digits))
	}
	return "+1" + digits, nil
}

// AddEmployee runs the interactive CLI for adding a new employee.
func AddEmployee() error {
	if err := database.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	firstName := helpers.GetUserInput("First name:                      ")
	lastName := helpers.GetUserInput("Last name:                       ")

	var phone string
	for {
		raw := helpers.GetUserInput("Phone (10 digits, e.g. 3035551234):")
		normalized, err := normalizePhone(raw)
		if err != nil {
			fmt.Printf("Invalid phone: %v. Please try again.\n", err)
			continue
		}
		phone = normalized
		break
	}

	email := helpers.GetUserInput("Email:                           ")
	superFirst := helpers.GetUserInput("Supervisor first name (or blank):")
	superLast := helpers.GetUserInput("Supervisor last name  (or blank):")
	adminInput := helpers.GetUserInput("Admin? (y/n):                    ")

	isAdmin := strings.HasPrefix(strings.ToLower(adminInput), "y")

	fmt.Printf("\nname:       %s %s\n", helpers.Title.String(firstName), helpers.Title.String(lastName))
	fmt.Printf("phone:      %s\n", phone)
	fmt.Printf("email:      %s\n", email)
	if superFirst != "" {
		fmt.Printf("supervisor: %s %s\n", helpers.Title.String(superFirst), helpers.Title.String(superLast))
	} else {
		fmt.Printf("supervisor: none\n")
	}
	fmt.Printf("admin:      %v\n\n", isAdmin)

	confirm := helpers.GetUserInput("Add employee? (y/n): ")
	fmt.Println()

	if !strings.HasPrefix(strings.ToLower(confirm), "y") {
		fmt.Println("Cancelled.")
		return nil
	}

	result, err := database.AddEmployee(firstName, lastName, email, phone, superFirst, superLast, isAdmin)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
