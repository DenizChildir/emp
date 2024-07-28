// internal/models/employee.go

package models

import (
	"database/sql"
	"math"
	"time"
)

type Employee struct {
	ID       int
	FullName string
	Active   bool
	Phone    int
	Email    string
	MsId     int
	ExcalID  int
	Rate     float64
}

func GetAllEmployees(db *sql.DB) ([]Employee, error) {
	rows, err := db.Query("SELECT * FROM employee")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		err := rows.Scan(&e.ID, &e.FullName, &e.Active, &e.Phone, &e.Email, &e.MsId, &e.ExcalID, &e.Rate)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func GetEmployee(db *sql.DB, id int) (Employee, error) {
	var e Employee
	err := db.QueryRow("SELECT * FROM employee WHERE id = ?", id).Scan(
		&e.ID, &e.FullName, &e.Active, &e.Phone, &e.Email, &e.MsId, &e.ExcalID, &e.Rate)
	return e, err
}

func (e *Employee) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO employee (fullName, active, phone, email, MsId, excalID, rate) VALUES (?, ?, ?, ?, ?, ?, ?)",
		e.FullName, e.Active, e.Phone, e.Email, e.MsId, e.ExcalID, e.Rate)
	return err
}

func (e *Employee) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE employee SET fullName=?, active=?, phone=?, email=?, MsId=?, excalID=?, rate=? WHERE id=?",
		e.FullName, e.Active, e.Phone, e.Email, e.MsId, e.ExcalID, e.Rate, e.ID)
	return err
}

func DeleteEmployee(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM employee WHERE id=?", id)
	return err
}

type EmployeeWeekData struct {
	Employee    Employee
	Days        [7]Day
	TotalHours  float64
	TotalIncome float64
}

func GetEmployeeWeekData(db *sql.DB, startDate time.Time) ([]EmployeeWeekData, error) {
	employees, err := GetAllEmployees(db)
	if err != nil {
		return nil, err
	}

	var employeeWeekData []EmployeeWeekData

	for _, emp := range employees {
		weekDays, err := GetWeekWorkdays(db, emp.ID, startDate)
		if err != nil {
			return nil, err
		}

		var weekData EmployeeWeekData
		weekData.Employee = emp
		weekData.Days = [7]Day{} // Initialize with empty Days

		totalHours := 0.0
		for _, day := range weekDays {
			// Calculate the day index (0 for Sunday, 1 for Monday, etc.)
			dayIndex := int(day.Date.Weekday())
			weekData.Days[dayIndex] = day
			totalHours += day.TotalHours()
		}

		weekData.TotalHours = math.Round(totalHours*100) / 100
		weekData.TotalIncome = math.Round(totalHours*emp.Rate*100) / 100

		employeeWeekData = append(employeeWeekData, weekData)
	}

	return employeeWeekData, nil
}
