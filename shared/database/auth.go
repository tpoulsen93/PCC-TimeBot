package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// GetEmployeeByEmail looks up an employee by email (case-insensitive).
// Returns (nil, nil) when no employee matches, so callers can distinguish
// "not found" from a real error without leaking which emails exist.
func GetEmployeeByEmail(email string) (*Employee, error) {
	employee := &Employee{}
	var phone, dbEmail sql.NullString

	err := db.QueryRow(`
		SELECT id, first_name, last_name, phone, email, supervisor_id, is_admin, timestamp
		FROM employees WHERE lower(email) = lower($1)`,
		strings.TrimSpace(email),
	).Scan(
		&employee.ID,
		&employee.FirstName,
		&employee.LastName,
		&phone,
		&dbEmail,
		&employee.SupervisorID,
		&employee.IsAdmin,
		&employee.Timestamp,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get employee by email: %w", err)
	}

	employee.Phone = phone.String
	employee.Email = dbEmail.String
	return employee, nil
}

// ListEmployees returns all employees ordered by last then first name.
func ListEmployees() ([]Employee, error) {
	rows, err := db.Query(`
		SELECT id, first_name, last_name, phone, email, supervisor_id, is_admin, timestamp
		FROM employees
		ORDER BY last_name, first_name`)
	if err != nil {
		return nil, fmt.Errorf("failed to list employees: %w", err)
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		var phone, email sql.NullString
		if err := rows.Scan(
			&e.ID, &e.FirstName, &e.LastName, &phone, &email,
			&e.SupervisorID, &e.IsAdmin, &e.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		e.Phone = phone.String
		e.Email = email.String
		employees = append(employees, e)
	}
	return employees, nil
}

// CountEmployees returns the total number of employee records.
// Used by the bootstrap tool to decide whether a first admin is needed.
func CountEmployees() (int, error) {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count employees: %w", err)
	}
	return count, nil
}

// CreateEmployee inserts a new employee and returns the created record.
// firstName/lastName/email are normalized to lower case to match lookup behavior.
func CreateEmployee(firstName, lastName, email, phone string, supervisorID *int, isAdmin bool) (*Employee, error) {
	firstName = strings.ToLower(strings.TrimSpace(firstName))
	lastName = strings.ToLower(strings.TrimSpace(lastName))
	email = strings.ToLower(strings.TrimSpace(email))

	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("first and last name are required")
	}

	var id int
	err := db.QueryRow(`
		INSERT INTO employees (first_name, last_name, supervisor_id, phone, email, is_admin, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id`,
		firstName, lastName, supervisorID,
		sql.NullString{String: phone, Valid: phone != ""},
		sql.NullString{String: email, Valid: email != ""},
		isAdmin,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	return GetEmployee(id)
}

// SetEmployeeAdmin sets the is_admin flag for an employee.
func SetEmployeeAdmin(id int, isAdmin bool) error {
	_, err := db.Exec("UPDATE employees SET is_admin = $1 WHERE id = $2", isAdmin, id)
	if err != nil {
		return fmt.Errorf("failed to set admin flag: %w", err)
	}
	return nil
}

// GetPayrollForEmployee returns payroll entries for a single employee within an
// inclusive date range, most recent first.
func GetPayrollForEmployee(id int, start, end time.Time) ([]PayrollEntry, error) {
	rows, err := db.Query(`
		SELECT transaction_id, id, time, date, COALESCE(message, ''), COALESCE(location, ''), timestamp
		FROM payroll
		WHERE id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC`,
		id, start, end,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payroll entries: %w", err)
	}
	defer rows.Close()

	var entries []PayrollEntry
	for rows.Next() {
		var e PayrollEntry
		if err := rows.Scan(
			&e.ID, &e.EmployeeID, &e.Time, &e.Date, &e.Message, &e.Location, &e.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payroll entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// CreateLoginToken stores a single-use magic-link token hash for an employee.
func CreateLoginToken(employeeID int, tokenHash string, expiresAt time.Time) error {
	_, err := db.Exec(`
		INSERT INTO login_tokens (employee_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`,
		employeeID, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create login token: %w", err)
	}
	return nil
}

// ConsumeLoginToken atomically validates and consumes a magic-link token.
// It returns the associated employee ID only if the token exists, has not
// expired, and has not already been used. The update marks the token used in
// the same statement to prevent reuse races.
func ConsumeLoginToken(tokenHash string) (int, error) {
	var employeeID int
	err := db.QueryRow(`
		UPDATE login_tokens
		SET used_at = NOW()
		WHERE token_hash = $1
		  AND used_at IS NULL
		  AND expires_at > NOW()
		RETURNING employee_id`,
		tokenHash,
	).Scan(&employeeID)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("invalid or expired token")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to consume login token: %w", err)
	}
	return employeeID, nil
}

// CreateSession stores a new authenticated session.
func CreateSession(id string, employeeID int, expiresAt time.Time) error {
	_, err := db.Exec(`
		INSERT INTO sessions (id, employee_id, expires_at)
		VALUES ($1, $2, $3)`,
		id, employeeID, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

// GetSessionEmployeeID returns the employee ID for a non-expired session,
// refreshes last_seen, and extends expires_at to newExpiry (rolling window).
// Returns (0, nil) when the session is missing or expired.
func GetSessionEmployeeID(id string, newExpiry time.Time) (int, error) {
	var employeeID int
	err := db.QueryRow(`
		UPDATE sessions
		SET last_seen = NOW(), expires_at = $2
		WHERE id = $1 AND expires_at > NOW()
		RETURNING employee_id`,
		id, newExpiry,
	).Scan(&employeeID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get session: %w", err)
	}
	return employeeID, nil
}

// DeleteSession removes a session (logout).
func DeleteSession(id string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
