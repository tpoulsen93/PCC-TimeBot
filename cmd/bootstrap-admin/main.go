// Command bootstrap-admin creates the first admin employee so that someone can
// sign in and provision the rest of the team. It is safe to run repeatedly: it
// will create the employee if the email is new, or promote an existing employee
// to admin if the email already exists.
//
// Usage:
//
//	bootstrap-admin -first Jane -last Doe -email jane@example.com [-phone 5551234567]
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tpoulsen/pcc-timebot/shared/database"
)

func main() {
	_ = godotenv.Overload()

	first := flag.String("first", "", "First name (required)")
	last := flag.String("last", "", "Last name (required)")
	emailFlag := flag.String("email", "", "Email address, used for sign-in (required)")
	phone := flag.String("phone", "", "Phone number (optional)")
	flag.Parse()

	if strings.TrimSpace(*first) == "" || strings.TrimSpace(*last) == "" || strings.TrimSpace(*emailFlag) == "" {
		fmt.Println("Error: -first, -last, and -email are required")
		flag.Usage()
		os.Exit(1)
	}

	if err := database.Initialize(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	email := strings.ToLower(strings.TrimSpace(*emailFlag))

	existing, err := database.GetEmployeeByEmail(email)
	if err != nil {
		fmt.Printf("Failed to look up employee: %v\n", err)
		os.Exit(1)
	}

	if existing != nil {
		if err := database.SetEmployeeAdmin(existing.ID, true); err != nil {
			fmt.Printf("Failed to promote employee to admin: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Existing employee %s (id %d) promoted to admin.\n", email, existing.ID)
		return
	}

	created, err := database.CreateEmployee(*first, *last, email, *phone, nil, true)
	if err != nil {
		fmt.Printf("Failed to create admin employee: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Admin employee created: %s %s <%s> (id %d).\n",
		*first, *last, email, created.ID)
	fmt.Println("They can now sign in via the web app to request a magic link.")
}
