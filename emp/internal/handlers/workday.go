// handler/workday.go
package handlers

import (
	"database/sql"
	"emp/internal/database"
	"emp/internal/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UpdateWorkday(c *fiber.Ctx) error {
	// Print all form values
	c.Context().PostArgs().VisitAll(func(key, value []byte) {
		fmt.Printf("Form field: %s = %s\n", string(key), string(value))
	})

	employeeID, err := strconv.Atoi(c.Params("employeeId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid employee ID")
	}
	fmt.Printf("employeeID: %d\n", employeeID)

	dateStr := c.FormValue("date")
	fmt.Printf("dateStr: %s\n", dateStr)

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid date format: " + err.Error())
	}
	fmt.Printf("parsed date: %v\n", date)

	// Try to get existing workday
	day, err := models.GetWorkday(database.DB, employeeID, date)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting workday: " + err.Error())
	}
	fmt.Printf("existing day before update: %+v\n", day)

	startTimeStr := c.FormValue("startTime")
	endTimeStr := c.FormValue("endTime")
	startMealTimeStr := c.FormValue("startMealTime")
	endMealTimeStr := c.FormValue("endMealTime")
	customerIDStr := c.FormValue("customerId")

	// Validate time formats
	_, err = time.Parse("15:04", startTimeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid start time format: " + err.Error())
	}
	_, err = time.Parse("15:04", endTimeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid end time format: " + err.Error())
	}
	_, err = time.Parse("15:04", startMealTimeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid start meal time format: " + err.Error())
	}
	_, err = time.Parse("15:04", endMealTimeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid end meal time format: " + err.Error())
	}

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid customer ID: " + err.Error())
	}

	startTime, _ := time.Parse("15:04", startTimeStr)
	endTime, _ := time.Parse("15:04", endTimeStr)
	startMealTime, _ := time.Parse("15:04", startMealTimeStr)
	endMealTime, _ := time.Parse("15:04", endMealTimeStr)

	shiftHours := calculateShiftHours(startTime, endTime, startMealTime, endMealTime)

	// Update or create the day struct
	day.Date = date
	day.EmployeeID = employeeID // Ensure this is set correctly
	day.StartTime = startTimeStr
	day.EndTime = endTimeStr
	day.StartMealTime = startMealTimeStr
	day.EndMealTime = endMealTimeStr
	day.CustomerID = customerID
	day.ShiftHours = shiftHours

	fmt.Printf("Final day struct before save: %+v\n", day)

	if day.ID == 0 {
		fmt.Println("Creating new workday")
		err = day.Create(database.DB)
		if err != nil {
			fmt.Printf("Error creating workday: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating workday: " + err.Error())
		}
	} else {
		fmt.Println("Updating existing workday")
		err = day.Update(database.DB)
		if err != nil {
			fmt.Printf("Error updating workday: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error updating workday: " + err.Error())
		}
	}

	fmt.Printf("Workday saved successfully. ID: %d\n", day.ID)

	return c.Redirect("/week-overview/" + strconv.Itoa(employeeID))
}

// Helper function to combine date and time
func combineDateAndTime(date time.Time, timeStr string) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, date.Location()), nil
}

func NewWorkday(c *fiber.Ctx) error {
	fmt.Println(c)
	employeeID, err := strconv.Atoi(c.Params("employeeId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid employee ID")
	}

	// Parse the date from the query parameter
	dateStr := c.Query("date")
	var date time.Time
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid date format")
		}
	} else {
		// If no date is provided, use today's date
		date = time.Now()
	}

	customers, err := models.GetAllCustomers(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("workday_form", fiber.Map{
		"Day": models.Day{
			EmployeeID: employeeID,
			Date:       date,
		},
		"Customers": customers,
	})
}

func GetWorkday(c *fiber.Ctx) error {
	employeeID, err := strconv.Atoi(c.Params("employeeId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid employee ID")
	}

	date, err := time.Parse("2006-01-02", c.Params("date"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid date format")
	}

	day, err := models.GetWorkday(database.DB, employeeID, date)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no workday exists, create a new one
			day = models.Day{EmployeeID: employeeID, Date: date}
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	customers, err := models.GetAllCustomers(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("workday_form", fiber.Map{
		"Day":       day,
		"Customers": customers,
	})
}

func calculateShiftHours(startTime, endTime, startMealTime, endMealTime time.Time) float64 {
	// Ensure end times are after start times
	if endTime.Before(startTime) {
		endTime = endTime.Add(24 * time.Hour)
	}
	if endMealTime.Before(startMealTime) {
		endMealTime = endMealTime.Add(24 * time.Hour)
	}

	shiftDuration := endTime.Sub(startTime)
	mealDuration := endMealTime.Sub(startMealTime)
	shiftHours := shiftDuration.Hours() - mealDuration.Hours()

	return shiftHours
}
