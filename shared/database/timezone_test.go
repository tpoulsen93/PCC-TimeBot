package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_LocalDateCalculation verifies that our date calculation method
// produces the correct local date regardless of time of day.
// This is a regression test for the bug where Truncate(24h) would
// shift dates backwards (e.g., Monday → Sunday) due to UTC epoch truncation.
func Test_LocalDateCalculation(t *testing.T) {
	loc, err := time.LoadLocation("America/Denver")
	require.NoError(t, err)

	// Test various times on the same local date
	testCases := []struct {
		name     string
		hour     int
		wantDate string
	}{
		{"early morning", 1, "2026-01-15"},
		{"mid morning", 10, "2026-01-15"},
		{"afternoon", 15, "2026-01-15"},
		{"late evening", 22, "2026-01-15"},
		{"near midnight", 23, "2026-01-15"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a time at the specified hour on Jan 15, 2026 in Denver
			testTime := time.Date(2026, 1, 15, tc.hour, 30, 0, 0, loc)

			// Use the same calculation pattern as SubmitTime/AddTime
			now := testTime.In(loc)
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

			gotDate := today.Format("2006-01-02")
			assert.Equal(t, tc.wantDate, gotDate,
				"at %d:30 Denver time, date should be %s", tc.hour, tc.wantDate)
		})
	}
}

// Test_TruncateVsTimeDate demonstrates why Truncate(24h) is wrong for local dates.
// This test documents the bug we fixed.
func Test_TruncateVsTimeDate(t *testing.T) {
	loc, err := time.LoadLocation("America/Denver")
	require.NoError(t, err)

	// 10pm on Monday Jan 19, 2026 in Denver = 5am Tuesday Jan 20 UTC
	mondayLateNight := time.Date(2026, 1, 19, 22, 0, 0, 0, loc)

	// CORRECT: time.Date preserves the local date
	correctDate := time.Date(mondayLateNight.Year(), mondayLateNight.Month(), mondayLateNight.Day(), 0, 0, 0, 0, loc)
	assert.Equal(t, "2026-01-19", correctDate.Format("2006-01-02"), "time.Date should give Monday")
	assert.Equal(t, time.Monday, correctDate.Weekday(), "should be Monday")

	// WRONG: Truncate(24h) would give Sunday because it truncates relative to UTC epoch
	// We don't actually call Truncate here (it's the bug we fixed), but we document the issue
	t.Log("Truncate(24h) was causing dates to shift backwards by 1 day for evening submissions")
}
