package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tpoulsen/pcc-timebot/internal/admin"
	"github.com/tpoulsen/pcc-timebot/shared/database"
	"github.com/tpoulsen/pcc-timebot/shared/helpers"
)

// AdminHandler exposes admin-only employee and payroll management endpoints.
// All routes using it must be protected by the RequireAdmin middleware.
type AdminHandler struct{}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListEmployees returns all employees.
func (h *AdminHandler) ListEmployees(c *gin.Context) {
	employees, err := database.ListEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list employees"})
		return
	}

	items := make([]gin.H, 0, len(employees))
	for i := range employees {
		items = append(items, employeeToJSON(&employees[i]))
	}
	c.JSON(http.StatusOK, gin.H{"employees": items})
}

type createEmployeeBody struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	SupervisorID *int   `json:"supervisorId"`
	IsAdmin      bool   `json:"isAdmin"`
}

// CreateEmployee creates a new employee record.
func (h *AdminHandler) CreateEmployee(c *gin.Context) {
	var body createEmployeeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	employee, err := database.CreateEmployee(
		body.FirstName, body.LastName, body.Email, body.Phone, body.SupervisorID, body.IsAdmin,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employeeToJSON(employee))
}

type updateEmployeeBody struct {
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	SupervisorID *int    `json:"supervisorId"`
	IsAdmin      *bool   `json:"isAdmin"`
}

// UpdateEmployee updates whitelisted fields of an employee by ID.
func (h *AdminHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	var body updateEmployeeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	updates := map[string]string{}
	if body.FirstName != nil {
		updates["first_name"] = *body.FirstName
	}
	if body.LastName != nil {
		updates["last_name"] = *body.LastName
	}
	if body.Email != nil {
		updates["email"] = *body.Email
	}
	if body.Phone != nil {
		updates["phone"] = *body.Phone
	}
	if body.SupervisorID != nil {
		updates["supervisor_id"] = strconv.Itoa(*body.SupervisorID)
	}

	for field, value := range updates {
		if err := database.UpdateEmployeeField(id, field, value); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if body.IsAdmin != nil {
		if err := database.SetEmployeeAdmin(id, *body.IsAdmin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update admin flag"})
			return
		}
	}

	employee, err := database.GetEmployee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load employee"})
		return
	}
	c.JSON(http.StatusOK, employeeToJSON(employee))
}

// GetAllTimecards returns per-employee time cards for a date range. Defaults to
// the current week (most recent Sunday through Saturday).
func (h *AdminHandler) GetAllTimecards(c *gin.Context) {
	start, end := adminRange(c)

	timeCards, payday, err := admin.BuildTimeCards(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build timecards"})
		return
	}

	ids := make([]int, 0, len(timeCards))
	for id := range timeCards {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	var totalHours float64
	cards := make([]gin.H, 0, len(ids))
	for _, id := range ids {
		tc := timeCards[id]
		totalHours += tc.TotalHours
		cards = append(cards, gin.H{
			"employeeId": id,
			"name":       tc.Name,
			"totalHours": tc.TotalHours,
		})
	}

	resp := gin.H{
		"start":      start.Format(dateLayout),
		"end":        end.Format(dateLayout),
		"timecards":  cards,
		"totalHours": totalHours,
		"cost":       totalHours * helpers.GlobalWage,
	}
	if !payday.IsZero() {
		resp["payday"] = payday.Format(dateLayout)
	}
	c.JSON(http.StatusOK, resp)
}

// SendTimecards builds and emails time cards for a date range.
func (h *AdminHandler) SendTimecards(c *gin.Context) {
	start, end := adminRange(c)

	sent, err := admin.SendTimeCardsForRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "time cards sent",
		"sent":    sent,
		"start":   start.Format(dateLayout),
		"end":     end.Format(dateLayout),
	})
}

// adminRange parses ?start= and ?end= query params, defaulting to the current
// week (most recent Sunday through the following Saturday) in Denver time.
func adminRange(c *gin.Context) (time.Time, time.Time) {
	loc, _ := time.LoadLocation("America/Denver")
	now := time.Now().In(loc)
	defStart := now.AddDate(0, 0, -int(now.Weekday()))
	defStart = time.Date(defStart.Year(), defStart.Month(), defStart.Day(), 0, 0, 0, 0, loc)
	defEnd := defStart.AddDate(0, 0, 6)
	return parseRange(c, defStart, defEnd)
}
