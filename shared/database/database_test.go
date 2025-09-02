package database

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
)

const testID = 1
const testFirstName = "admin"
const testLastName = "admin"

func TestMain(m *testing.M) {
	// Ensure we're using a test database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" || !strings.Contains(strings.ToLower(dbURL), "test") {
		fmt.Println("WARNING: Tests should use a test database. Set DATABASE_URL to a test database URL containing 'test'")
		fmt.Println("Skipping database tests to avoid affecting production data.")
		os.Exit(0)
	}

	// Setup test database connection
	err := Initialize()
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()

	os.Exit(m.Run())
}

func Test_GetEmployeeID(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantID    int
		wantErr   bool
	}{
		{
			name:      "existing employee",
			firstName: testFirstName,
			lastName:  testLastName,
			wantID:    testID,
			wantErr:   false,
		},
		{
			name:      "non-existent employee",
			firstName: "nonexistent",
			lastName:  "person",
			wantID:    0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GetEmployeeID(tt.firstName, tt.lastName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}
		})
	}
}

func Test_SubmitTime(t *testing.T) {
	// Clean up any existing time entry for the test employee first
	// Get today's date in Mountain Time (same as SubmitTime function)
	loc, _ := time.LoadLocation("America/Denver")
	today := time.Now().In(loc).Truncate(24 * time.Hour)

	_, err := db.Exec("DELETE FROM payroll WHERE id = $1 AND date = $2", testID, today)
	require.NoError(t, err)

	// Setup test data
	hours := 8.5
	message := "Test time submission"

	result, err := SubmitTime(testID, hours, message)
	require.NoError(t, err)
	assert.Contains(t, result, "Submitted hours: 8.50")

	// Test duplicate submission
	newHours := 9.0
	result, err = SubmitTime(testID, newHours, message)
	require.NoError(t, err)
	assert.Contains(t, result, "Updated hours: 8.50 to 9.00")
}

func Test_AddEmployee(t *testing.T) {
	// Use a unique name to avoid conflicts with existing data
	timestamp := time.Now().UnixNano()
	firstName := "testuser"
	lastName := fmt.Sprintf("unique%d", timestamp)
	email := fmt.Sprintf("test%d@example.com", timestamp)
	phone := "1234567890"

	// Add employee without supervisor
	result, err := AddEmployee(firstName, lastName, email, phone, "", "")
	require.NoError(t, err)
	assert.Contains(t, result, fmt.Sprintf("%s %s was successfully added",
		helpers.Title.String(firstName), helpers.Title.String(lastName)))

	// Verify employee was added
	id, err := GetEmployeeID(firstName, lastName)
	require.NoError(t, err)
	assert.NotZero(t, id)

	// Clean up - use defer to ensure cleanup happens even if test fails
	defer func() {
		// Clean up payroll records first to avoid foreign key constraint
		_, err = db.Exec("DELETE FROM payroll WHERE id = $1", id)
		if err != nil {
			t.Logf("Failed to clean up payroll records: %v", err)
		}

		_, err = db.Exec("DELETE FROM employees WHERE id = $1", id)
		if err != nil {
			t.Logf("Failed to clean up employee record: %v", err)
		}
	}()
}

func Test_UpdateEmployee(t *testing.T) {
	// Setup test employee
	firstName := "update"
	lastName := "test"
	_, err := AddEmployee(firstName, lastName, "", "", "", "")
	require.NoError(t, err)

	// Get the employee ID that was just created
	empID, err := GetEmployeeID(firstName, lastName)
	require.NoError(t, err)
	require.NotZero(t, empID)

	// Test update
	newPhone := "9876543210"
	err = UpdateEmployee(firstName, lastName, "phone", newPhone)
	require.NoError(t, err)

	// Verify update
	emp, err := GetEmployee(empID)
	require.NoError(t, err)
	assert.Equal(t, newPhone, emp.Phone)

	// Clean up
	_, err = db.Exec("DELETE FROM employees WHERE first_name = $1 AND last_name = $2",
		firstName, lastName)
	require.NoError(t, err)
}

func Test_GetTimeCards(t *testing.T) {
	// Setup test dates - use date-only format to avoid time component issues
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day()-7, 0, 0, 0, 0, time.UTC)
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)

	entries, err := GetTimeCards(start, end)
	require.NoError(t, err)

	// Verify entries are within date range
	for _, entry := range entries {
		// Compare just the date part by formatting as strings
		entryDateStr := entry.Date.Format("2006-01-02")
		startDateStr := start.Format("2006-01-02")
		endDateStr := end.Format("2006-01-02")

		entryDate, _ := time.Parse("2006-01-02", entryDateStr)
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		endDate, _ := time.Parse("2006-01-02", endDateStr)

		withinRange := !entryDate.Before(startDate) && !entryDate.After(endDate)
		assert.True(t, withinRange,
			"Entry date %s should be between %s and %s (inclusive)",
			entryDateStr, startDateStr, endDateStr)
	}
}
