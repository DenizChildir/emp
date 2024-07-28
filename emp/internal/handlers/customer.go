// internal/handlers/customer.go

package handlers

import (
	"emp/internal/database"
	"emp/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetCustomers(c *fiber.Ctx) error {
	customers, err := models.GetAllCustomers(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Render("customers", fiber.Map{
		"Title":     "All Customers",
		"Customers": customers,
	})
}

func NewCustomer(c *fiber.Ctx) error {
	return c.Render("customer_form", fiber.Map{
		"Title": "New Customer",
	})
}

func CreateCustomer(c *fiber.Ctx) error {
	cust := new(models.Customer)
	if err := c.BodyParser(cust); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := cust.Create(database.DB); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/customers")
}

func GetCustomer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	customer, err := models.GetCustomer(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("customer_form", fiber.Map{
		"Title":    "Edit Customer",
		"Customer": customer,
	})
}

func EditCustomer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	customer, err := models.GetCustomer(database.DB, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("customer_form", fiber.Map{
		"Title":    "Edit Customer",
		"Customer": customer,
	})
}

func UpdateCustomer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	cust := new(models.Customer)
	if err := c.BodyParser(cust); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	cust.ID = id

	if err := cust.Update(database.DB); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/customers")
}

func DeleteCustomer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	if err := models.DeleteCustomer(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("Customer deleted successfully")
}
