// Package timecard provides functionality for creating and managing employee time cards.
package timecard

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
)

// DayEntry represents a single day's work entry with hours and optional location
type DayEntry struct {
	Hours    float64
	Location string
}

// TimeCard represents a time card for an employee during a specific pay period.
type TimeCard struct {
	ID         int                 // Employee ID
	Name       string              // Full employee name
	Email      string              // Employee email address
	Phone      string              // Employee phone number
	Days       map[string]DayEntry // Map of date (YYYY-MM-DD) to day entry
	TotalHours float64             // Total hours worked during the pay period
	PayDay     time.Time           // Expected payday
}

// computePayday returns the payday for a given pay period end date.
// Historically, PCC payroll is paid on the *second Friday* after the pay period ends.
// This keeps payday on a Friday regardless of whether the pay period ends on Saturday or Sunday.
func computePayday(periodEnd time.Time) time.Time {
	// Normalize to a date-only value in the same location.
	endDate := time.Date(periodEnd.Year(), periodEnd.Month(), periodEnd.Day(), 0, 0, 0, 0, periodEnd.Location())

	// Days until the next Friday.
	// If the period end date is already a Friday, the "next" Friday is the following week
	// because payday is defined as the second Friday *after* the pay period ends.
	// Go: Sunday=0 ... Saturday=6
	daysUntilFriday := (int(time.Friday) - int(endDate.Weekday()) + 7) % 7
	if daysUntilFriday == 0 {
		daysUntilFriday = 7
	}
	firstFriday := endDate.AddDate(0, 0, daysUntilFriday)

	// Payroll is on the second Friday after period end.
	return firstFriday.AddDate(0, 0, 7)
}

// NewTimeCard creates a new time card for an employee covering the specified date range.
// The date range is inclusive of both start and end dates.
//
// Parameters:
//   - id: Employee ID from the database
//   - startDate: Start date in "YYYY-MM-DD" format
//   - endDate: End date in "YYYY-MM-DD" format
//
// Returns a populated TimeCard or an error if the employee doesn't exist or dates are invalid.
func NewTimeCard(id int, startDate, endDate string) (*TimeCard, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start date: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end date: %w", err)
	}

	if end.Before(start) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	// Create the timecard with initial values
	tc := &TimeCard{
		ID:         id,
		Days:       make(map[string]DayEntry),
		TotalHours: 0,
		PayDay:     computePayday(end),
	}

	// Initialize days map with zero hours for each day
	current := start
	for !current.After(end) {
		tc.Days[current.Format("2006-01-02")] = DayEntry{Hours: 0, Location: ""}
		current = current.AddDate(0, 0, 1)
	}

	// Get employee info from database
	employee, err := database.GetEmployee(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee info: %w", err)
	}

	tc.Name = helpers.Title.String(employee.FirstName) + " " + helpers.Title.String(employee.LastName)
	tc.Email = employee.Email
	tc.Phone = employee.Phone

	return tc, nil
}

// AddHours adds hours to a specific date in the time card.
// The date must be within the pay period range and hours must be non-negative.
//
// Parameters:
//   - date: Date in "YYYY-MM-DD" format
//   - hours: Number of hours worked (must be >= 0)
//   - location: Optional job location/name
//
// Returns an error if the date is not in the pay period or hours are invalid.
func (tc *TimeCard) AddHours(date string, hours float64, location string) error {
	if hours < 0 {
		return fmt.Errorf("hours cannot be negative: %.2f", hours)
	}

	if _, exists := tc.Days[date]; !exists {
		return fmt.Errorf("date %s not in pay period", date)
	}

	// Round hours to 2 decimal places for consistency
	hours = helpers.Round(hours, 2)

	entry := tc.Days[date]
	entry.Hours += hours
	entry.Location = location
	tc.Days[date] = entry
	tc.TotalHours = helpers.Round(tc.TotalHours+hours, 2)
	return nil
}

// String returns a string representation of the time card
func (tc *TimeCard) String() string {
	var sb strings.Builder

	// Header
	sb.WriteString(tc.Name + "\n\n\n")
	sb.WriteString(fmt.Sprintf("%11s|%5s|%6s|%s\n", "Date", "Day", "Hours", "Job Name"))
	sb.WriteString(fmt.Sprintf("%s+%s+%s+%s\n", strings.Repeat("-", 11), strings.Repeat("-", 5), strings.Repeat("-", 6), strings.Repeat("-", 20)))

	// Days - sort them by date
	dates := make([]string, 0, len(tc.Days))
	for date := range tc.Days {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, date := range dates {
		t, _ := time.Parse("2006-01-02", date)
		dayName := t.Format("Mon")[:3]
		entry := tc.Days[date]
		location := entry.Location
		if location == "" {
			location = "-"
		}
		sb.WriteString(fmt.Sprintf("%10s | %4s| %5.2f| %s\n", date, dayName, entry.Hours, location))
	}

	// Footer
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("Total hours:  %.2f\n", tc.TotalHours))
	sb.WriteString(fmt.Sprintf("Payday:  %s\n", tc.PayDay.Format("2006-01-02")))

	return sb.String()
}

// ToHTML returns an HTML representation of the time card with proper table formatting
func (tc *TimeCard) ToHTML() string {
	var sb strings.Builder

	// HTML email header with inline CSS
	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
<style>
body { font-family: Arial, sans-serif; margin: 20px; }
table { border-collapse: collapse; width: 100%; max-width: 700px; margin: 20px 0; }
th, td { padding: 10px; text-align: left; border: 1px solid #ddd; }
th { background-color: #4CAF50; color: white; font-weight: bold; }
tr:nth-child(even) { background-color: #f2f2f2; }
.total-row { font-weight: bold; background-color: #e8f5e9 !important; }
.hours-col { text-align: right; }
.header { margin-bottom: 10px; }
.footer { margin-top: 20px; font-size: 14px; }
</style>
</head>
<body>
`)

	sb.WriteString(fmt.Sprintf("<div class='header'><h2>%s</h2></div>\n", tc.Name))
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>Date</th><th>Day</th><th class='hours-col'>Hours</th><th>Job Name</th></tr>\n")

	// Days - sort them by date
	dates := make([]string, 0, len(tc.Days))
	for date := range tc.Days {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, date := range dates {
		t, _ := time.Parse("2006-01-02", date)
		dayName := t.Format("Mon")
		entry := tc.Days[date]
		location := entry.Location
		if location == "" {
			location = "-"
		}
		sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td class='hours-col'>%.2f</td><td>%s</td></tr>\n",
			date, dayName, entry.Hours, location))
	}

	// Total row
	sb.WriteString(fmt.Sprintf("<tr class='total-row'><td colspan='2'>Total</td><td class='hours-col'>%.2f</td><td></td></tr>\n", tc.TotalHours))
	sb.WriteString("</table>\n")

	// Payday
	sb.WriteString(fmt.Sprintf("<div class='footer'><p><strong>Payday:</strong> %s</p></div>\n", tc.PayDay.Format("2006-01-02")))

	sb.WriteString("</body>\n</html>")
	return sb.String()
}
