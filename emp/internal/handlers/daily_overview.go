// internal/handlers/daily_overview.go

package handlers

import (
	"emp/internal/database"
	"emp/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func DailyOverview(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	var date time.Time
	var err error

	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid date format")
		}
	}

	overview, err := models.GetDailyOverview(database.DB, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("daily_overview", fiber.Map{
		"Date":     date,
		"Overview": overview,
	})
}
