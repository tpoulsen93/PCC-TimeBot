package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/tpoulsen/pcc-timebot/internal/admin"
)

func main() {
	_ = godotenv.Overload()

	// if the addTime flag is set, run the addTime function
	addTime := flag.Bool("addTime", false, "Run AddTime")
	// if the updateEmployee flag is set, run the updateEmployee function
	updateEmployee := flag.Bool("updateEmployee", false, "Run UpdateEmployee")
	// if the sendTimeCards flag is set, run the sendTimeCards function
	sendTimeCards := flag.Bool("sendTimeCards", false, "Run SendTimeCards")
	// optional start date for sendTimeCards
	startDate := flag.String("startDate", "", "Pay period start date (YYYY-MM-DD)")
	// optional end date for sendTimeCards
	endDate := flag.String("endDate", "", "Pay period end date (YYYY-MM-DD)")
	// optional flag to use last period's end date and calculate next 7-day period
	useLastPeriod := flag.Bool("useLastPeriod", false, "Use last period's end date to calculate next 7-day period")

	flag.Parse()

	switch {
	case *addTime:
		fmt.Println("Running AddTime...")
		admin.AddTime()
	case *updateEmployee:
		fmt.Println("Running UpdateEmployee...")
		admin.UpdateEmployee()
	case *sendTimeCards:
		fmt.Println("Running SendTimeCards...")
		admin.SendTimeCards(*startDate, *endDate, *useLastPeriod)
	default:
		fmt.Println("No valid command provided. Use -addTime, -updateEmployee, or -sendTimeCards.")
		fmt.Println("For -sendTimeCards, you can optionally provide -startDate and -endDate (YYYY-MM-DD), or use -useLastPeriod to calculate the next 7-day period from the last used end date.")
	}
}
