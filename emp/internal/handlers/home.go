// internal/handlers/home.go

package handlers

import "github.com/gofiber/fiber/v2"

func Home(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Home - Workforce Management System",
	})
}
