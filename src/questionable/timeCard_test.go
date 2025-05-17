package timecalc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpoulsen/pcc-timebot/src/helpers"
)

func TestNewTimeCard(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		start      string
		end        string
		wantErr    bool
		wantDays   int
		wantPayDay string
	}{
		{
			name:       "valid date range",
			id:         1,
			start:      "2025-01-01",
			end:        "2025-01-07",
			wantErr:    false,
			wantDays:   7,
			wantPayDay: "2025-01-19",
		},
		{
			name:    "invalid start date",
			id:      1,
			start:   "2025-13-01", // invalid month
			end:     "2025-01-07",
			wantErr: true,
		},
		{
			name:    "invalid end date",
			id:      1,
			start:   "2025-01-01",
			end:     "2025-01-32", // invalid day
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc, err := NewTimeCard(tt.id, tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err, "expected error")
				return
			}

			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tt.wantDays, len(tc.Days), "wrong number of days")
			assert.Equal(t, tt.wantPayDay, tc.PayDay.Format("2006-01-02"), "wrong payday")
		})
	}
}

func TestAddHours(t *testing.T) {
	tc, err := NewTimeCard(1, "2025-01-01", "2025-01-07")
	assert.NoError(t, err, "failed to create timecard")

	tests := []struct {
		date      string
		hours     float64
		wantTotal float64
	}{
		{"2025-01-01", 8.0, 8.0},
		{"2025-01-02", 7.5, 15.5},
		{"2025-01-03", 6.75, 22.25},
	}

	for _, tt := range tests {
		tc.AddHours(tt.date, tt.hours)
		assert.Equal(t, tt.hours, tc.Days[tt.date], "wrong hours for %s", tt.date)
		assert.Equal(t, tt.wantTotal, tc.TotalHours, "wrong total hours after adding %.2f", tt.hours)
	}
}

func TestString(t *testing.T) {
	tc, err := NewTimeCard(1, "2025-01-01", "2025-01-03")
	assert.NoError(t, err, "failed to create timecard")

	tc.Name = "John Doe"
	tc.AddHours("2025-01-01", 8.0)
	tc.AddHours("2025-01-02", 7.5)
	tc.AddHours("2025-01-03", 6.75)

	output := tc.String()

	// Check that the output contains expected elements
	expectedElements := []string{
		"John Doe",
		"Date",
		"Day",
		"Hours",
		"2025-01-01",
		"Wed",
		"8.00",
		"2025-01-02",
		"Thu",
		"7.50",
		"2025-01-03",
		"Fri",
		"6.75",
		"Total hours:  22.25",
		"Payday:  2025-01-15",
	}

	for _, expected := range expectedElements {
		assert.Contains(t, output, expected, "output missing expected element %q", expected)
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		input    float64
		places   int
		expected float64
	}{
		{1.234, 2, 1.23},
		{1.235, 2, 1.24},
		{8.5, 1, 8.5},
		{8.25, 1, 8.3},
		{7.555, 2, 7.56},
	}

	for _, tt := range tests {
		result := helpers.Round(tt.input, tt.places)
		assert.Equal(t, tt.expected, result, "round(%.3f, %d) = %.3f; want %.3f", tt.input, tt.places, result, tt.expected)
	}
}
