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
			name:        "Monday - should return last pay period Sunday to Saturday",
			currentDate: "2025-12-15", // Monday
			timezone:    "America/Denver",
			wantStart:   "2025-12-07",
			wantEnd:     "2025-12-13",
			wantErr:     false,
		},
		{
			name:        "Tuesday - should return last pay period Sunday to Saturday",
			currentDate: "2025-12-16", // Tuesday
			timezone:    "America/Denver",
			wantStart:   "2025-12-07",
			wantEnd:     "2025-12-13",
			wantErr:     false,
		},
		{
			name:        "Wednesday - should return last pay period Sunday to Saturday",
			currentDate: "2025-12-17", // Wednesday
			timezone:    "America/Denver",
			wantStart:   "2025-12-07",
			wantEnd:     "2025-12-13",
			wantErr:     false,
		},
		{
			name:        "Thursday - should return last pay period Sunday to Saturday",
			currentDate: "2025-12-18", // Thursday
			timezone:    "America/Denver",
			wantStart:   "2025-12-07",
			wantEnd:     "2025-12-13",
			wantErr:     false,
		},
		{
			name:        "Friday - should return last pay period Sunday to Saturday",
			currentDate: "2025-12-19", // Friday
			timezone:    "America/Denver",
			wantStart:   "2025-12-07",
			wantEnd:     "2025-12-13",
			wantErr:     false,
		},
		{
			name:        "Saturday - should return last pay period Sunday to Saturday",
			currentDate: "2025-12-20", // Saturday
			timezone:    "America/Denver",
			wantStart:   "2025-12-07",
			wantEnd:     "2025-12-13",
			wantErr:     false,
		},
		{
			name:        "Sunday - should return last pay period Sunday to Saturday",
			currentDate: "2025-12-21", // Sunday
			timezone:    "America/Denver",
			wantStart:   "2025-12-14",
			wantEnd:     "2025-12-20",
			wantErr:     false,
		},
		{
			name:        "First Monday of year",
			currentDate: "2026-01-05", // Monday
			timezone:    "America/Denver",
			wantStart:   "2025-12-28",
			wantEnd:     "2026-01-03",
			wantErr:     false,
		},
		{
			name:        "Year boundary - New Year's Day (Thursday)",
			currentDate: "2026-01-01", // Thursday
			timezone:    "America/Denver",
			wantStart:   "2025-12-21",
			wantEnd:     "2025-12-27",
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
			wantStart:   "2025-12-07",
			wantEnd:     "2025-12-13",
			wantErr:     false,
		},
		{
			name:        "Different timezone - America/New_York",
			currentDate: "2025-12-15", // Monday
			timezone:    "America/New_York",
			wantStart:   "2025-12-07",
			wantEnd:     "2025-12-13",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Construct a time that is definitively on the given local date in the requested timezone.
			// Using date-only parsing in UTC can shift the local date backwards for US timezones.
			var testDate time.Time
			if tt.timezone != "" {
				loc, locErr := time.LoadLocation(tt.timezone)
				if locErr == nil {
					dateOnly, parseErr := time.ParseInLocation("2006-01-02", tt.currentDate, loc)
					if parseErr != nil {
						t.Fatalf("Failed to parse test date in location: %v", parseErr)
					}
					// Midday avoids DST/offset edge cases around local midnight.
					testDate = dateOnly.Add(12 * time.Hour)
				} else {
					// For invalid timezone tests, fall back to a stable UTC value.
					dateOnly, parseErr := time.Parse("2006-01-02", tt.currentDate)
					if parseErr != nil {
						t.Fatalf("Failed to parse test date: %v", parseErr)
					}
					testDate = dateOnly.Add(12 * time.Hour)
				}
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

			// Additional validation: start date should be a Sunday
			if startTime.Weekday() != time.Sunday {
				t.Errorf("Start date should be a Sunday, got %v", startTime.Weekday())
			}

			// Additional validation: end date should be a Saturday
			if endTime.Weekday() != time.Saturday {
				t.Errorf("End date should be a Saturday, got %v", endTime.Weekday())
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
	// Test that the time of day doesn't affect the result (only the *local date* matters)
	date := "2025-12-15"
	timezone := "America/Denver"

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}

	dateOnly, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		t.Fatalf("Failed to parse date in location: %v", err)
	}

	// Same local date, different local times
	localEarly := dateOnly.Add(1 * time.Hour)
	localLate := dateOnly.Add(23 * time.Hour)

	startEarly, endEarly, err := calculateLastWeekDates(localEarly, timezone)
	if err != nil {
		t.Fatalf("Unexpected error for early time: %v", err)
	}

	startLate, endLate, err := calculateLastWeekDates(localLate, timezone)
	if err != nil {
		t.Fatalf("Unexpected error for late time: %v", err)
	}

	if startEarly != startLate {
		t.Errorf("Start date should not depend on time of day: %v != %v", startEarly, startLate)
	}
	if endEarly != endLate {
		t.Errorf("End date should not depend on time of day: %v != %v", endEarly, endLate)
	}

	// Verify returned range is Sunday-Saturday
	startT, _ := time.Parse("2006-01-02", startEarly)
	endT, _ := time.Parse("2006-01-02", endEarly)
	if startT.Weekday() != time.Sunday {
		t.Errorf("Start date should be Sunday, got %v", startT.Weekday())
	}
	if endT.Weekday() != time.Saturday {
		t.Errorf("End date should be Saturday, got %v", endT.Weekday())
	}
}
