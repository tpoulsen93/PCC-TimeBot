package admin

import (
	"fmt"
	"time"

	"github.com/tpoulsen/pcc-timebot/internal/email"
	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/timecard"
)

// BuildTimeCards builds per-employee time cards for the inclusive date range.
// It returns the time cards keyed by employee ID and the computed payday. This
// is the HTTP-safe counterpart to the CLI SendTimeCards flow: it never calls
// os.Exit and never reads from stdin.
func BuildTimeCards(start, end time.Time) (map[int]*timecard.TimeCard, time.Time, error) {
	startDate := start.Format("2006-01-02")
	endDate := end.Format("2006-01-02")

	entries, err := database.GetTimeCards(start, end)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to get time cards: %w", err)
	}

	timeCards := make(map[int]*timecard.TimeCard)
	var payday time.Time

	// Stored dates are UTC; convert to Denver local dates to match pay-period days.
	loc, _ := time.LoadLocation("America/Denver")

	for _, entry := range entries {
		tc, exists := timeCards[entry.EmployeeID]
		if !exists {
			tc, err = timecard.NewTimeCard(entry.EmployeeID, startDate, endDate)
			if err != nil {
				return nil, time.Time{}, fmt.Errorf("failed to create time card for employee %d: %w", entry.EmployeeID, err)
			}
			timeCards[entry.EmployeeID] = tc
			if payday.IsZero() {
				payday = tc.PayDay
			}
		}

		localDate := entry.Date.In(loc).Format("2006-01-02")
		if err := tc.AddHours(localDate, entry.Time, entry.Location); err != nil {
			return nil, time.Time{}, fmt.Errorf("failed to add hours for employee %d: %w", entry.EmployeeID, err)
		}
	}

	return timeCards, payday, nil
}

// SendTimeCardsForRange builds and emails individual time cards for the date
// range and emails a payroll summary to the admin. It returns the number of
// time cards sent. Intended to be called from an authenticated admin HTTP
// endpoint, so it returns errors instead of exiting.
func SendTimeCardsForRange(start, end time.Time) (int, error) {
	timeCards, payday, err := BuildTimeCards(start, end)
	if err != nil {
		return 0, err
	}
	if len(timeCards) == 0 {
		return 0, nil
	}

	smtpConfig := email.NewSMTPConfig()
	if smtpConfig.Username == "" || smtpConfig.Password == "" {
		return 0, fmt.Errorf("SMTP credentials not configured")
	}

	sent := 0
	for _, tc := range timeCards {
		if tc.Email == "" {
			continue
		}
		if err := email.SendTimeCard(smtpConfig,
			smtpConfig.Username,
			tc.Email,
			tc.Name,
			tc.ToHTML(),
			tc.PayDay); err != nil {
			return sent, fmt.Errorf("failed to send time card to %s: %w", tc.Name, err)
		}
		sent++
	}

	// Send the payroll summary to the configured admin address (SMTP user).
	if err := email.SendPayrollSummary(smtpConfig,
		smtpConfig.Username,
		smtpConfig.Username,
		timeCards,
		start,
		end,
		payday); err != nil {
		return sent, fmt.Errorf("failed to send payroll summary: %w", err)
	}

	return sent, nil
}
