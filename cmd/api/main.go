package main

import "visa-hunter/internal/db"

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
}
