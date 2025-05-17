package helpers

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const GlobalWage = 30

var Title = cases.Title(language.English)

// Round helps round float64 to specified decimal places
func Round(num float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(num*shift) / shift
}

func GetUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
