// internal/handlers/week_overview.go

package handlers

import (
	"emp/internal/database"
	"emp/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func WeekOverview(c *fiber.Ctx) error {
	employeeID, err := strconv.Atoi(c.Params("employeeId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid employee ID")
	}

	// Get the start of the current week (Sunday)
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))

	workdays, err := models.GetWeekWorkdays(database.DB, employeeID, startOfWeek)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	employee, err := models.GetEmployee(database.DB, employeeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Create a map of workdays keyed by date string
	workdayMap := make(map[string]models.Day)
	for _, day := range workdays {
		workdayMap[day.Date.Format("2006-01-02")] = day
	}

	return c.Render("week_overview", fiber.Map{
		"Employee":   employee,
		"WorkdayMap": workdayMap,
		"WeekStart":  startOfWeek,
	})
}

func WeeklyOverview(c *fiber.Ctx) error {
	startDate := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))

	employeeWeekData, err := models.GetEmployeeWeekData(database.DB, startDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Create a slice of weekdays
	weekdays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	data := fiber.Map{
		"StartDate":        startDate,
		"EmployeeWeekData": employeeWeekData,
		"Weekdays":         weekdays,
	}

	return c.Render("weekly_overview", data)
}
