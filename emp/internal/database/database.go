// internal/database/database.go

package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() error {
	var err error
	DB, err = sql.Open("mysql", "root:uuuu@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		return err
	}
	return DB.Ping()
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
