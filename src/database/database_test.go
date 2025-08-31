package database

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testID = 1
const testFirstName = "admin"
const testLastName = "admin"

func TestMain(m *testing.M) {
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
	firstName := "test"
	lastName := "employee"
	email := "test@example.com"
	phone := "1234567890"

	// Add employee without supervisor
	result, err := AddEmployee(firstName, lastName, email, phone, "", "")
	require.NoError(t, err)
	assert.Contains(t, result, "Test Employee was successfully added")

	// Verify employee was added
	id, err := GetEmployeeID(firstName, lastName)
	require.NoError(t, err)
	assert.NotZero(t, id)

	// Clean up
	_, err = db.Exec("DELETE FROM employees WHERE first_name = $1 AND last_name = $2",
		firstName, lastName)
	require.NoError(t, err)
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
	// Setup test dates
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()

	entries, err := GetTimeCards(start, end)
	require.NoError(t, err)

	// Verify entries are within date range
	for _, entry := range entries {
		assert.True(t, !entry.Date.Before(start) && !entry.Date.After(end))
	}
}
