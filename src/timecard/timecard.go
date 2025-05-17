package timecard

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/src/database"
	"github.com/tpoulsen/pcc-timebot/src/helpers"
)

// TimeCard represents a time card for an employee
type TimeCard struct {
	ID         int
	Name       string
	Email      string
	Phone      string
	Days       map[string]float64
	TotalHours float64
	PayDay     time.Time
}

// NewTimeCard creates a new time card for an employee
func NewTimeCard(id int, startDate, endDate string) (*TimeCard, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start date: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end date: %w", err)
	}

	// Create the timecard with initial values
	tc := &TimeCard{
		ID:         id,
		Days:       make(map[string]float64),
		TotalHours: 0,
		PayDay:     end.AddDate(0, 0, 12), // 12 days after end date
	}

	// Initialize days map with zero hours for each day
	current := start
	for !current.After(end) {
		tc.Days[current.Format("2006-01-02")] = 0
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

// AddHours adds hours to a specific date in the time card
func (tc *TimeCard) AddHours(date string, hours float64) error {
	if _, exists := tc.Days[date]; !exists {
		return fmt.Errorf("date %s not in pay period", date)
	}

	tc.Days[date] = hours
	tc.TotalHours += hours
	return nil
}

// String returns a string representation of the time card
func (tc *TimeCard) String() string {
	var sb strings.Builder

	// Header
	sb.WriteString(tc.Name + "\n\n\n")
	sb.WriteString(fmt.Sprintf("%11s|%5s|%6s\n", "Date", "Day", "Hours"))
	sb.WriteString(fmt.Sprintf("%s+%s+%s\n", strings.Repeat("-", 11), strings.Repeat("-", 5), strings.Repeat("-", 6)))

	// Days - sort them by date
	dates := make([]string, 0, len(tc.Days))
	for date := range tc.Days {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, date := range dates {
		t, _ := time.Parse("2006-01-02", date)
		dayName := t.Format("Mon")[:3]
		sb.WriteString(fmt.Sprintf("%10s | %4s| %5.2f\n", date, dayName, tc.Days[date]))
	}

	// Footer
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("Total hours:  %.2f\n", tc.TotalHours))
	sb.WriteString(fmt.Sprintf("Payday:  %s\n", tc.PayDay.Format("2006-01-02")))

	return sb.String()
}
