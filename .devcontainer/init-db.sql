-- Database initialization script for PCC-TimeBot development environment

-- Create the timebot user
CREATE USER timebot_user WITH PASSWORD 'timebot_password';

-- Create the employees table
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    supervisor_id INTEGER REFERENCES employees(id),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the payroll table
CREATE TABLE IF NOT EXISTS payroll (
    transaction_id SERIAL PRIMARY KEY,
    id INTEGER NOT NULL REFERENCES employees(id),
    time DECIMAL(4,2) NOT NULL,
    date DATE NOT NULL,
    message TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT submission UNIQUE(id, date)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_employees_name ON employees(first_name, last_name);
CREATE INDEX IF NOT EXISTS idx_employees_supervisor ON employees(supervisor_id);
CREATE INDEX IF NOT EXISTS idx_payroll_employee ON payroll(id);
CREATE INDEX IF NOT EXISTS idx_payroll_date ON payroll(date);
CREATE INDEX IF NOT EXISTS idx_payroll_employee_date ON payroll(id, date);

-- Insert some sample data for development
INSERT INTO employees (first_name, last_name, phone, email) VALUES
    ('admin', 'admin', '555-0001', 'admin@pcc-timebot.dev'),
    ('john', 'doe', '555-0002', 'john.doe@pcc-timebot.dev'),
    ('jane', 'smith', '555-0003', 'jane.smith@pcc-timebot.dev')
ON CONFLICT DO NOTHING;

-- Insert some sample supervisor relationships
UPDATE employees SET supervisor_id = 1 WHERE first_name = 'john' AND last_name = 'doe';
UPDATE employees SET supervisor_id = 1 WHERE first_name = 'jane' AND last_name = 'smith';

-- Insert some sample payroll data
INSERT INTO payroll (id, time, date, message) VALUES
    (2, 8.0, CURRENT_DATE - INTERVAL '1 day', 'Regular workday'),
    (3, 8.5, CURRENT_DATE - INTERVAL '1 day', 'Regular workday with overtime'),
    (2, 7.5, CURRENT_DATE - INTERVAL '2 days', 'Left early'),
    (3, 8.0, CURRENT_DATE - INTERVAL '2 days', 'Regular workday')
ON CONFLICT (id, date) DO NOTHING;

-- Grant necessary permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO timebot_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO timebot_user;
