package timecard

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeCard_AddHours(t *testing.T) {
	// Create a timecard for testing
	tc := &TimeCard{
		ID:         1,
		Name:       "John Doe",
		Email:      "john@example.com",
		Phone:      "555-1234",
		Days:       make(map[string]DayEntry),
		TotalHours: 0,
		PayDay:     time.Now(),
	}

	// Initialize some days
	tc.Days["2025-01-01"] = DayEntry{Hours: 0, Location: ""}
	tc.Days["2025-01-02"] = DayEntry{Hours: 0, Location: ""}

	// Test adding hours to valid date
	err := tc.AddHours("2025-01-01", 8.5, "Main Street")
	assert.NoError(t, err)
	assert.Equal(t, 8.5, tc.Days["2025-01-01"].Hours)
	assert.Equal(t, "Main Street", tc.Days["2025-01-01"].Location)
	assert.Equal(t, 8.5, tc.TotalHours)

	// Test adding more hours to same date
	err = tc.AddHours("2025-01-01", 2.0, "Oak Avenue")
	assert.NoError(t, err)
	assert.Equal(t, 10.5, tc.Days["2025-01-01"].Hours)
	assert.Equal(t, "Oak Avenue", tc.Days["2025-01-01"].Location)
	assert.Equal(t, 10.5, tc.TotalHours)

	// Test adding hours to another date
	err = tc.AddHours("2025-01-02", 7.0, "")
	assert.NoError(t, err)
	assert.Equal(t, 7.0, tc.Days["2025-01-02"].Hours)
	assert.Equal(t, 17.5, tc.TotalHours)

	// Test adding hours to invalid date
	err = tc.AddHours("2025-01-03", 5.0, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in pay period")
	assert.Equal(t, 17.5, tc.TotalHours) // Should not change
}

func TestTimeCard_String(t *testing.T) {
	// Create a timecard with some data
	payday := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	tc := &TimeCard{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "555-1234",
		Days: map[string]DayEntry{
			"2025-01-01": {Hours: 8.0, Location: "Main Office"},
			"2025-01-02": {Hours: 7.5, Location: ""},
			"2025-01-03": {Hours: 0.0, Location: ""},
		},
		TotalHours: 15.5,
		PayDay:     payday,
	}

	result := tc.String()

	// Check that the result contains expected elements
	assert.Contains(t, result, "John Doe")
	assert.Contains(t, result, "Total hours:  15.50")
	assert.Contains(t, result, "Payday:  2025-01-15")

	// Check that dates are included (order may vary due to sorting)
	lines := strings.Split(result, "\n")
	dataLines := []string{}
	for _, line := range lines {
		// Only count lines that have the date format and contain a pipe (data rows)
		if strings.Contains(line, "2025-01-") && strings.Contains(line, "|") && !strings.Contains(line, "Date") {
			dataLines = append(dataLines, line)
		}
	}

	// Should have 3 data lines
	assert.Len(t, dataLines, 3)

	// Check that the table format is correct
	assert.Contains(t, result, "Date")
	assert.Contains(t, result, "Day")
	assert.Contains(t, result, "Hours")
}

func TestTimeCard_String_Empty(t *testing.T) {
	// Create an empty timecard
	payday := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	tc := &TimeCard{
		ID:         1,
		Name:       "Jane Smith",
		Email:      "jane@example.com",
		Phone:      "555-5678",
		Days:       map[string]DayEntry{},
		TotalHours: 0,
		PayDay:     payday,
	}

	result := tc.String()

	assert.Contains(t, result, "Jane Smith")
	assert.Contains(t, result, "Total hours:  0.00")
	assert.Contains(t, result, "Payday:  2025-01-15")
}

func TestComputePayday_IsSecondFridayAfterPeriodEnd(t *testing.T) {
	t.Run("period ends Sunday -> payday is second Friday", func(t *testing.T) {
		// Period end: Sunday 2025-12-14
		end := time.Date(2025, 12, 14, 0, 0, 0, 0, time.UTC)
		payday := computePayday(end)

		assert.Equal(t, time.Friday, payday.Weekday())
		assert.Equal(t, "2025-12-26", payday.Format("2006-01-02"))
	})

	t.Run("period ends Saturday -> payday is second Friday", func(t *testing.T) {
		// Period end: Saturday 2025-12-13
		end := time.Date(2025, 12, 13, 0, 0, 0, 0, time.UTC)
		payday := computePayday(end)

		assert.Equal(t, time.Friday, payday.Weekday())
		assert.Equal(t, "2025-12-26", payday.Format("2006-01-02"))
	})

	t.Run("period ends Friday -> payday is second Friday", func(t *testing.T) {
		// Period end: Friday 2025-12-12
		end := time.Date(2025, 12, 12, 0, 0, 0, 0, time.UTC)
		payday := computePayday(end)

		assert.Equal(t, time.Friday, payday.Weekday())
		assert.Equal(t, "2025-12-26", payday.Format("2006-01-02"))
	})
}

// Note: NewTimeCard cannot be easily unit tested because it depends on database calls.
// Consider creating an interface for database operations to enable mocking in tests.
