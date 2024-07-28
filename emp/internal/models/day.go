// internal/models/day.go

package models

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"
)

type Day struct {
	ID            int
	Date          time.Time
	StartTime     string
	EndTime       string
	EmployeeID    int
	ShiftHours    float64
	CustomerID    int
	StartMealTime string
	EndMealTime   string
}

func GetWorkday(db *sql.DB, employeeID int, date time.Time) (Day, error) {
	var d Day
	var dateStr string
	err := db.QueryRow(`
        SELECT id, date, startTime, endTime, startMealTime, endMealTime,
        employeeId, shiftHours, customerId FROM day WHERE employeeId = ? AND date = ?`,
		employeeID, date.Format("2006-01-02")).Scan(
		&d.ID, &dateStr, &d.StartTime, &d.EndTime, &d.StartMealTime, &d.EndMealTime,
		&d.EmployeeID, &d.ShiftHours, &d.CustomerID)

	if err != nil {
		return d, err
	}

	d.Date, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return d, err
	}

	return d, nil
}

func (d *Day) Create(db *sql.DB) error {
	result, err := db.Exec(`
        INSERT INTO day (date, startTime, endTime, startMealTime, endMealTime, employeeId, shiftHours, customerId) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.Date.Format("2006-01-02"), d.StartTime, d.EndTime, d.StartMealTime, d.EndMealTime, d.EmployeeID, d.ShiftHours, d.CustomerID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = int(id)
	return nil
}

func (d *Day) Update(db *sql.DB) error {
	_, err := db.Exec(`
        UPDATE day SET date=?, startTime=?, endTime=?, startMealTime=?, endMealTime=?, shiftHours=?, customerId=? 
        WHERE id=?`,
		d.Date.Format("2006-01-02"), d.StartTime, d.EndTime, d.StartMealTime, d.EndMealTime, d.ShiftHours, d.CustomerID, d.ID)
	return err
}

func GetWeekWorkdays(db *sql.DB, employeeID int, startDate time.Time) ([]Day, error) {
	endDate := startDate.AddDate(0, 0, 6)
	rows, err := db.Query("SELECT * FROM day WHERE employeeId = ? AND date BETWEEN ? AND ? ORDER BY date",
		employeeID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []Day
	for rows.Next() {
		var d Day
		var dateStr string
		err := rows.Scan(&d.ID, &dateStr, &d.StartTime, &d.EndTime, &d.EmployeeID, &d.ShiftHours, &d.CustomerID, &d.StartMealTime, &d.EndMealTime)
		if err != nil {
			return nil, err
		}
		d.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		days = append(days, d)
	}
	return days, nil
}
func GetDailyOverview(db *sql.DB, date time.Time) ([]struct {
	Day
	EmployeeName string
	CustomerName string
}, error) {
	query := `
		SELECT d.*, e.fullName as employeeName, c.company as customerName
		FROM day d
		JOIN employee e ON d.employeeId = e.id
		LEFT JOIN customer c ON d.customerId = c.id
		WHERE d.date = ?
		ORDER BY d.startTime
	`
	rows, err := db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overview []struct {
		Day
		EmployeeName string
		CustomerName string
	}

	for rows.Next() {
		var o struct {
			Day
			EmployeeName string
			CustomerName string
		}
		err := rows.Scan(
			&o.ID, &o.Date, &o.StartTime, &o.EndTime,
			&o.EmployeeID, &o.ShiftHours, &o.CustomerID,
			&o.EmployeeName, &o.CustomerName, &o.StartMealTime, &o.EndMealTime,
		)
		if err != nil {
			return nil, err
		}
		overview = append(overview, o)
	}

	return overview, nil
}

func (d *Day) TotalHours() float64 {
	if d.StartTime == "" || d.EndTime == "" {
		return 0
	}

	layout := "15:04:05" // Changed to include seconds

	startTime, err := time.Parse(layout, d.StartTime)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		return 0
	}
	endTime, err := time.Parse(layout, d.EndTime)
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		return 0
	}

	// Handle cases where end time is on the next day
	if endTime.Before(startTime) {
		endTime = endTime.Add(24 * time.Hour)
	}

	duration := endTime.Sub(startTime)

	// Subtract meal time if present
	if d.StartMealTime != "" && d.EndMealTime != "" {
		startMeal, err := time.Parse(layout, d.StartMealTime)
		if err == nil {
			endMeal, err := time.Parse(layout, d.EndMealTime)
			if err == nil {
				// Handle cases where meal end time is on the next day
				if endMeal.Before(startMeal) {
					endMeal = endMeal.Add(24 * time.Hour)
				}
				mealDuration := endMeal.Sub(startMeal)
				duration -= mealDuration
			}
		}
	}

	hours := duration.Hours()
	fmt.Println(math.Round(hours*100) / 100)
	return math.Round(hours*100) / 100 // Round to 2 decimal places
}
