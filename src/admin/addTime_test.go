package admin

import (
	"os"
	"testing"
	"time"

	"github.com/tpoulsen/pcc-timebot/src/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddTime(t *testing.T) {
	// Ensure database connection is available
	err := database.Initialize()
	require.NoError(t, err, "failed to initialize database")

	// Create a test employee
	result, err := database.AddEmployee("test", "employee", "", "", "", "")
	require.NoError(t, err)
	assert.Contains(t, result, "Test Employee was successfully added")

	// Get the test employee's ID
	id, err := database.GetEmployeeID("test", "employee")
	require.NoError(t, err)
	require.NotZero(t, id)

	// Test adding time for today
	today := time.Now()
	result, err = database.AddTime(id, today, 8.5)
	require.NoError(t, err)
	assert.Contains(t, result, "Submitted 8.5 hours for Test Employee")

	// Test updating the same day
	result, err = database.AddTime(id, today, 9.0)
	require.NoError(t, err)
	assert.Contains(t, result, "Updated submission for Test Employee from 8.5 to 9.0 hours")

	// Clean up test data
	_, err = database.GetEmployee(id)
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	// Setup
	if err := database.Initialize(); err != nil {
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup can be added here if needed

	os.Exit(code)
}
