package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
)

// Employee represents an employee record in the database
type Employee struct {
	ID           int
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	SupervisorID *int
	Timestamp    time.Time
}

// PayrollEntry represents a payroll record in the database
type PayrollEntry struct {
	ID         int `db:"transaction_id"`
	EmployeeID int `db:"id"`
	Time       float64
	Date       time.Time
	Message    string
	Location   string
	Timestamp  time.Time
}

var db *sql.DB

// Initialize sets up the database connection
func Initialize() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable not set")
	}

	// Replace postgres:// with postgresql:// in the connection string if needed
	dbURL = strings.Replace(dbURL, "postgres://", "postgresql://", 1)

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return db.Ping()
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// GetSupervisorID returns the supervisor ID for a given employee
func GetSupervisorID(employeeID int) (int, error) {
	var supervisorID int
	err := db.QueryRow("SELECT supervisor_id FROM employees WHERE id = $1", employeeID).Scan(&supervisorID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get supervisor ID: %w", err)
	}
	return supervisorID, nil
}

// GetEmployeeID returns the ID for an employee given their first and last name
func GetEmployeeID(firstName, lastName string) (int, error) {
	var id int
	err := db.QueryRow(
		"SELECT id FROM employees WHERE first_name = $1 AND last_name = $2",
		strings.ToLower(firstName),
		strings.ToLower(lastName),
	).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get employee ID: %w", err)
	}
	return id, nil
}

// GetEmployeeName returns the full name of an employee given their ID
func GetEmployeeName(id int) (string, error) {
	var firstName, lastName string
	err := db.QueryRow(
		"SELECT first_name, last_name FROM employees WHERE id = $1",
		id,
	).Scan(&firstName, &lastName)
	if err != nil {
		return "", fmt.Errorf("failed to get employee name: %w", err)
	}
	return fmt.Sprintf("%s %s", helpers.Title.String(firstName), helpers.Title.String(lastName)), nil
}

// GetEmployeePhone returns the phone number for an employee
func GetEmployeePhone(id int) (string, error) {
	var phone string
	err := db.QueryRow("SELECT phone FROM employees WHERE id = $1", id).Scan(&phone)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get employee phone: %w", err)
	}
	return phone, nil
}

// GetEmployeeEmail returns the email address for an employee
func GetEmployeeEmail(id int) (string, error) {
	var email string
	err := db.QueryRow("SELECT email FROM employees WHERE id = $1", id).Scan(&email)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get employee email: %w", err)
	}
	return email, nil
}

// GetEmployee returns all information about an employee
func GetEmployee(id int) (*Employee, error) {
	employee := &Employee{}
	var phone, email sql.NullString

	err := db.QueryRow(`
		SELECT id, first_name, last_name, phone, email, supervisor_id, timestamp
		FROM employees WHERE id = $1`,
		id,
	).Scan(
		&employee.ID,
		&employee.FirstName,
		&employee.LastName,
		&phone,
		&email,
		&employee.SupervisorID,
		&employee.Timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Convert sql.NullString to regular strings
	employee.Phone = phone.String
	employee.Email = email.String

	return employee, nil
}

// DuplicateSubmission checks if there's already a time submission for the given date
func DuplicateSubmission(id int, date time.Time) (float64, error) {
	var hours float64
	err := db.QueryRow(
		"SELECT time FROM payroll WHERE id = $1 AND date = $2",
		id, date,
	).Scan(&hours)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to check for duplicate submission: %w", err)
	}
	return hours, nil
}

// SubmitTime adds a new time entry or updates an existing one
func SubmitTime(id int, hours float64, message string, location string) (string, error) {
	// Convert current time to Mountain Time (UTC-7)
	loc, _ := time.LoadLocation("America/Denver")
	today := time.Now().In(loc).Truncate(24 * time.Hour)

	// Check for duplicate submission
	existingHours, err := DuplicateSubmission(id, today)
	if err != nil {
		return "", err
	}

	var result string
	if existingHours == 0 {
		result = fmt.Sprintf("Submitted hours: %.2f", hours)
	} else {
		result = fmt.Sprintf("Updated hours: %.2f to %.2f", existingHours, hours)
	}

	// Use upsert to either insert new record or update existing one
	_, err = db.Exec(`
		INSERT INTO payroll (id, time, date, message, location)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT ON CONSTRAINT submission
		DO UPDATE SET time = EXCLUDED.time, message = EXCLUDED.message, location = EXCLUDED.location`,
		id, hours, today, message, location,
	)
	if err != nil {
		return "", fmt.Errorf("failed to submit time: %w", err)
	}

	return result, nil
}

// AddTime manually adds or updates time for an employee on a specific date
func AddTime(id int, date time.Time, hours float64, location string) (string, error) {
	name, err := GetEmployeeName(id)
	if err != nil {
		return "", err
	}

	existingHours, err := DuplicateSubmission(id, date)
	if err != nil {
		return "", err
	}

	loc, _ := time.LoadLocation("America/Denver")
	today := time.Now().In(loc).Truncate(24 * time.Hour)

	var message string
	var result string
	if existingHours == 0 {
		message = fmt.Sprintf("Submitted manually for %s on %s", name, today.Format("2006-01-02"))
		result = fmt.Sprintf("Submitted %.2f hours for %s on %s", hours, name, date.Format("2006-01-02"))
	} else {
		message = fmt.Sprintf("Updated manually for %s on %s", name, today.Format("2006-01-02"))
		result = fmt.Sprintf("Updated submission for %s from %.2f to %.2f hours on %s",
			name, existingHours, hours, date.Format("2006-01-02"))
	}

	_, err = db.Exec(`
		INSERT INTO payroll (id, time, date, message, location)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT ON CONSTRAINT submission
		DO UPDATE SET time = EXCLUDED.time, message = EXCLUDED.message, location = EXCLUDED.location`,
		id, hours, date, message, location,
	)
	if err != nil {
		return "", fmt.Errorf("failed to add time: %w", err)
	}

	return result, nil
}

// AddEmployee creates a new employee record
func AddEmployee(firstName, lastName string, email, phone string, superFirstName, superLastName string) (string, error) {
	var supervisorID *int

	if superFirstName != "" && superLastName != "" {
		id, err := GetEmployeeID(superFirstName, superLastName)
		if err != nil {
			return "", err
		}
		if id == 0 {
			return "", fmt.Errorf("supervisor not found")
		}
		supervisorID = &id
	}

	// Normalize inputs
	firstName = strings.ToLower(firstName)
	lastName = strings.ToLower(lastName)
	if email != "" {
		email = strings.ToLower(email)
	}

	_, err := db.Exec(`
		INSERT INTO employees (first_name, last_name, supervisor_id, phone, email, timestamp)
		VALUES ($1, $2, $3, $4, $5, NOW())`,
		firstName, lastName, supervisorID,
		sql.NullString{String: phone, Valid: phone != ""},
		sql.NullString{String: email, Valid: email != ""},
	)
	if err != nil {
		return "", fmt.Errorf("failed to add employee: %w", err)
	}

	return fmt.Sprintf("%s %s was successfully added",
		helpers.Title.String(firstName), helpers.Title.String(lastName)), nil
}

// UpdateEmployee updates an employee's information
func UpdateEmployee(firstName, lastName, field, value string) error {
	id, err := GetEmployeeID(firstName, lastName)
	if err != nil {
		return err
	}
	if id == 0 {
		return fmt.Errorf("employee not found")
	}

	query := fmt.Sprintf("UPDATE employees SET %s = $1 WHERE id = $2", field)
	_, err = db.Exec(query, value, id)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

// GetTimeCards retrieves all time card entries between start and end dates
func GetTimeCards(start, end time.Time) ([]PayrollEntry, error) {
	rows, err := db.Query(`
		SELECT id, time, date FROM payroll
		WHERE date >= $1 AND date <= $2
		ORDER BY id`,
		start, end,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get time cards: %w", err)
	}
	defer rows.Close()

	var entries []PayrollEntry
	for rows.Next() {
		var entry PayrollEntry
		err := rows.Scan(&entry.EmployeeID, &entry.Time, &entry.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan time card entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
