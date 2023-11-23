package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Organisation struct {
	ID     int64
	Name   string
	City   string
	County string
	Jobs   []int64
}

type Job struct {
	ID        int64
	Type      string
	Rating    string
	VisaRoute string
}

func MustConnectDB() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "visa_sponsors_db",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("db connection successful")
	return db, nil
}
