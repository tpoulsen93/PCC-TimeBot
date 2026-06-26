package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tpoulsen/pcc-timebot/internal/auth"
	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/timecalc"
	"github.com/tpoulsen/pcc-timebot/shared/timecard"
)

const dateLayout = "2006-01-02"

// TimecardHandler handles time submission and timecard viewing for the
// authenticated employee.
type TimecardHandler struct{}

// NewTimecardHandler creates a new TimecardHandler.
func NewTimecardHandler() *TimecardHandler {
	return &TimecardHandler{}
}

type submitHoursBody struct {
	Date     string `json:"date"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Lunch    string `json:"lunch"`
	Extra    string `json:"extra"`
	Location string `json:"location"`
}

// SubmitHours records hours for the authenticated employee. The employee ID is
// always taken from the session, never from the request body.
func (h *TimecardHandler) SubmitHours(c *gin.Context) {
	employee, ok := auth.CurrentEmployee(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var body submitHoursBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	date, err := time.Parse(dateLayout, body.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected YYYY-MM-DD"})
		return
	}

	lunch := body.Lunch
	if lunch == "" {
		lunch = "0"
	}
	extra := body.Extra
	if extra == "" {
		extra = "0"
	}

	hours, err := timecalc.CalculateTime(body.Start, body.End, lunch, extra)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.AddTime(employee.ID, date, hours, body.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit hours"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": result,
		"hours":   hours,
		"date":    body.Date,
	})
}

// GetHistory returns the authenticated employee's submissions within an
// optional date range. Defaults to the last 60 days.
func (h *TimecardHandler) GetHistory(c *gin.Context) {
	employee, ok := auth.CurrentEmployee(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	loc, _ := time.LoadLocation("America/Denver")
	now := time.Now().In(loc)

	start, end := parseRange(c, now.AddDate(0, 0, -60), now)

	entries, err := database.GetPayrollForEmployee(employee.ID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load history"})
		return
	}

	items := make([]gin.H, 0, len(entries))
	for _, e := range entries {
		items = append(items, gin.H{
			"date":     e.Date.In(loc).Format(dateLayout),
			"hours":    e.Time,
			"location": e.Location,
			"message":  e.Message,
		})
	}

	c.JSON(http.StatusOK, gin.H{"entries": items})
}

// GetWeeklySummary returns a weekly timecard for the authenticated employee.
// The week starts on the date given by ?weekStart=YYYY-MM-DD (defaults to the
// most recent Sunday) and covers 7 days.
func (h *TimecardHandler) GetWeeklySummary(c *gin.Context) {
	employee, ok := auth.CurrentEmployee(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	loc, _ := time.LoadLocation("America/Denver")
	now := time.Now().In(loc)

	var weekStart time.Time
	if ws := c.Query("weekStart"); ws != "" {
		parsed, err := time.ParseInLocation(dateLayout, ws, loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid weekStart"})
			return
		}
		weekStart = parsed
	} else {
		// Most recent Sunday.
		weekStart = now.AddDate(0, 0, -int(now.Weekday()))
	}
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, loc)
	weekEnd := weekStart.AddDate(0, 0, 6)

	tc, err := timecard.NewTimeCard(employee.ID, weekStart.Format(dateLayout), weekEnd.Format(dateLayout))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build timecard"})
		return
	}

	entries, err := database.GetPayrollForEmployee(employee.ID, weekStart, weekEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load entries"})
		return
	}
	for _, e := range entries {
		localDate := e.Date.In(loc).Format(dateLayout)
		_ = tc.AddHours(localDate, e.Time, e.Location)
	}

	c.JSON(http.StatusOK, timecardToJSON(tc, weekStart, weekEnd))
}

// parseRange reads optional ?start= and ?end= query params (YYYY-MM-DD),
// falling back to the provided defaults.
func parseRange(c *gin.Context, defStart, defEnd time.Time) (time.Time, time.Time) {
	loc := defStart.Location()
	start := defStart
	end := defEnd
	if s := c.Query("start"); s != "" {
		if parsed, err := time.ParseInLocation(dateLayout, s, loc); err == nil {
			start = parsed
		}
	}
	if e := c.Query("end"); e != "" {
		if parsed, err := time.ParseInLocation(dateLayout, e, loc); err == nil {
			end = parsed
		}
	}
	return start, end
}

// timecardToJSON serializes a timecard with its days sorted ascending.
func timecardToJSON(tc *timecard.TimeCard, start, end time.Time) gin.H {
	dates := make([]string, 0, len(tc.Days))
	for date := range tc.Days {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	days := make([]gin.H, 0, len(dates))
	for _, date := range dates {
		entry := tc.Days[date]
		days = append(days, gin.H{
			"date":     date,
			"hours":    entry.Hours,
			"location": entry.Location,
		})
	}

	return gin.H{
		"name":       tc.Name,
		"weekStart":  start.Format(dateLayout),
		"weekEnd":    end.Format(dateLayout),
		"days":       days,
		"totalHours": tc.TotalHours,
		"payday":     tc.PayDay.Format(dateLayout),
	}
}
