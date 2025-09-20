package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tpoulsen/pcc-timebot/internal/admin"
)

func main() {
	// Define command line flags
	var (
		startDate     = flag.String("start", "", "Pay period start date (YYYY-MM-DD)")
		endDate       = flag.String("end", "", "Pay period end date (YYYY-MM-DD)")
		useLastPeriod = flag.Bool("lastperiod", false, "Use the last period's end date to calculate the next 7-day period")
		showHelp      = flag.Bool("help", false, "Show usage information")
	)

	// Customize usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Send time cards for a specified pay period.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                                    # Interactive mode - prompts for dates\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -start 2024-01-01 -end 2024-01-07  # Specify exact date range\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -lastperiod                        # Use next 7-day period after last run\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -help                              # Show this help message\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDate format: YYYY-MM-DD\n")
		fmt.Fprintf(os.Stderr, "\nRequired environment variables:\n")
		fmt.Fprintf(os.Stderr, "  SMTP_USERNAME - Email username for sending\n")
		fmt.Fprintf(os.Stderr, "  SMTP_PASSWORD - Email password for sending\n")
	}

	// Parse command line flags
	flag.Parse()

	// Show help if requested
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Validate flag combinations
	if *useLastPeriod && (*startDate != "" || *endDate != "") {
		fmt.Fprintf(os.Stderr, "Error: Cannot use -lastperiod with -start or -end flags\n")
		fmt.Fprintf(os.Stderr, "Use -help for usage information\n")
		os.Exit(1)
	}

	if (*startDate != "" && *endDate == "") || (*startDate == "" && *endDate != "") {
		fmt.Fprintf(os.Stderr, "Error: Both -start and -end must be provided together\n")
		fmt.Fprintf(os.Stderr, "Use -help for usage information\n")
		os.Exit(1)
	}

	// Call SendTimeCards with appropriate parameters
	admin.SendTimeCards(*startDate, *endDate, *useLastPeriod)
	fmt.Println("Time cards sent successfully!")
}
