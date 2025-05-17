package timecalc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tpoulsen/pcc-timebot/src/helpers"
)

var (
	ErrTimeFormat  = errors.New("time formatted incorrectly")
	ErrHours       = errors.New("hours spot is wrong")
	ErrMeridiem    = errors.New("meridiem is wrong (am/pm)")
	ErrMinutes     = errors.New("minutes spot is wrong")
	ErrIllegalTime = errors.New("end time is before start time")
	ErrLunch       = errors.New("subtracted hours formatted incorrectly")
	ErrExtra       = errors.New("additional hours formatted incorrectly")
)

// CalculateTime calculates the total hours worked based on start time, end time, lunch break, and extra time
func CalculateTime(start, end, less, more string) (float64, error) {
	startTime, err := getDuration(start)
	if err != nil {
		return 0, fmt.Errorf("invalid start time: %w", err)
	}

	endTime, err := getDuration(end)
	if err != nil {
		return 0, fmt.Errorf("invalid end time: %w", err)
	}

	if endTime < startTime {
		return 0, ErrIllegalTime
	}

	subtract, err := strconv.ParseFloat(less, 64)
	if err != nil {
		return 0, ErrLunch
	}

	add, err := strconv.ParseFloat(more, 64)
	if err != nil {
		return 0, ErrExtra
	}

	// Calculate total hours
	totalDuration := endTime - startTime
	totalHours := totalDuration.Hours()

	// Apply lunch break subtraction and extra time addition
	totalHours = totalHours - subtract + add

	return helpers.Round(totalHours, 2), nil
}

// getDuration converts a time string (like "9:00am" or "9am") into a time.Duration
func getDuration(timeStr string) (time.Duration, error) {
	timeStr = strings.ToLower(timeStr)
	var hours, minutes int

	if strings.Contains(timeStr, ":") { // Format: "9:00am"
		if len(timeStr) < 6 || len(timeStr) > 7 {
			return 0, ErrTimeFormat
		}

		parts := strings.Split(timeStr, ":")
		var err error
		hours, err = strconv.Atoi(parts[0])
		if err != nil || hours < 1 || hours > 12 {
			return 0, ErrHours
		}

		// Handle AM/PM and minutes
		minutesPart := parts[1]
		if strings.HasSuffix(minutesPart, "am") {
			if hours == 12 {
				hours = 0
			}
			minutesPart = strings.TrimSuffix(minutesPart, "am")
		} else if strings.HasSuffix(minutesPart, "pm") {
			if hours != 12 {
				hours += 12
			}
			minutesPart = strings.TrimSuffix(minutesPart, "pm")
		} else {
			return 0, ErrMeridiem
		}

		minutes, err = strconv.Atoi(minutesPart)
		if err != nil || minutes < 0 || minutes > 59 {
			return 0, ErrMinutes
		}
	} else { // Format: "9am"
		if len(timeStr) < 3 || len(timeStr) > 4 {
			return 0, ErrTimeFormat
		}

		minutes = 0
		var numStr string
		if strings.HasSuffix(timeStr, "am") {
			numStr = strings.TrimSuffix(timeStr, "am")
			hours, _ = strconv.Atoi(numStr)
			if hours == 12 {
				hours = 0
			}
		} else if strings.HasSuffix(timeStr, "pm") {
			numStr = strings.TrimSuffix(timeStr, "pm")
			hours, _ = strconv.Atoi(numStr)
			if hours != 12 {
				hours += 12
			}
		} else {
			return 0, ErrMeridiem
		}

		if hours < 1 || hours > 24 {
			return 0, ErrHours
		}
	}

	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute, nil
}
