package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		num      float64
		places   int
		expected float64
	}{
		{
			name:     "round to 2 decimal places",
			num:      3.14159,
			places:   2,
			expected: 3.14,
		},
		{
			name:     "round to 0 decimal places",
			num:      3.7,
			places:   0,
			expected: 4.0,
		},
		{
			name:     "round to 3 decimal places",
			num:      2.71828,
			places:   3,
			expected: 2.718,
		},
		{
			name:     "round negative number",
			num:      -1.23456,
			places:   2,
			expected: -1.23,
		},
		{
			name:     "round zero",
			num:      0.0,
			places:   2,
			expected: 0.0,
		},
		{
			name:     "round large number",
			num:      123456.789,
			places:   1,
			expected: 123456.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round(tt.num, tt.places)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRound_EdgeCases(t *testing.T) {
	// Test very small numbers
	assert.Equal(t, 0.0, Round(0.000001, 2))

	// Test very large numbers
	assert.Equal(t, 1000000.0, Round(999999.999, 0))

	// Test rounding at boundary
	assert.Equal(t, 1.0, Round(0.5, 0))
	assert.Equal(t, 2.0, Round(1.5, 0))
	assert.Equal(t, -1.0, Round(-0.5, 0))
	assert.Equal(t, -2.0, Round(-1.5, 0))
}

func TestGlobalWage(t *testing.T) {
	assert.Equal(t, 30.0, GlobalWage)
	assert.IsType(t, 30.0, GlobalWage)
}

func TestTitle(t *testing.T) {
	// Test title casing
	assert.Equal(t, "John", Title.String("john"))
	assert.Equal(t, "Mary Smith", Title.String("mary smith"))
	assert.Equal(t, "A", Title.String("a"))
	assert.Equal(t, "", Title.String(""))
	assert.Equal(t, "Hello World", Title.String("hello world"))
}

// Note: GetUserInput is difficult to unit test because it reads from stdin.
// In a real application, you might want to refactor it to accept an io.Reader
// for better testability, or test it through integration tests.
