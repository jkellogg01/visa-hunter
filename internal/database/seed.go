package database

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
)

func SeedDB() {
	db := MustConnectDB()
	defer db.Close()

	// err := initTable(db)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	file, err := os.Open("internal/database/data.csv")
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal("error reading file: ", err)
	}

	for i, row := range data {
		if i == 0 {
			continue
		}
		err := addCompany(
			db,
			row[0],
			row[1],
			row[2],
			row[3],
			row[4],
		)
		if err != nil {
			log.Fatal("error writing new row: ", err)
		}
	}
}

func addCompany(db *sql.DB, org_name string, city string, county string, job_type_rating string, visa_route string) error {
	result, err := db.Exec(
		"INSERT INTO companies (org_name, city, county, job_type_rating, visa_route) VALUES (?, ?, ?, ?, ?)",
		org_name, city, county, job_type_rating, visa_route,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()

	log.Printf("addCompany completed affecting %d row(s)", rows)
	return err
}

func initTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS companies")
	if err != nil {
		return err
	}

	log.Println("Cleared DB")

	_, err = db.Exec(`
		CREATE TABLE companies (
			org_id int AUTO INCREMENT,
			org_name varchar(255) NOT NULL,
			city varchar(255),
			county varchar(255),
			job_type_rating varchar(255),
			primary key (org_id)
		)
	`)
	if err != nil {
		return err
	}

	log.Println("Initialized DB")

	return nil
}
