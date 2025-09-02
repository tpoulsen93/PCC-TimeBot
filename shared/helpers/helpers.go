// Package helpers provides utility functions used throughout the PCC-TimeBot application.
package helpers

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/tpoulsen/pcc-timebot/shared/constants"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	// GlobalWage represents the standard hourly wage rate
	// Deprecated: Use constants.DefaultWage instead
	GlobalWage = constants.DefaultWage
)

var (
	// Title provides proper case conversion for names and text
	Title = cases.Title(language.English)
)

// Round rounds a float64 to the specified number of decimal places.
// This is useful for monetary calculations and time precision.
func Round(num float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(num*shift) / shift
}

// GetUserInput prompts the user with the given message and returns their input.
// The input is trimmed of leading and trailing whitespace.
func GetUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
