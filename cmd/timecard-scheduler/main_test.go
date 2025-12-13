package main

import (
	"testing"
	"time"
)

func TestCalculateLastWeekDates(t *testing.T) {
	tests := []struct {
		name        string
		currentDate string // YYYY-MM-DD format
		timezone    string
		wantStart   string
		wantEnd     string
		wantErr     bool
	}{
		{
			name:        "Monday - should return last week Monday to Sunday",
			currentDate: "2025-12-15", // Monday
			timezone:    "America/Denver",
			wantStart:   "2025-12-01",
			wantEnd:     "2025-12-07",
			wantErr:     false,
		},
		{
			name:        "Tuesday - should return last week Monday to Sunday",
			currentDate: "2025-12-16", // Tuesday
			timezone:    "America/Denver",
			wantStart:   "2025-12-08",
			wantEnd:     "2025-12-14",
			wantErr:     false,
		},
		{
			name:        "Wednesday - should return last week Monday to Sunday",
			currentDate: "2025-12-17", // Wednesday
			timezone:    "America/Denver",
			wantStart:   "2025-12-08",
			wantEnd:     "2025-12-14",
			wantErr:     false,
		},
		{
			name:        "Thursday - should return last week Monday to Sunday",
			currentDate: "2025-12-18", // Thursday
			timezone:    "America/Denver",
			wantStart:   "2025-12-08",
			wantEnd:     "2025-12-14",
			wantErr:     false,
		},
		{
			name:        "Friday - should return last week Monday to Sunday",
			currentDate: "2025-12-19", // Friday
			timezone:    "America/Denver",
			wantStart:   "2025-12-08",
			wantEnd:     "2025-12-14",
			wantErr:     false,
		},
		{
			name:        "Saturday - should return last week Monday to Sunday",
			currentDate: "2025-12-20", // Saturday
			timezone:    "America/Denver",
			wantStart:   "2025-12-08",
			wantEnd:     "2025-12-14",
			wantErr:     false,
		},
		{
			name:        "Sunday - should return last week Monday to Sunday",
			currentDate: "2025-12-21", // Sunday
			timezone:    "America/Denver",
			wantStart:   "2025-12-08",
			wantEnd:     "2025-12-14",
			wantErr:     false,
		},
		{
			name:        "First Monday of year",
			currentDate: "2026-01-05", // Monday
			timezone:    "America/Denver",
			wantStart:   "2025-12-22",
			wantEnd:     "2025-12-28",
			wantErr:     false,
		},
		{
			name:        "Year boundary - New Year's Day (Thursday)",
			currentDate: "2026-01-01", // Thursday
			timezone:    "America/Denver",
			wantStart:   "2025-12-22",
			wantEnd:     "2025-12-28",
			wantErr:     false,
		},
		{
			name:        "Invalid timezone",
			currentDate: "2025-12-15",
			timezone:    "Invalid/Timezone",
			wantStart:   "",
			wantEnd:     "",
			wantErr:     true,
		},
		{
			name:        "Different timezone - UTC",
			currentDate: "2025-12-15", // Monday
			timezone:    "UTC",
			wantStart:   "2025-12-08",
			wantEnd:     "2025-12-14",
			wantErr:     false,
		},
		{
			name:        "Different timezone - America/New_York",
			currentDate: "2025-12-15", // Monday
			timezone:    "America/New_York",
			wantStart:   "2025-12-01",
			wantEnd:     "2025-12-07",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the test date
			testDate, err := time.Parse("2006-01-02", tt.currentDate)
			if err != nil {
				t.Fatalf("Failed to parse test date: %v", err)
			}

			// Call the function
			gotStart, gotEnd, err := calculateLastWeekDates(testDate, tt.timezone)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateLastWeekDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, we're done
			if tt.wantErr {
				return
			}

			// Check start date
			if gotStart != tt.wantStart {
				t.Errorf("calculateLastWeekDates() gotStart = %v, want %v", gotStart, tt.wantStart)
			}

			// Check end date
			if gotEnd != tt.wantEnd {
				t.Errorf("calculateLastWeekDates() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}

			// Additional validation: end date should be 6 days after start date
			startTime, _ := time.Parse("2006-01-02", gotStart)
			endTime, _ := time.Parse("2006-01-02", gotEnd)
			daysDiff := endTime.Sub(startTime).Hours() / 24
			if daysDiff != 6 {
				t.Errorf("End date should be 6 days after start date, got %v days", daysDiff)
			}

			// Additional validation: start date should be a Monday
			if startTime.Weekday() != time.Monday {
				t.Errorf("Start date should be a Monday, got %v", startTime.Weekday())
			}

			// Additional validation: end date should be a Sunday
			if endTime.Weekday() != time.Sunday {
				t.Errorf("End date should be a Sunday, got %v", endTime.Weekday())
			}
		})
	}
}

func TestCalculateLastWeekDates_ConsistentResults(t *testing.T) {
	// Test that calling the function multiple times with the same input gives the same result
	testDate := time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC)

	start1, end1, err1 := calculateLastWeekDates(testDate, "America/Denver")
	start2, end2, err2 := calculateLastWeekDates(testDate, "America/Denver")

	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: err1=%v, err2=%v", err1, err2)
	}

	if start1 != start2 {
		t.Errorf("Inconsistent start dates: %v != %v", start1, start2)
	}

	if end1 != end2 {
		t.Errorf("Inconsistent end dates: %v != %v", end1, end2)
	}
}

func TestCalculateLastWeekDates_TimeOfDayDoesNotMatter(t *testing.T) {
	// Test that the time of day doesn't affect the result (only the date matters)
	date := "2025-12-15"
	timezone := "America/Denver"

	// Test at midnight
	midnight, _ := time.Parse("2006-01-02 15:04:05", date+" 00:00:00")
	startMidnight, endMidnight, err := calculateLastWeekDates(midnight, timezone)
	if err != nil {
		t.Fatalf("Unexpected error at midnight: %v", err)
	}

	// Verify the function returns valid Monday-Sunday pairs
	startT, _ := time.Parse("2006-01-02", startMidnight)
	endT, _ := time.Parse("2006-01-02", endMidnight)

	if startT.Weekday() != time.Monday {
		t.Errorf("Start date should be Monday, got %v", startT.Weekday())
	}
	if endT.Weekday() != time.Sunday {
		t.Errorf("End date should be Sunday, got %v", endT.Weekday())
	}
}
