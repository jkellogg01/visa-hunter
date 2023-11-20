package db

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
)

func SeedDB() {
	db := MustConnectDB()
	defer db.Close()

	file, err := os.Open("db/data.csv")
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
