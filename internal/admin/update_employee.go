package admin

import (
	"fmt"
	"strings"

	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
)

func UpdateEmployee() {
	if err := database.Initialize(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	first := helpers.GetUserInput("Employee first name:             ")
	last := helpers.GetUserInput("Employee last name:              ")
	target := helpers.GetUserInput("Field to update (first_name | last_name | email | phone | supervisor_id): ")

	var value string
	if strings.ToLower(strings.TrimSpace(target)) == "phone" {
		for {
			raw := helpers.GetUserInput("New phone (10 digits, e.g. 3035551234): ")
			normalized, err := normalizePhone(raw)
			if err != nil {
				fmt.Printf("Invalid phone: %v. Please try again.\n", err)
				continue
			}
			value = normalized
			break
		}
	} else {
		value = helpers.GetUserInput("New value:                       ")
	}

	fmt.Printf("\nname:   %s %s\n", helpers.Title.String(first), helpers.Title.String(last))
	fmt.Printf("field:  %s\n", target)
	fmt.Printf("value:  %s\n\n", value)

	confirm := helpers.GetUserInput("Submit? (y/n): ")
	fmt.Println()

	if strings.HasPrefix(strings.ToLower(confirm), "y") {
		err := database.UpdateEmployee(first, last, target, value)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("%s %s's %s updated to %s\n",
			helpers.Title.String(first), helpers.Title.String(last), target, value)
	} else {
		fmt.Println("Cancelled.")
	}
}
