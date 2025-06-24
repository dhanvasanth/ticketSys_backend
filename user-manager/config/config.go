package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/ticketing_db?parseTime=true")
	if err != nil {
		return err
	}
	
	if err = DB.Ping(); err != nil {
		return err
	}
	
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}