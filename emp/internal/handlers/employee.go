// internal/handlers/employee.go

package handlers

import (
	"emp/internal/database"
	"emp/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetEmployees(c *fiber.Ctx) error {
	employees, err := models.GetAllEmployees(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Render("employees", fiber.Map{
		"Title":     "All Employees",
		"Employees": employees,
	})
}

func NewEmployee(c *fiber.Ctx) error {
	return c.Render("employee_form", fiber.Map{
		"Title": "New Employee",
	})
}

func CreateEmployee(c *fiber.Ctx) error {
	e := new(models.Employee)
	if err := c.BodyParser(e); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := e.Create(database.DB); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/employees")
}

func GetEmployee(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	employee, err := models.GetEmployee(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("employee_form", fiber.Map{
		"Title":    "Edit Employee",
		"Employee": employee,
	})
}

func EditEmployee(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	employee, err := models.GetEmployee(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("employee_form", fiber.Map{
		"Title":    "Edit Employee",
		"Employee": employee,
	})
}

func UpdateEmployee(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	e := new(models.Employee)
	if err := c.BodyParser(e); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	e.ID = id

	if err := e.Update(database.DB); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/employees")
}

func DeleteEmployee(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	if err := models.DeleteEmployee(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("Employee deleted successfully")
}
