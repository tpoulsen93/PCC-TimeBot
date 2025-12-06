package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/internal/admin"
	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
	"github.com/tpoulsen/pcc-timebot/shared/timecalc"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/client"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwiMLResponse struct {
	XMLName xml.Name `xml:"Response"`
	Message []string `xml:"Message"`
}

func main() {
	// If the heroku flag is set, run the server on Heroku
	heroku := flag.Bool("heroku", false, "Run on Heroku")

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

	if *heroku {
		fmt.Println("Server running on Heroku")
		fmt.Println("Starting server...")
		runOnHeroku()
		fmt.Println("Server terminated")
		return
	}

	fmt.Println("Running as admin")
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

func textUsage(c *gin.Context) {
	resp := TwiMLResponse{
		Message: []string{
			"Usage: Time <first name> <last name> <start time> <end time> <subtracted hours(lunch)> [<additional hours(drive time)>] [\"<job location>\"]\nExample:",
			"Time Taylor Poulsen 11:46am 5:04pm 1.25 3.6 \"Main Street Project\"",
		},
	}
	c.XML(http.StatusOK, resp)
}

func runOnHeroku() {
	// Check if running on Heroku
	if os.Getenv("DYNO") != "" {
		os.Exit(1)
	}

	// Initialize database connection
	if err := database.Initialize(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	r := gin.Default()

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"PCC": "Poulsen Concrete Contractors Inc.",
		})
	})

	// Submit hours endpoint
	r.GET("/submitHours", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		date := c.Query("date")
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format")
			return
		}

		hours, err := strconv.ParseFloat(c.Query("hours"), 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid hours")
			return
		}

		result, err := database.AddTime(id, parsedDate, hours, "")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Send confirmation
		if err := twilioConfirmSubmission(id, result, "TimeApp"); err != nil {
			fmt.Printf("Failed to send confirmation: %v\n", err)
		}

		c.String(http.StatusOK, result)
	})	// SMS webhook endpoint
	r.POST("/sms", func(c *gin.Context) {
		// Validate Twilio signature
		if !validateTwilioRequest(c) {
			fmt.Println("unexpected user encountered")
			c.String(http.StatusBadRequest, "Error in Twilio Signature")
			return
		}

		from := c.PostForm("From")
		body := c.PostForm("Body")

		// Process the message
		msg, err := processMessage(body, from)
		if err != nil {
			resp := TwiMLResponse{
				Message: []string{"Encountered an unexpected error. Check your format and try again."},
			}
			fmt.Printf("Encountered unexpected error in message: [%s]\n%v\n", body, err)
			c.XML(http.StatusOK, resp)
			return
		}

		if msg == "" {
			fmt.Printf("Ignored message from %s:\n[%s]\n", from, body)
			c.String(http.StatusOK, "")
			return
		}

		// Build response
		if strings.HasPrefix(msg, "Help") {
			textUsage(c)
		} else if strings.HasPrefix(msg, "Error") {
			resp := TwiMLResponse{Message: []string{msg}}
			c.XML(http.StatusOK, resp)
			textUsage(c)
		} else {
			resp := TwiMLResponse{Message: []string{msg}}
			c.XML(http.StatusOK, resp)
		}

		fmt.Printf("Processed message from %s:\n[%s]\nResponded:\n[%s]\n", from, body, msg)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server started on port", port)
}

func validateTwilioRequest(c *gin.Context) bool {
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioAuthToken == "" {
		return false
	}

	validator := client.NewRequestValidator(twilioAuthToken)
	url := c.Request.URL.String()
	params := make(map[string]string)
	for key, values := range c.Request.PostForm {
		params[key] = values[0]
	}
	signature := c.GetHeader("X-Twilio-Signature")

	return validator.Validate(url, params, signature)
}

func processMessage(message, from string) (string, error) {
	msg := strings.ToLower(message)

	// Handle time submission
	if strings.HasPrefix(msg, "time") || strings.HasPrefix(msg, "hours") {
		if strings.Contains(msg, "help") {
			return "Help", nil
		}
		return processTime(message, from)
	}

	// Message not meant for us
	return "", nil
}

func processTime(message, from string) (string, error) {
	// Extract quoted location if present (e.g., "Main Street Project")
	location := ""
	locationRegex := regexp.MustCompile(`"([^"]+)"`)
	if match := locationRegex.FindStringSubmatch(message); len(match) > 1 {
		location = match[1]
		// Remove the quoted location from the message for normal parsing
		message = locationRegex.ReplaceAllString(message, "")
	}

	parts := strings.Fields(message)
	if len(parts) < 6 {
		return "Error. Time formatted incorrectly. Too few parameters.", nil
	}
	if len(parts) > 7 {
		return "Error. Time formatted incorrectly. Too many parameters.", nil
	}

	// Get employee ID
	employeeID, err := database.GetEmployeeID(parts[1], parts[2])
	if err != nil {
		return "", fmt.Errorf("failed to get employee ID: %w", err)
	}
	if employeeID == 0 {
		return "Error. Employee not found.", nil
	}

	// Get time parameters
	start := parts[3]
	end := parts[4]
	less := parts[5]
	more := "0"
	if len(parts) == 7 {
		more = parts[6]
	}

	// Calculate hours
	hours, err := timecalc.CalculateTime(start, end, less, more)
	if err != nil {
		switch err {
		case timecalc.ErrHours:
			return "Error. Time formatted incorrectly. Hours spot is wrong.", nil
		case timecalc.ErrMeridiem:
			return "Error. Time formatted incorrectly. Meridiem is wrong. (am/pm)", nil
		case timecalc.ErrMinutes:
			return "Error. Time formatted incorrectly. Minutes spot is wrong.", nil
		case timecalc.ErrIllegalTime:
			return "Error. Time formatted incorrectly. End time is before start time...", nil
		case timecalc.ErrLunch:
			return "Error. Subtracted hours formatted incorrectly.", nil
		case timecalc.ErrExtra:
			return "Error. Additional hours formatted incorrectly.", nil
		default:
			return "Error. Time formatted incorrectly.", nil
		}
	}

	// Submit the time
	submission, err := database.SubmitTime(employeeID, hours, message, location)
	if err != nil {
		return "", fmt.Errorf("failed to submit time: %w", err)
	}

	// Build response message
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		loc = time.UTC
	}
	today := time.Now().In(loc).Truncate(24 * time.Hour)

	moreHours, _ := strconv.ParseFloat(more, 64)
	lessHours, _ := strconv.ParseFloat(less, 64)

	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s\n", today.Format("2006-01-02")))
	result.WriteString(fmt.Sprintf("%s %s\n", helpers.Title.String(parts[1]), helpers.Title.String(parts[2])))
	result.WriteString(fmt.Sprintf("Start: %s\n", start))
	result.WriteString(fmt.Sprintf("End: %s\n", end))
	if lessHours > 0 {
		result.WriteString(fmt.Sprintf("Lunch hours: %s\n", less))
	}
	if moreHours > 0 {
		result.WriteString(fmt.Sprintf("Extra hours: %s\n", more))
	}
	if location != "" {
		result.WriteString(fmt.Sprintf("Location: %s\n", location))
	}
	result.WriteString(submission)

	// Send confirmations
	if err := twilioConfirmSubmission(employeeID, result.String(), from); err != nil {
		return fmt.Sprintf("%s\n%s", result.String(), err.Error()), nil
	}

	return result.String(), nil
}

func twilioConfirmSubmission(employeeID int, msg, from string) error {
	adminID, err := database.GetEmployeeID("admin", "admin")
	if err != nil {
		return fmt.Errorf("failed to get admin ID: %w", err)
	}

	ownerID, err := database.GetEmployeeID("jr", "poulsen")
	if err != nil {
		return fmt.Errorf("failed to get owner ID: %w", err)
	}

	supervisorID, err := database.GetSupervisorID(employeeID)
	if err != nil {
		return fmt.Errorf("failed to get supervisor ID: %w", err)
	}

	isAdmin := employeeID < 3

	// Send confirmations
	if !isAdmin {
		if err := sendConfirmation(adminID, "Admin", msg); err != nil {
			return err
		}
		if err := sendConfirmation(ownerID, "Owner", msg); err != nil {
			return err
		}
		if err := sendConfirmation(supervisorID, "Supervisor", msg); err != nil {
			return err
		}
	}

	// Send to recipient if they aren't already texted in the response
	phone, err := database.GetEmployeePhone(employeeID)
	if err != nil {
		return fmt.Errorf("failed to get employee phone: %w", err)
	}

	if phone != strings.TrimPrefix(from, "+1") {
		if err := sendText(employeeID, msg); err != nil {
			return fmt.Errorf("failed to send text to employee: %w", err)
		}
	}

	return nil
}

func sendConfirmation(id int, userType string, msg string) error {
	if id <= 0 {
		return nil
	}

	phone, err := database.GetEmployeePhone(id)
	if err != nil {
		return fmt.Errorf("failed to get %s phone: %w", userType, err)
	}
	if phone == "" {
		return fmt.Errorf("error: %s phone not found", userType)
	}

	return sendText(id, msg)
}

func sendText(id int, msg string) error {
	phone, err := database.GetEmployeePhone(id)
	if err != nil || phone == "" {
		return nil
	}

	clientParams := twilio.ClientParams{
		AccountSid: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password:   os.Getenv("TWILIO_AUTH_TOKEN"),
	}

	twilioClient := twilio.NewRestClientWithParams(clientParams)
	twilioPhone := os.Getenv("TWILIO_PHONE")

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(fmt.Sprintf("+1%s", phone))
	params.SetFrom(fmt.Sprintf("+1%s", twilioPhone))
	params.SetBody(msg)

	_, err = twilioClient.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return err
	}

	return nil
}
