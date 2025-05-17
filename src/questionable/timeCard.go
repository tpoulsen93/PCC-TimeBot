package timecalc

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/src/helpers"
)

// TimeCard represents an employee's time card for a pay period
type TimeCard struct {
	ID         int
	Name       string
	Email      string
	Phone      string
	Days       map[string]float64
	TotalHours float64
	PayDay     time.Time
}

// NewTimeCard creates a new TimeCard instance for the given employee and date range
func NewTimeCard(id int, startDate, endDate string) (*TimeCard, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Initialize timecard with empty days map
	tc := &TimeCard{
		ID:         id,
		Days:       make(map[string]float64),
		TotalHours: 0,
		PayDay:     end.AddDate(0, 0, 12),
	}

	// Fill days map with zero hours for each day in range
	current := start
	for !current.After(end) {
		tc.Days[current.Format("2006-01-02")] = 0
		current = current.AddDate(0, 0, 1)
	}

	return tc, nil
}

// AddHours adds hours worked for a specific date
func (tc *TimeCard) AddHours(date string, hours float64) {
	tc.Days[date] = helpers.Round(hours, 2)
	tc.TotalHours += hours
}

// buildDayLine creates a formatted string for a single day's entry
func (tc *TimeCard) buildDayLine(date string) string {
	t, _ := time.Parse("2006-01-02", date)
	dayOfWeek := t.Format("Mon")
	hours := tc.Days[date]

	return fmt.Sprintf("%10s | %4s| %5.2f\n", date, dayOfWeek, hours)
}

// String implements the Stringer interface for TimeCard
func (tc *TimeCard) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n\n\n", tc.Name))
	sb.WriteString(fmt.Sprintf("%11s|%5s|%6s\n", "Date", "Day", "Hours"))
	sb.WriteString(fmt.Sprintf("%11s+%5s+%6s\n",
		strings.Repeat("-", 11),
		strings.Repeat("-", 5),
		strings.Repeat("-", 6)))

	// Sort and write days
	dates := make([]string, 0, len(tc.Days))
	for date := range tc.Days {
		dates = append(dates, date)
	}
	// Sort dates chronologically
	sort.Strings(dates)

	for _, date := range dates {
		sb.WriteString(tc.buildDayLine(date))
	}

	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("Total hours:  %.2f\n", tc.TotalHours))
	sb.WriteString(fmt.Sprintf("Payday:  %s\n", tc.PayDay.Format("2006-01-02")))

	return sb.String()
}
