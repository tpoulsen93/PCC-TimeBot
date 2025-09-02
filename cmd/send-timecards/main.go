package main

import (
	"fmt"

	"github.com/tpoulsen/pcc-timebot/internal/admin"
)

func main() {
	// Call SendTimeCards with empty args to trigger interactive mode
	admin.SendTimeCards("", "", false)
	fmt.Println("Time cards sent successfully!")
}
