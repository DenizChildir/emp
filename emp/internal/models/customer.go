// internal/models/customer.go

package models

import (
	"database/sql"
)

type Customer struct {
	ID           int
	FullName     string
	Company      string
	ManagerEmail string
	CCEmail      string
	DutyLunch    bool
}

func GetAllCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query("SELECT * FROM customer")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.ID, &c.FullName, &c.Company, &c.ManagerEmail, &c.CCEmail, &c.DutyLunch)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func GetCustomer(db *sql.DB, id int) (Customer, error) {
	var c Customer
	err := db.QueryRow("SELECT * FROM customer WHERE id = ?", id).Scan(
		&c.ID, &c.FullName, &c.Company, &c.ManagerEmail, &c.CCEmail, &c.DutyLunch)
	return c, err
}

func (c *Customer) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO customer (fullName, company, managerEmail, ccEmail, dutyLunch) VALUES (?, ?, ?, ?, ?)",
		c.FullName, c.Company, c.ManagerEmail, c.CCEmail, c.DutyLunch)
	return err
}

func (c *Customer) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE customer SET fullName=?, company=?, managerEmail=?, ccEmail=?, dutyLunch=? WHERE id=?",
		c.FullName, c.Company, c.ManagerEmail, c.CCEmail, c.DutyLunch, c.ID)
	return err
}

func DeleteCustomer(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM customer WHERE id=?", id)
	return err
}
