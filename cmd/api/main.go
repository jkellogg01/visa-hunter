package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"visa-hunter/internal/db"
)

type Company struct {
	ID         int
	Name       string
	City       string
	County     string
	TypeRating string
	Route      string
}

func main() {
	db := db.MustConnectDB()
	defer db.Close()

	jobs, err := getFiltered(db, map[string]string{
		"city": "london",
	}, 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	for line, comp := range jobs {
		fmt.Println(line, comp)
	}
}

func getFiltered(db *sql.DB, params map[string]string, size int, page int) ([]Company, error) {
	var companies []Company

	var where strings.Builder
	if len(params) > 0 {
		where.WriteString("WHERE ")
	}
	for k, v := range params {
		where.WriteString(fmt.Sprintf("%s LIKE '%%%s%%'", k, v))
	}
	// log.Println(where.String())

	query := fmt.Sprintf("SELECT * FROM companies %s LIMIT %d, %d", where.String(), page*size, size)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comp Company
		if err := rows.Scan(&comp.ID, &comp.Name, &comp.City, &comp.County, &comp.TypeRating, &comp.Route); err != nil {
			return nil, err
		}
		companies = append(companies, comp)
	}

	return companies, nil
}
