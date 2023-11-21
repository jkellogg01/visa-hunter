package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"visa-hunter/internal/database"
)

func main() {
	// database.SeedDB()

	http.HandleFunc("/", handleIndex)

	http.HandleFunc("/index", handleIndexPage)

	log.Print(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	db := database.MustConnectDB()
	defer db.Close()

	pageqp := r.URL.Query().Get("page")
	pagenum, err := strconv.Atoi(pageqp)
	if err != nil {
		log.Println("page number absent or invalid, returning page zero")
		pagenum = 0
	}

	tmpl := template.Must(template.ParseFiles("views/index.html"))

	data, err := database.GetFiltered(db, map[string]string{}, 48, pagenum*48)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(w, data)
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	db := database.MustConnectDB()
	defer db.Close()

	pageqp := r.URL.Query().Get("page")
	pagenum, err := strconv.Atoi(pageqp)
	if err != nil {
		log.Println("Invalid page number, returning page zero")
		pagenum = 0
	}

	tmpl := template.Must(template.ParseFiles("views/index.html"))

	data, err := database.GetFiltered(db, map[string]string{}, 48, pagenum*48)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "results", data)
}
